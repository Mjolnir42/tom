/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package meta // // import "github.com/mjolnir42/tom/internal/model/meta"

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/rest"
	"github.com/mjolnir42/tom/internal/stmt"
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdNamespaceAttrAdd)

	registry = append(registry, function{
		cmd:    proto.CmdNamespaceAttrRemove,
		handle: namespaceAttrRemove,
	})
}

func namespaceAttrRemove(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.NamespaceAttrRemove)
}

func exportNamespaceAttrRemove(result *proto.Result, r *msg.Result) {
	result.Namespace = &[]proto.Namespace{}
	*result.Namespace = append(*result.Namespace, r.Namespace...)
}

// NamespaceAttrAdd function
func (m *Model) NamespaceAttrRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer rest.PanicCatcher(w, m.x.LM)

	request := msg.New(
		r, params,
		proto.CmdNamespaceAttrRemove,
		msg.SectionNamespace,
		proto.ActionAttrRemove,
	)

	req := proto.Request{}
	if err := rest.DecodeJSONBody(r, &req); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}
	request.Namespace = *req.Namespace
	request.Namespace.TomID = params.ByName(`tomID`)
	if err := request.Namespace.ParseTomID(); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}

	if err := proto.ValidNamespace(
		request.Namespace.Name,
	); err != nil {
		m.x.ReplyBadRequest(&w, &request, err)
		return
	}

	if !m.x.IsAuthorized(&request) {
		m.x.ReplyForbidden(&w, &request)
		return
	}

	m.x.HM.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	m.x.Send(&w, &result, exportNamespaceAttrRemove)
}

// attributeRemove ...
func (h *NamespaceWriteHandler) attributeRemove(q *msg.Request, mr *msg.Result) {

	var (
		err                      error
		tx                       *sql.Tx
		attrID, dictID, attrType string
	)

	if tx, err = h.conn.Begin(); err != nil {
		mr.ServerError(err)
		return
	}

	if _, err = tx.Exec(`SET CONSTRAINTS ALL DEFERRED;`); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}

	// loop over attributes to delete
	for _, attribute := range q.Namespace.Attributes {
		if err = tx.QueryRow(
			stmt.NamespaceAttributeSelect,
			q.Namespace.Name,
			attribute.Key,
		).Scan(
			&attrID,
			&dictID,
			&attrType,
		); err == sql.ErrNoRows {
			mr.NotFound(err)
			tx.Rollback()
			return
		} else if err != nil {
			mr.ServerError(err)
			tx.Rollback()
			return
		}
		switch {
		case attrType == proto.AttributeUnique && attribute.Unique:
		case attrType == proto.AttributeStandard && !attribute.Unique:
		default:
			mr.ServerError(
				fmt.Errorf("Mismatched attribute type: %s",
					attribute.Key,
				))
			tx.Rollback()
			return
		}

		switch attrType {
		case proto.AttributeUnique:
			for _, statement := range []string{
				stmt.ContainerDelNamespaceUniqValues,
				stmt.OrchestrationUniqAttrRemove,
				stmt.RuntimeDelNamespaceUniqValues,
				stmt.ServerUniqAttrRemove,
				stmt.SocketUniqAttrRemove,
				stmt.NamespaceUniqAttrRemoveValue,
				// TODO ix.deployment_group_unique_attribute_values
				// TODO ix.endpoint_unique_attribute_values
				// TODO ix.functional_component_unique_attribute_values
				// TODO ix.product_unique_attribute_values
				// TODO ix.technical_service_unique_attribute_values
				// TODO ix.top_level_service_unique_attribute_values
				// TODO yp.corporate_domain_unique_attribute_values
				// TODO yp.domain_unique_attribute_values
				// TODO yp.information_system_unique_attribute_values
				// TODO yp.service_unique_attribute_values
			} {
				if _, err = tx.Exec(
					statement,
					attrID,
					dictID,
				); err != nil {
					mr.ServerError(err)
					tx.Rollback()
					return
				}
			}
		case proto.AttributeStandard:
			for _, statement := range []string{
				stmt.ContainerDelNamespaceStdValues,
				stmt.OrchestrationStdAttrRemove,
				stmt.RuntimeDelNamespaceStdValues,
				stmt.ServerStdAttrRemove,
				stmt.SocketStdAttrRemove,
				stmt.NamespaceStdAttrRemoveValue,
				// TODO ix.deployment_group_standard_attribute_values
				// TODO ix.endpoint_standard_attribute_values
				// TODO ix.functional_component_standard_attribute_values
				// TODO ix.product_standard_attribute_values
				// TODO ix.technical_service_standard_attribute_values
				// TODO ix.top_level_service_standard_attribute_values
				// TODO yp.corporate_domain_standard_attribute_values
				// TODO yp.domain_standard_attribute_values
				// TODO yp.information_system_standard_attribute_values
				// TODO yp.service_standard_attribute_values
			} {
				if _, err = tx.Exec(
					statement,
					attrID,
					dictID,
				); err != nil {
					mr.ServerError(err)
					tx.Rollback()
					return
				}
			}
		default:
			mr.ServerError(
				fmt.Errorf("Unhandled attribute type: %s",
					attribute.Key,
				))
			tx.Rollback()
			return
		}
		for _, statement := range []string{
			stmt.NamespaceStdAttrRemove,
			stmt.NamespaceUniqAttrRemove,
		} {
			if _, err = tx.Exec(
				statement,
				attrID,
				dictID,
			); err != nil {
				mr.ServerError(err)
				tx.Rollback()
				return
			}
		}
		if _, err = tx.Exec(
			stmt.NamespaceAttrRemove,
			attribute.Key,
			dictID,
		); err != nil {
			mr.ServerError(err)
			tx.Rollback()
			return
		}
	}

	if err = tx.Commit(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.Namespace = append(mr.Namespace, q.Namespace)
	mr.OK()

}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
