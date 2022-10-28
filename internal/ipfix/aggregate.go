/*-
 * Copyright (c) 2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package ipfix

import (
	"sync"

	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/config"
)

type procAggregate struct {
	conf   config.SettingsIPFIX
	mirror chan []byte
	pool   *sync.Pool
	lm     *lhm.LogHandleMap
}

func newAggregate(conf config.SettingsIPFIX, mirror chan []byte, pool *sync.Pool, lm *lhm.LogHandleMap) (*procAggregate, error) {
	m := &procAggregate{
		conf:   conf,
		mirror: mirror,
		pool:   pool,
		lm:     lm,
	}
	return m, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
