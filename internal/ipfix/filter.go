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

type procFilter struct {
	conf    config.SettingsIPFIX
	inpipe  chan []byte
	outpipe chan []byte
	mirror  chan []byte
	pool    *sync.Pool
	lm      *lhm.LogHandleMap
}

func newFilter(conf config.SettingsIPFIX, inpipe, outpipe, mirror chan []byte, pool *sync.Pool, lm *lhm.LogHandleMap) (*procFilter, error) {
	f := &procFilter{
		conf:    conf,
		inpipe:  inpipe,
		outpipe: outpipe,
		mirror:  mirror,
		pool:    pool,
		lm:      lm,
	}
	// read from inpipe
	// filter
	// if processing contains aggregate
	//		copy to mirror
	// write to outpipe

	return f, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
