/*-
 * Copyright (c) 2020-2021, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

type Entity interface {
	String() string
	FormatDNS() string
	FormatTomID() string
	ParseTomID() error
	PropertyIterator() <-chan PropertyDetail
	SetTomID() Entity
	ExportName() string
	ExportNamespace() string
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
