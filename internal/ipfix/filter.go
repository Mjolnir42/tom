/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
	"sync"
	"time"

	flow "github.com/EdgeCast/vflow/ipfix"
	"github.com/mjolnir42/flowdata"
	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/config"
	"github.com/vmware/go-ipfix/pkg/entities"
	"github.com/vmware/go-ipfix/pkg/registry"
)

var (
	// templates memory cache
	mCache flow.MemCache
)

type procFilter struct {
	conf         config.SettingsIPFIX
	quit         chan interface{}
	exit         chan interface{}
	err          chan error
	wg           sync.WaitGroup
	pool         *sync.Pool
	lock         sync.RWMutex
	lm           *lhm.LogHandleMap
	pipe         chan IPFIXMessage
	pipeConvert  chan IPFIXMessage
	pipeFilter   chan *MessagePack
	pipeOutput   chan *MessagePack
	outpipeIPFIX chan IPFIXMessage
	outpipeJSON  chan []byte
	outpipeFDR   chan flowdata.Record
	outpipeRawJS chan []byte
	fOutJSN      bool
	fRawJSN      bool
	fFmtJSN      string
	fOutIPFIX    bool
	mux          *ipfixMux
	parsedRules  []config.Rule
	sequences    map[uint32]uint32
	templates    map[uint16]entities.Set
	TemplateInt  time.Duration
}

func newFilter(conf config.SettingsIPFIX, mux *ipfixMux, pool *sync.Pool, lm *lhm.LogHandleMap) (*procFilter, error) {
	f := &procFilter{
		conf:      conf,
		quit:      make(chan interface{}),
		exit:      make(chan interface{}),
		err:       make(chan error),
		mux:       mux,
		pool:      pool,
		lm:        lm,
		templates: make(map[uint16]entities.Set),
		sequences: make(map[uint32]uint32),
	}
	switch f.conf.Refresh {
	case ``:
		f.TemplateInt = 120 * time.Second
	default:
		d, err := time.ParseDuration(f.conf.Refresh)
		if err != nil {
			return nil, err
		}
		f.TemplateInt = d
	}
	// read from inpipe
	// filter
	// if processing contains aggregate
	//		copy to mirror
	// write to outpipe
	if conf.Forwarding {
		for _, c := range conf.Clients {
			if !c.Enabled {
				continue
			}
			switch c.ForwardProto {
			case ProtoUDP, ProtoTCP, ProtoTLS:
				f.fOutIPFIX = true
			case ProtoJSON:
				f.fOutJSN = true
				f.fRawJSN = c.Unfiltered
				f.fFmtJSN = c.Format
			}
		}
	}

	if err := f.parseRules(); err != nil {
		return nil, err
	}

	// outFLT is the mux output to the filter
	f.pipe = f.mux.pipe(`outFLT`)

	f.pipeConvert = make(chan IPFIXMessage, 64)
	f.pipeFilter = make(chan *MessagePack, 64)
	f.pipeOutput = make(chan *MessagePack, 64)

	// inFLT is the mux input for filtered ipfix
	f.outpipeIPFIX = f.mux.pipe(`inFLT`)
	// inFLJ is the mux input for filtered JSON
	f.outpipeJSON = f.mux.jsonPipe(`inFLJ`)
	// inFLR is the mux input for filtered flowdata.Record
	f.outpipeFDR = f.mux.flowdataPipe(`inFLR`)
	// outJSN is the mux output for JSON data to the client
	f.outpipeRawJS = f.mux.jsonPipe(`outJSN`)

	// setup in-memory template cache
	mCache = flow.GetCache(f.conf.TemplFile)

	for i := 0; i < 8; i++ {
		f.wg.Add(1)
		go f.conversionWorker()

		f.wg.Add(1)
		go f.filterWorker()

		f.wg.Add(1)
		go f.outputWorker()
	}

	f.createTemplate()

	f.wg.Add(1)
	go f.run()
	return f, nil
}

func (f *procFilter) run() {
	defer f.wg.Done()
	f.lm.GetLogger(`application`).Infoln(`ipfix.filter: module running`)

runloop:
	for {
		select {
		case <-f.quit:
			f.lm.GetLogger(`application`).Infoln(`ipfix.filter: shutdown signal received`)
			break runloop
		case frame := <-f.pipe:
			select {
			case f.pipeConvert <- frame:
			default:
			}
		case err := <-f.err:
			f.lm.GetLogger(`error`).Errorln(`ipfix.filter:`, err)
		case <-f.exit:
			break runloop
		case <-time.Tick(f.TemplateInt):
			for msg := range f.GetTemplateMsg(uint16(4739)) {
				f.outpipeIPFIX <- msg
			}
		case <-time.Tick(5 * time.Minute):
			if err := mCache.Dump(f.conf.TemplFile); err != nil {
				// save in memory template cache on disk
				f.err <- err
			}
		}
	}
}

func (f *procFilter) Err() chan error {
	return f.err
}

func (f *procFilter) Exit() chan interface{} {
	return f.exit
}

func (f *procFilter) Stop() chan error {
	go func(e chan error) {
		close(f.quit)
		f.wg.Wait()
		close(e)
	}(f.Err())
	return f.Err()
}

func (f *procFilter) createTemplate() {
	templateID := uint16(4379)
	templateSet := entities.NewSet(false)
	if err := templateSet.PrepareSet(entities.Template, templateID); err != nil {
		f.err <- err
		close(f.exit)
		return
	}

	elements := make([]entities.InfoElementWithValue, 0)
	for _, field := range []string{
		`octetDeltaCount`,
		`packetDeltaCount`,
		`protocolIdentifier`,
		`tcpControlBits`,
		`sourceTransportPort`,
		`sourceIPv4Address`,
		`destinationTransportPort`,
		`destinationIpv4Address`,
		`ingressInterface`,
		`egressInterface`,
		`sourceIPv6Address`,
		`destinationIPv6Address`,
		`ipVersion`,
		`flowDirection`,
		`exporterIPv4Address`,
		`exporterIPv6Address`,
		`exportingProcessID`,
		`flowStartMilliseconds`,
		`flowEndMilliseconds`,
	} {
		var element *entities.InfoElement
		var err error

		if element, err = registry.GetInfoElement(field, registry.IANAEnterpriseID); err != nil {
			f.err <- err
			close(f.exit)
			return
		}
		ie, _ := entities.DecodeAndCreateInfoElementWithValue(element, nil)
		elements = append(elements, ie)
	}

	templateSet.AddRecord(elements, templateID)
	templateSet.UpdateLenInHeader()

	f.lock.Lock()
	defer f.lock.Unlock()

	f.templates[templateID] = templateSet
}

func (f *procFilter) GetTemplateMsg(tID uint16) chan IPFIXMessage {
	ret := make(chan IPFIXMessage)

	go func(out chan IPFIXMessage) {
		f.lock.RLock()
		defer f.lock.RUnlock()

		set := f.templates[tID]

		var obsDomID, seqNo uint32
		for obsDomID = range f.sequences {
			seqNo = f.sequences[obsDomID]

			msg := entities.NewMessage(false)
			msgLen := entities.MsgHeaderLength + set.GetSetLength()
			msg.SetVersion(10)
			msg.SetObsDomainID(obsDomID)
			msg.SetMessageLen(uint16(msgLen))
			msg.SetExportTime(uint32(time.Now().Unix()))
			msg.SetSequenceNum(seqNo)

			byteSlice := make([]byte, msgLen)
			copy(
				byteSlice[:entities.MsgHeaderLength],
				msg.GetMsgHeader(),
			)
			copy(
				byteSlice[entities.MsgHeaderLength:entities.MsgHeaderLength+entities.SetHeaderLen],
				set.GetHeaderBuffer(),
			)
			index := entities.MsgHeaderLength + entities.SetHeaderLen
			for _, record := range set.GetRecords() {
				len := record.GetRecordLength()
				copy(byteSlice[index:index+len], record.GetBuffer())
				index += len
			}
			out <- IPFIXMessage{
				body: byteSlice,
			}
		}
	}(ret)
	return ret
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
