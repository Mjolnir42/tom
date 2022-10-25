/*-
 * Copyright (c) 2021-2022, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto //

const (
	CmdUser       = ModelIAM + `::` + EntityUser + `:`
	CmdUserAdd    = ModelIAM + `::` + EntityUser + `:` + ActionAdd
	CmdUserList   = ModelIAM + `::` + EntityUser + `:` + ActionList
	CmdUserRemove = ModelIAM + `::` + EntityUser + `:` + ActionRemove
	CmdUserShow   = ModelIAM + `::` + EntityUser + `:` + ActionShow
	CmdUserUpdate = ModelIAM + `::` + EntityUser + `:` + ActionUpdate
	CmdMachEnrol  = ModelIAM + `::` + EntityMachine + `:` + ActionEnrolment
)

func init() {
	Commands[CmdUserAdd] = CmdDef{
		Method:      MethodPOST,
		Path:        `/user/`,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{},
	}
	Commands[CmdUserList] = CmdDef{
		Method:      MethodGET,
		Path:        `/user/`,
		Body:        false,
		ResultTmpl:  TemplateList,
		Placeholder: []string{},
	}
	Commands[CmdUserRemove] = CmdDef{
		Method:      MethodDELETE,
		Path:        `/user/` + PlHoldTomID,
		Body:        false,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdUserShow] = CmdDef{
		Method:      MethodGET,
		Path:        `/user/` + PlHoldTomID,
		Body:        false,
		ResultTmpl:  TemplateDetail,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdUserUpdate] = CmdDef{
		Method:      MethodPATCH,
		Path:        `/user/` + PlHoldTomID,
		Body:        false,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
	Commands[CmdMachEnrol] = CmdDef{
		Method:      MethodPUT,
		Path:        `/machine/` + PlHoldTomID,
		Body:        true,
		ResultTmpl:  TemplateCommand,
		Placeholder: []string{PlHoldTomID},
	}
}

// User ...
type User struct {
	LibraryName    string      `json:"library-name"`
	FirstName      string      `json:"first-name,omitempty"`
	LastName       string      `json:"last-name,omitempty"`
	UserName       string      `json:"user-name"`
	EmployeeNumber string      `json:"employee-number,omitempty"`
	MailAddress    string      `json:"mailaddress,omitempty"`
	ExternalID     string      `json:"external-ref,omitempty"`
	Credential     *Credential `json:"credential,omitempty"`
	IsActive       bool        `json:"is-active"`
	IsDeleted      bool        `json:"is-deleted"`
	CreatedAt      string      `json:"createdAt,omitempty"`
	CreatedBy      string      `json:"createdBy,omitempty"`
	ID             string      `json:"-"`
	LibraryID      string      `json:"-"`
	TomID          string      `json:"-"`
}

type Credential struct {
	Category string `json:"category"`
	Value    string `json:"value"`
}

func NewUser() *User {
	return &User{
		Credential: &Credential{},
	}
}

func (u *User) ParseTomID() error {
	var typeID string
	switch {
	case u.TomID == ``:
		return ErrEmptyTomID
	case isTomIDFormatDNS(u.TomID):
		u.UserName, u.LibraryName, typeID = parseTomIDFormatDNS(u.TomID)
		return assessTomID(EntityMachine, typeID)
	case isTomIDFormatURI(u.TomID):
		u.UserName, u.LibraryName, typeID = parseTomIDFormatURI(u.TomID)
		return assessTomID(EntityMachine, typeID)
	default:
		return ErrInvalidTomID
	}
}

// Serialize ...
func (u *User) Serialize() []byte {
	data := make([]byte, 0)
	data = append(data, []byte(u.LibraryName)...)
	data = append(data, []byte(u.FirstName)...)
	data = append(data, []byte(u.LastName)...)
	data = append(data, []byte(u.UserName)...)
	data = append(data, []byte(u.EmployeeNumber)...)
	data = append(data, []byte(u.MailAddress)...)
	data = append(data, []byte(u.ExternalID)...)
	if u.Credential != nil {
		data = append(data, []byte(u.Credential.Serialize())...)
	}
	return data
}

// Serialize ...
func (c *Credential) Serialize() []byte {
	data := make([]byte, 0)
	data = append(data, []byte(c.Category)...)
	data = append(data, []byte(c.Value)...)
	return data
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
