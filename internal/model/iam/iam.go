/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package iam // import "github.com/mjolnir42/tom/internal/model/iam"

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/tom/internal/rest"
)

type Model struct {
	x *rest.Rest
}

func New(x *rest.Rest) *Model {
	m := &Model{
		x: x,
	}
	return m
}

func (m *Model) RouteRegister(rt *httprouter.Router) {
	m.routeRegisterLibrary(rt)
	m.routeRegisterUser(rt)
	m.routeRegisterTeam(rt)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
