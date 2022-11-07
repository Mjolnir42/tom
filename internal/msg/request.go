/*-
 * Copyright (c) 2020-2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package msg // import "github.com/mjolnir42/tom/internal/msg"

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/pkg/proto"
	uuid "github.com/satori/go.uuid"
)

type Request struct {
	ID         uuid.UUID
	Command    string
	Section    string
	Action     string
	RemoteAddr string
	UserIDLib  string
	AuthUser   string
	RequestURI string
	Reply      chan Result `json:"-"`
	Verbose    bool

	Auth Super

	Update UpdateData

	Container     proto.Container
	Library       proto.Library
	Namespace     proto.Namespace
	Orchestration proto.Orchestration
	Runtime       proto.Runtime
	Server        proto.Server
	Team          proto.Team
	User          proto.User
}

// New returns a Request
func New(r *http.Request, params httprouter.Params, cmd, sec, ac string) Request {
	returnChannel := make(chan Result, 1)
	identity := authUser(params)
	rq := Request{
		ID:         requestID(params),
		Command:    cmd,
		Section:    sec,
		Action:     ac,
		RequestURI: requestURI(params),
		RemoteAddr: remoteAddr(r),
		UserIDLib:  identity[0],
		AuthUser:   identity[1],
		Reply:      returnChannel,
	}
	rq.Verbose, _ = strconv.ParseBool(r.URL.Query().Get(`verbose`))
	switch sec {
	case SectionContainer:
		rq.Container = *(proto.NewContainer())
	case SectionLibrary:
		rq.Library = proto.Library{} // TODO when implenting model
	case SectionNamespace:
		rq.Namespace = *(proto.NewNamespace())
	case SectionOrchestration:
		rq.Orchestration = *(proto.NewOrchestration())
	case SectionRuntime:
		rq.Runtime = *(proto.NewRuntime())
	case SectionServer:
		rq.Server = *(proto.NewServer())
	case SectionTeam:
		rq.Team = proto.Team{} // TODO when implementing model
	case SectionUser:
		rq.User = proto.User{} // TODO when implementing model
	case SectionMachine:
		rq.User = *(proto.NewUser())
	}
	return rq
}

type UpdateData struct {
	Team proto.Team
	User proto.User
}

type Super struct {
	Nonce      []byte
	RequestURI string
	IDLib      string
	UserID     string
	Sig        []byte
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
