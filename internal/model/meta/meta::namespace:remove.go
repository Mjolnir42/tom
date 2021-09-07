/*-
 * Copyright (c) 2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package meta // import "github.com/mjolnir42/tom/internal/model/meta"

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/msg"
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

	request := msg.New(r, params)
	request.Section = msg.SectionNamespace
	request.Action = proto.ActionRemove
	request.Namespace = proto.Namespace{
		TomID: params.ByName(`tomID`),
	}

	if err := request.Namespace.ParseTomID(); err != nil {
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
	mr.NotImplemented() // TODO
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
