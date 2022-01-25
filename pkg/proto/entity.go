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

func ParseTomID(s string) (error, string, Entity) {
	var name, namespace, entity string
	switch {
	case s == ``:
		return ErrEmptyTomID, ``, nil
	case isTomIDFormatDNS(s):
		name, namespace, entity = parseTomIDFormatDNS(s)
	case isTomIDFormatURI(s):
		name, namespace, entity = parseTomIDFormatURI(s)
	default:
		return ErrInvalidTomID, ``, nil
	}

	switch entity {
	case EntityContainer:
		return nil, entity, (&Container{
			Namespace: namespace,
			Name:      name,
		}).SetTomID()
	case EntityNamespace:
		return nil, entity, (&Namespace{
			Name: name,
		}).SetTomID()
	case EntityOrchestration:
		return nil, entity, (&Orchestration{
			Namespace: namespace,
			Name:      name,
		}).SetTomID()
	case EntityRuntime:
		return nil, entity, (&Runtime{
			Namespace: namespace,
			Name:      name,
		}).SetTomID()
	case EntityServer:
		return nil, entity, (&Server{
			Namespace: namespace,
			Name:      name,
		}).SetTomID()
	case EntitySocket:
		return nil, entity, (&Socket{
			Namespace: namespace,
			Name:      name,
		}).SetTomID()
	default:
		return ErrInvalidTomID, ``, nil
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
