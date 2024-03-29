/*-
 * Copyright (c) 2020, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

// Result is the response wrapper of Tom's public API
type Result struct {
	RequestID  string `json:"requestID"`
	Command    string `json:"command"`
	StatusCode uint16 `json:"status"`
	ErrorText  string `json:"error"`
	// Job information is set for StatusCode 202 (async processing)
	JobID   string `json:"jobId,omitempty"`
	JobType string `json:"jobType,omitempty"`

	Container           *[]Container           `json:"container,omitempty"`
	ContainerHeader     *[]ContainerHeader     `json:"container-list,omitempty"`
	Flow                *[]Flow                `json:"flow,omitempty"`
	Library             *[]Library             `json:"library,omitempty"`
	Namespace           *[]Namespace           `json:"namespace,omitempty"`
	NamespaceHeader     *[]NamespaceHeader     `json:"namespace-list,omitempty"`
	Orchestration       *[]Orchestration       `json:"orchestration,omitempty"`
	OrchestrationHeader *[]OrchestrationHeader `json:"orchestration-list,omitempty"`
	Runtime             *[]Runtime             `json:"runtime,omitempty"`
	RuntimeHeader       *[]RuntimeHeader       `json:"runtime-list,omitempty"`
	Server              *[]Server              `json:"server,omitempty"`
	ServerHeader        *[]ServerHeader        `json:"server-list,omitempty"`
	Socket              *[]Socket              `json:"socket,omitempty"`
	Team                *[]Team                `json:"team,omitempty"`
	User                *[]User                `json:"user,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
