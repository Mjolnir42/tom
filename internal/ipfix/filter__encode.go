//go:build exclude

/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
	"github.com/vmware/go-ipfix/pkg/entities"
	"github.com/vmware/go-ipfix/pkg/registry"
)

func init() {
	registry.LoadRegistry()
}

func (f *procFilter) createIPFIX(pack *MessagePack) {
// create template
	templateID = uint16(4379)
	templateSet := entities.NewSet(false)
	err := templateSet.PrepareSet(entities.Template, templateID)
	if err != nil {
		f.err <- err
	}

	elements := make([]entities.InfoElementWithValue, 0)
	element, err := registry.GetInfoElement("sourceIPv4Address", registry.IANAEnterpriseID)
	if err != nil {
		f.err <- err
	}
	ie, _ := entities.DecodeAndCreateInfoElementWithValue(element, nil)
	elements = append(elements, ie)

	element, err = registry.GetInfoElement("destinationIPv4Address", registry.IANAEnterpriseID)
	if err != nil {
		f.err <- err
	}
	ie, _ := entities.DecodeAndCreateInfoElementWithValue(element, nil)
	elements = append(elements, ie)

	templateSet.AddRecord(elements, templateID)

	// XXX
	bytesSent, err := exporter.SendSet(templateSet)
		setType := set.GetSetType()
		ep.updateTemplate(record.GetTemplateID(), record.GetOrderedElementList(), record.GetMinDataRecordLen())

		set.UpdateLenInHeader()
		bytesSent, err ep.createAndSendIPFIXMsg(set)

		msg := entities.NewMessage(false)
		msgLen := entities.MsgHeaderLength + set.GetSetLength()
		msg.SetVersion(10)
		msg.SetObsDomainID(ep.obsDomainID)
		msg.SetMessageLen(uint16(msgLen))
		msg.SetExportTime(uint32(time.Now().Unix()))
		if set.GetSetType() == entities.Data {
			ep.seqNumber = ep.seqNumber + set.GetNumberOfRecords()
		}
		msg.SetSequenceNum(ep.seqNumber)

		bytesSlice := make([]byte, msgLen)
		copy(bytesSlice[:entities.MsgHeaderLength], msg.GetMsgHeader())
		copy(bytesSlice[entities.MsgHeaderLength:entities.MsgHeaderLength+entities.SetHeaderLen], set.GetHeaderBuffer())
		index := entities.MsgHeaderLength + entities.SetHeaderLen
		for _, record := range set.GetRecords() {
				len := record.GetRecordLength()
				copy(bytesSlice[index:index+len], record.GetBuffer())
				index += len
		}


		// create data message
		dataSet := entities.NewSet(false)
		err = dataSet.PrepareSet(entities.Data, templateID)

		elements := make([]entities.InfoElementWithValue, 0)

		element, err = registry.GetInfoElement("sourceIPv4Address", registry.IANAEnterpriseID)
		ie, _ := entities.DecodeAndCreateInfoElementWithValue(element, net.ParseIP("1.2.3.4"))
		elements = append(elements, ie)

		element, err = registry.GetInfoElement("destinationIPv4Address", registry.IANAEnterpriseID)
		ie, _ = entities.DecodeAndCreateInfoElementWithValue(element, net.ParseIP("5.6.7.8"))
		elements = append(elements, ie)

		dataSet.AddRecord(elements, templateID)
		dataRecBuff := dataSet.GetRecords()[0].GetBuffer()

		bytesSent, err := exporter.SendSet(dataSet)

}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
