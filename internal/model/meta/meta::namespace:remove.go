/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package meta // import "github.com/mjolnir42/tom/internal/model/meta"

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/stmt"
	"github.com/mjolnir42/tom/pkg/proto"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdNamespaceRemove)

	registry = append(registry, function{
		cmd:    proto.CmdNamespaceRemove,
		handle: namespaceRemove,
	})
}

func namespaceRemove(m *Model) httprouter.Handle {
	return m.x.Authenticated(m.NamespaceRemove)
}

func exportNamespaceRemove(result *proto.Result, r *msg.Result) {
	result.Namespace = &[]proto.Namespace{}
	*result.Namespace = append(*result.Namespace, r.Namespace...)
}

// NamespaceRemove function
func (m *Model) NamespaceRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(
		r, params,
		proto.CmdNamespaceRemove,
		msg.SectionNamespace,
		proto.ActionRemove,
	)
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
	m.x.Send(&w, &result, exportNamespaceRemove)
}

// remove deletes a specific namespace
func (h *NamespaceWriteHandler) remove(q *msg.Request, mr *msg.Result) {
	var (
		err                                                    error
		tx                                                     *sql.Tx
		rows                                                   *sql.Rows
		attrID, dictID, attrType, dictName, createdAt, userUID string
		attributes                                             []string
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

	// query dictionaryID
	if err = tx.QueryRow(
		stmt.NamespaceTxShow,
		q.Namespace.Name,
	).Scan(
		&dictID,
		&dictName,
		&createdAt,
		&userUID,
	); err == sql.ErrNoRows {
		mr.NotFound(err)
		tx.Rollback()
		return
	} else if err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}

	// query all attributes
	if rows, err = tx.Query(
		stmt.NamespaceTxSelectAttributes,
		dictID,
	); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}

	for rows.Next() {
		var attr string
		if err = rows.Scan(
			&attr,
			&attrType,
			&createdAt,
			&userUID,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			tx.Rollback()
			return
		}
		attributes = append(attributes, attr)
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}

	// loop over all namespace attributes
	for _, attribute := range attributes {
		if err = tx.QueryRow(
			stmt.NamespaceAttributeSelect,
			q.Namespace.Name,
			attribute,
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

		// remove attribute values
		switch attrType {
		case proto.AttributeUnique:
			for _, statement := range []string{
				stmt.ContainerDelNamespaceUniqValues,
				stmt.OrchestrationDelNamespaceUniqValues,
				stmt.RuntimeDelNamespaceUniqValues,
				stmt.ServerDelNamespaceUniqValues,
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
				stmt.OrchestrationDelNamespaceStdValues,
				stmt.RuntimeDelNamespaceStdValues,
				stmt.ServerDelNamespaceStdValues,
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
				fmt.Errorf("Unhandled attribute type: %s (%s)",
					attrType,
					attribute,
				))
			tx.Rollback()
			return
		}

		// remove attribute from meta.standard_attribute and meta.unique_attribute
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
		// remove attribute from meta.attribute
		if _, err = tx.Exec(
			stmt.NamespaceAttrRemove,
			attribute,
			dictID,
		); err != nil {
			mr.ServerError(err)
			tx.Rollback()
			return
		}
	}

	// remove all containers in the namespace
	for _, statement := range []string{
		// TODO asset.Socket
		stmt.ContainerDelNamespaceLinking,
		stmt.ContainerDelNamespaceParent,
		stmt.ContainerDelNamespace,
		stmt.OrchestrationDelNamespaceLinking,
		stmt.OrchestrationDelNamespaceParent,
		stmt.OrchestrationDelNamespace,
		stmt.RuntimeDelNamespaceLinking,
		stmt.RuntimeDelNamespaceParent,
		stmt.RuntimeDelNamespace,
		stmt.ServerDelNamespaceLinking,
		stmt.ServerDelNamespaceParent,
		stmt.ServerDelNamespace,
	} {
		if _, err = tx.Exec(
			statement,
			dictID,
		); err != nil {
			mr.ServerError(err)
			tx.Rollback()
			return
		}
	}

	// remove the namespace
	if _, err = tx.Stmt(
		h.stmtRemove,
	).Exec(
		dictID,
		q.Namespace.Name,
	); err != nil {
		mr.ServerError(err)
		tx.Rollback()
		return
	}

	if err = tx.Commit(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.Namespace = append(mr.Namespace, q.Namespace)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
