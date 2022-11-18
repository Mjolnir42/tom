/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
	"encoding/binary"
	"net"
	"time"

	"github.com/vmware/go-ipfix/pkg/entities"
	"github.com/vmware/go-ipfix/pkg/registry"
)

func init() {
	registry.LoadRegistry()
}

func (f *procFilter) createIPFIX(pack *MessagePack) error {
	var err error
	var b []byte

	data := entities.NewSet(false)
	templateID := uint16(4379)
	if err = data.PrepareSet(entities.Data, templateID); err != nil {
		return err
	}

	var element *entities.InfoElement
	var ieval entities.InfoElementWithValue

	for i := range pack.records {
		elements := make([]entities.InfoElementWithValue, 0)

		if pack.records[i].IPVersion != 0 {
			if element, err = registry.GetInfoElement(`ipVersion`, registry.IANAEnterpriseID); err != nil {
				return err
			}
			if ieval, err = entities.DecodeAndCreateInfoElementWithValue(element, []byte{pack.records[i].IPVersion}); err != nil {
				return err
			}
			elements = append(elements, ieval)
		}
		switch pack.records[i].IPVersion {
		case 4:
			if element, err = registry.GetInfoElement(`sourceIPv4Address`, registry.IANAEnterpriseID); err != nil {
				return err
			}
			if ieval, err = entities.DecodeAndCreateInfoElementWithValue(element, net.ParseIP(pack.records[i].SrcAddress)); err != nil {
				return err
			}
			elements = append(elements, ieval)

			if element, err = registry.GetInfoElement(`destinationIPv4Address`, registry.IANAEnterpriseID); err != nil {
				return err
			}
			if ieval, err = entities.DecodeAndCreateInfoElementWithValue(element, net.ParseIP(pack.records[i].DstAddress)); err != nil {
				return err
			}
			elements = append(elements, ieval)
		case 6:
			if element, err = registry.GetInfoElement(`sourceIPv6Address`, registry.IANAEnterpriseID); err != nil {
				return err
			}
			if ieval, err = entities.DecodeAndCreateInfoElementWithValue(element, net.ParseIP(pack.records[i].SrcAddress)); err != nil {
				return err
			}
			elements = append(elements, ieval)

			if element, err = registry.GetInfoElement(`destinationIPv6Address`, registry.IANAEnterpriseID); err != nil {
				return err
			}
			if ieval, err = entities.DecodeAndCreateInfoElementWithValue(element, net.ParseIP(pack.records[i].DstAddress)); err != nil {
				return err
			}
			elements = append(elements, ieval)
		}
		if pack.records[i].ProtocolID != 0 {
			if element, err = registry.GetInfoElement(`protocolIdentifier`, registry.IANAEnterpriseID); err != nil {
				return err
			}
			if ieval, err = entities.DecodeAndCreateInfoElementWithValue(element, []byte{pack.records[i].ProtocolID}); err != nil {
				return err
			}
			elements = append(elements, ieval)
		}

		if pack.records[i].SrcPort != 0 {
			b = make([]byte, 2)
			binary.BigEndian.PutUint16(b, pack.records[i].SrcPort)
			if element, err = registry.GetInfoElement(`sourceTransportPort`, registry.IANAEnterpriseID); err != nil {
				return err
			}
			if ieval, err = entities.DecodeAndCreateInfoElementWithValue(element, b); err != nil {
				return err
			}
			elements = append(elements, ieval)
		}

		if pack.records[i].DstPort != 0 {
			b = make([]byte, 2)
			binary.BigEndian.PutUint16(b, pack.records[i].DstPort)
			if element, err = registry.GetInfoElement(`destinationTransportPort`, registry.IANAEnterpriseID); err != nil {
				return err
			}
			if ieval, err = entities.DecodeAndCreateInfoElementWithValue(element, b); err != nil {
				return err
			}
			elements = append(elements, ieval)
		}

		b = make([]byte, 8)
		binary.BigEndian.PutUint64(b, pack.records[i].OctetCount)
		if element, err = registry.GetInfoElement(`octetDeltaCount`, registry.IANAEnterpriseID); err != nil {
			return err
		}
		if ieval, err = entities.DecodeAndCreateInfoElementWithValue(element, b); err != nil {
			return err
		}
		elements = append(elements, ieval)

		b = make([]byte, 8)
		binary.BigEndian.PutUint64(b, pack.records[i].PacketCount)
		if element, err = registry.GetInfoElement(`packetDeltaCount`, registry.IANAEnterpriseID); err != nil {
			return err
		}
		if ieval, err = entities.DecodeAndCreateInfoElementWithValue(element, b); err != nil {
			return err
		}
		elements = append(elements, ieval)

		b = make([]byte, 8)
		binary.BigEndian.PutUint64(b, uint64(pack.records[i].StartMilli.UnixMilli()))
		if element, err = registry.GetInfoElement(`flowStartMilliseconds`, registry.IANAEnterpriseID); err != nil {
			return err
		}
		if ieval, err = entities.DecodeAndCreateInfoElementWithValue(element, b); err != nil {
			return err
		}
		elements = append(elements, ieval)

		b = make([]byte, 8)
		binary.BigEndian.PutUint64(b, uint64(pack.records[i].EndMilli.UnixMilli()))
		if element, err = registry.GetInfoElement(`flowEndMilliseconds`, registry.IANAEnterpriseID); err != nil {
			return err
		}
		if ieval, err = entities.DecodeAndCreateInfoElementWithValue(element, b); err != nil {
			return err
		}
		elements = append(elements, ieval)

		b = make([]byte, 2)
		binary.BigEndian.PutUint16(b, uint16(pack.records[i].TcpControlBits))
		if element, err = registry.GetInfoElement(`tcpControlBits`, registry.IANAEnterpriseID); err != nil {
			return err
		}
		if ieval, err = entities.DecodeAndCreateInfoElementWithValue(element, b); err != nil {
			return err
		}
		elements = append(elements, ieval)

		data.AddRecord(elements, templateID)
	}

	data.UpdateLenInHeader()
	msg := entities.NewMessage(false)
	msgLen := entities.MsgHeaderLength + data.GetSetLength()
	msg.SetVersion(10)
	msg.SetObsDomainID(pack.header.ClientID)
	msg.SetMessageLen(uint16(msgLen))
	msg.SetExportTime(uint32(time.Now().Unix()))

	f.lock.RLock()
	seq, ok := f.sequences[pack.header.ClientID]
	f.lock.RUnlock()
	if !ok {
		seq = uint32(0)
	}
	seq = seq + data.GetNumberOfRecords()
	msg.SetSequenceNum(seq)
	f.lock.Lock()
	f.sequences[pack.header.ClientID] = seq
	f.lock.Unlock()

	pack.ipfix = make([]byte, msgLen)
	copy(pack.ipfix[:entities.MsgHeaderLength], msg.GetMsgHeader())
	copy(pack.ipfix[entities.MsgHeaderLength:entities.MsgHeaderLength+entities.SetHeaderLen], data.GetHeaderBuffer())
	index := entities.MsgHeaderLength + entities.SetHeaderLen
	for _, record := range data.GetRecords() {
		len := record.GetRecordLength()
		copy(pack.ipfix[index:index+len], record.GetBuffer())
		index += len
	}
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
