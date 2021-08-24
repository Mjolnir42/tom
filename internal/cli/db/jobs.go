/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * Copyright (c) 2021, Jörg Pernfuß <joerg.pernfuss@ionos.com>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package db

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

func (d *DB) SaveJob(jid, jtype string) error {
	if err := d.Open(); err != nil {
		return err
	}
	defer d.Close()
	now := time.Now().UTC().Format(rfc3339Milli)

	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(`jobs`)).Bucket([]byte(`active`))
		id, _ := b.NextSequence()
		return b.Put(
			uitob(id),
			[]byte(fmt.Sprintf("%s|%s|%s", jid, now, jtype)),
		)
	})
}

func uitob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func btoui(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func btos(b []byte) string {
	return strconv.FormatUint(binary.BigEndian.Uint64(b), 10)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
