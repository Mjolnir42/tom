/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main // import "github.com/mjolnir42/tom/"

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/mjolnir42/lhm"
	"github.com/mjolnir42/tom/internal/msg"
)

func connectToDatabase(lm *lhm.LogHandleMap) *sql.DB {
	var err error
	var conn *sql.DB

	driver := `postgres`

	connect := fmt.Sprintf("dbname='%s' user='%s' password='%s' host='%s' port='%s' sslmode='%s' connect_timeout='%s'",
		// TomCfg is a global variable
		TomCfg.Database.Name,
		TomCfg.Database.User,
		TomCfg.Database.Pass,
		TomCfg.Database.Host,
		TomCfg.Database.Port,
		TomCfg.Database.TLSMode,
		TomCfg.Database.Timeout,
	)

	// enable handling of infinity timestamps
	pq.EnableInfinityTs(msg.NegTimeInf, msg.PosTimeInf)

	if conn, err = sql.Open(driver, connect); err != nil {
		lm.GetLogger(`error`).Fatal(err)
	}
	if err = conn.Ping(); err != nil {
		lm.GetLogger(`error`).Fatal(err)
	}

	lm.GetLogger(`application`).Print(`Connected to database`)
	if _, err = conn.Exec(`SET TIME ZONE 'UTC';`); err != nil {
		lm.GetLogger(`error`).Fatal(err)
	}
	if _, err = conn.Exec(`SET SESSION CHARACTERISTICS AS TRANSACTION ISOLATION LEVEL SERIALIZABLE;`); err != nil {
		lm.GetLogger(`error`).Fatal(err)
	}

	// size the connection pool
	conn.SetMaxIdleConns(5)
	conn.SetMaxOpenConns(15)
	conn.SetConnMaxLifetime(12 * time.Hour)

	return conn
}

func pingDatabase(lm *lhm.LogHandleMap, conn *sql.DB) {
	ticker := time.NewTicker(time.Second).C

	for {
		<-ticker
		if err := conn.Ping(); err != nil {
			lm.GetLogger(`error`).Print(err)
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
