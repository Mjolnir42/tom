/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * Copyright (c) 2021, Jörg Pernfuß <joerg.pernfuss@ionos.com>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package adm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/mjolnir42/tom/pkg/proto"
	"github.com/urfave/cli/v2"
)

// Exported functions

// WRAPPER
func Perform(rqType, path, tmpl string, body interface{}, c *cli.Context) error {
	var (
		err  error
		resp *resty.Response
	)

	if strings.HasSuffix(rqType, `body`) && body == nil {
		goto noattachment
	}

	switch rqType {
	case `get`:
		resp, err = GetReq(path)
	case `head`:
		resp, err = HeadReq(path)
	case `delete`:
		resp, err = DeleteReq(path)
	case `deletebody`:
		resp, err = DeleteReqBody(body, path)
	case `putbody`:
		resp, err = PutReqBody(body, path)
	case `postbody`:
		resp, err = PostReqBody(body, path)
	case `patchbody`:
		resp, err = PatchReqBody(body, path)
	}

	if err != nil {
		return err
	}
	return FormatOut(c, resp.Body(), tmpl)

noattachment:
	return fmt.Errorf(`Missing body to client request that requires it.`)
}

func MockOK(tmpl string, c *cli.Context) error {
	mock := `{"statusCode":200,"statusText":"OK","errors":[]}`
	return FormatOut(c, []byte(mock), tmpl)
}

func DecodedResponse(resp *resty.Response, res *proto.Result) error {
	if err := decodeResponse(resp, res); err != nil {
		return err
	}
	return checkApplicationError(res)
}

// DELETE
func DeleteReq(p string) (*resty.Response, error) {
	return handleRequestOptions(client.R().Delete(p))
}

func DeleteReqBody(body interface{}, p string) (*resty.Response, error) {
	return handleRequestOptions(
		client.R().SetBody(body).SetContentLength(true).Delete(p))
}

// GET
func GetReq(p string) (*resty.Response, error) {
	return handleRequestOptions(client.R().Get(p))
}

// HEAD
func HeadReq(p string) (*resty.Response, error) {
	return handleRequestOptions(client.R().Head(p))
}

// PATCH
func PatchReqBody(body interface{}, p string) (*resty.Response, error) {
	return handleRequestOptions(
		client.R().SetBody(body).SetContentLength(true).Patch(p))
}

// POST
func PostReqBody(body interface{}, p string) (*resty.Response, error) {
	return handleRequestOptions(
		client.R().SetBody(body).SetContentLength(true).Post(p))
}

// PUT
func PutReq(p string) (*resty.Response, error) {
	return handleRequestOptions(client.R().Put(p))
}

func PutReqBody(body interface{}, p string) (*resty.Response, error) {
	return handleRequestOptions(
		client.R().SetBody(body).SetContentLength(true).Put(p))
}

// Private functions

func handleRequestOptions(resp *resty.Response, err error) (*resty.Response, error) {
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() >= 300 {
		return resp, fmt.Errorf("Request error: %s, %s", resp.Status(), resp.String())
	}

	if !(async || jobSave) {
		return resp, nil
	}

	var result *proto.Result
	if err = decodeResponse(resp, result); err != nil {
		return nil, err
	}

	if jobSave {
		if result.StatusCode == 202 && result.JobID != "" {
			cache.SaveJob(result.JobID, result.JobType)
		}
	}

	if async {
		asyncWait(result)
	}
	return resp, nil
}

func asyncWait(result *proto.Result) {
	if !async {
		return
	}

	if result.StatusCode == 202 && result.JobID != "" {
		fmt.Fprintf(os.Stderr, "Waiting for job: %s\n", result.JobID)
		_, err := GetReq(fmt.Sprintf("/job/byID/%s/_processed", result.JobID))
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "Wait error: %s\n", err.Error())
		}
	}
}

func decodeResponse(resp *resty.Response, res *proto.Result) error {
	decoder := json.NewDecoder(bytes.NewReader(resp.Body()))
	return decoder.Decode(res)
}

// checkApplicationError tests the server result for
// application errors
func checkApplicationError(result *proto.Result) error {
	if result.StatusCode >= 300 {
		var s string
		// application errors
		if result.StatusCode == 404 {
			s = fmt.Sprintf("Object lookup error: %d - %s",
				result.StatusCode, result.ErrorText)
		} else {
			s = fmt.Sprintf("Application error: %d - %s",
				result.StatusCode, result.ErrorText)
		}
		m := []string{s}

		return fmt.Errorf(combineStrings(m...))
	}
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
