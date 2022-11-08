/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package supervisor // import "github.com/mjolnir42/tom/internal/model/supervisor/"

import (
	"database/sql"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/cred"
	"github.com/mjolnir42/tom/internal/msg"
	"github.com/mjolnir42/tom/internal/stmt"
	"github.com/mjolnir42/tom/pkg/proto"
	"golang.org/x/crypto/ed25519"
)

func init() {
	proto.AssertCommandIsDefined(proto.CmdSupervisorAuthEPK)

	registry = append(registry, function{
		cmd:    proto.CmdSupervisorAuthEPK,
		handle: supervisorAuthEPK,
	})
}

// REST calls into the supervisor are invalid
func supervisorAuthEPK(m *Model) httprouter.Handle {
	return m.x.Deny(m.SupervisorAuth)
}

// SupervisorAuth function
func (m *Model) SupervisorAuth(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	// wrapped in Rest.Deny() this function should never be reached
	request := msg.New(
		r, params,
		proto.CmdSupervisorAuthEPK,
		proto.EntitySupervisor,
		proto.ActionAuthenticateEPK,
	)
	m.x.ReplyForbidden(&w, &request)
}

// authenticateEPK ...
func (h *CoreHandler) authenticateEPK(q *msg.Request, mr *msg.Result) {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
