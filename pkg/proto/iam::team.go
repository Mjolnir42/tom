/*-
 * Copyright (c) 2021-2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

const (
	CmdTeam          = ModelIAM + `::` + EntityTeam + `:`
	CmdTeamAdd       = ModelIAM + `::` + EntityTeam + `:` + ActionAdd
	CmdTeamHdSet     = ModelIAM + `::` + EntityTeam + `:` + ActionHdSet
	CmdTeamHdUnset   = ModelIAM + `::` + EntityTeam + `:` + ActionHdUnset
	CmdTeamList      = ModelIAM + `::` + EntityTeam + `:` + ActionList
	CmdTeamMbrAdd    = ModelIAM + `::` + EntityTeam + `:` + ActionMbrAdd
	CmdTeamMbrList   = ModelIAM + `::` + EntityTeam + `:` + ActionMbrList
	CmdTeamMbrRemove = ModelIAM + `::` + EntityTeam + `:` + ActionMbrRemove
	CmdTeamMbrSet    = ModelIAM + `::` + EntityTeam + `:` + ActionMbrSet
	CmdTeamRemove    = ModelIAM + `::` + EntityTeam + `:` + ActionRemove
	CmdTeamShow      = ModelIAM + `::` + EntityTeam + `:` + ActionShow
	CmdTeamUpdate    = ModelIAM + `::` + EntityTeam + `:` + ActionUpdate
)

// User ...
type Team struct {
	LibraryName string  `json:"library-name"`
	TeamName    string  `json:"team-name"`
	ExternalID  string  `json:"external-ref,omitempty"`
	IsDeleted   bool    `json:"is-deleted"`
	CreatedAt   string  `json:"createdAt,omitempty"`
	CreatedBy   string  `json:"createdBy,omitempty"`
	TeamLead    *User   `json:"team-lead,omitempty"`
	Member      *[]User `json:"team-member,omitempty"`
	ID          string  `json:"-"`
	LibraryID   string  `json:"-"`
}

// Serialize ...
func (t *Team) Serialize() []byte {
	data := make([]byte, 0)
	data = append(data, []byte(t.LibraryName)...)
	data = append(data, []byte(t.TeamName)...)
	data = append(data, []byte(t.ExternalID)...)
	switch t.IsDeleted {
	case true:
		data = append(data, []byte(`true`)...)
	default:
		data = append(data, []byte(`false`)...)
	}
	if t.TeamLead != nil {
		data = append(data, t.TeamLead.Serialize()...)
	}
	if t.Member != nil {
		data = append(data, SerializeUserSlice(t.Member)...)
	}
	return data
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
