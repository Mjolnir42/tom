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

type Specification struct {
	Name        string
	Placeholder map[string]string
	QueryParams *map[string]string
	Body        interface{}
}

// Exported functions

// WRAPPER
func Perform(cmd Specification, c *cli.Context) error {
	var (
		err  error
		path string
		resp *resty.Response
	)

	if _, ok := proto.Commands[cmd.Name]; !ok {
		goto unknownCommand
	}

	if proto.Commands[cmd.Name].Body && cmd.Body == nil {
		goto missingBody
	}

	path = proto.Commands[cmd.Name].Path
	switch {
	case len(proto.Commands[cmd.Name].Placeholder) != 0:
		if cmd.Placeholder == nil {
			goto improperSpec
		}

		for _, ph := range proto.Commands[cmd.Name].Placeholder {
			if _, ok := cmd.Placeholder[ph]; !ok {
				goto improperSpec
			}

			path = strings.Replace(
				path,
				ph,
				cmd.Placeholder[ph],
				1,
			)
		}
	}

	switch proto.Commands[cmd.Name].Method {
	case proto.MethodGET:
		resp, err = getReq(path, cmd.QueryParams)
	case proto.MethodHEAD:
		resp, err = headReq(path)
	case proto.MethodDELETE:
		switch {
		case proto.Commands[cmd.Name].Body:
			resp, err = deleteReqBody(cmd.Body, path)
		default:
			resp, err = deleteReq(path)
		}
	case proto.MethodPUT:
		switch {
		case proto.Commands[cmd.Name].Body:
			resp, err = putReqBody(cmd.Body, path)
		default:
			goto unhandledMethod
		}
	case proto.MethodPOST:
		switch {
		case proto.Commands[cmd.Name].Body:
			resp, err = postReqBody(cmd.Body, path)
		default:
			goto unhandledMethod
		}
	case proto.MethodPATCH:
		switch {
		case proto.Commands[cmd.Name].Body:
			resp, err = patchReqBody(cmd.Body, path)
		default:
			goto unhandledMethod
		}
	default:
		goto unhandledMethod
	}

	if err != nil {
		return err
	}
	return FormatOut(c, resp.Body(), proto.Commands[cmd.Name].ResultTmpl)

unknownCommand:
	return fmt.Errorf("Unknown command definition requested: %s",
		cmd.Name,
	)

missingBody:
	return fmt.Errorf(
		`Missing body to client request that requires it.`,
	)

unhandledMethod:
	return fmt.Errorf("Unhandled: Method:%s/Body:%t",
		proto.Commands[cmd.Name].Method,
		proto.Commands[cmd.Name].Body,
	)

improperSpec:
	return fmt.Errorf(
		`Specification contains uninitialized Placeholder map.`,
	)
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

// Private functions

// DELETE
func deleteReq(p string) (*resty.Response, error) {
	return handleRequestOptions(client.R().Delete(p))
}

func deleteReqBody(body interface{}, p string) (*resty.Response, error) {
	return handleRequestOptions(
		client.R().SetBody(body).SetContentLength(true).Delete(p))
}

// GET
func getReq(p string, q *map[string]string) (*resty.Response, error) {
	switch q {
	case nil:
		return handleRequestOptions(client.R().Get(p))
	default:
		return handleRequestOptions(client.R().SetQueryParams(*q).Get(p))
	}
}

// HEAD
func headReq(p string) (*resty.Response, error) {
	return handleRequestOptions(client.R().Head(p))
}

// PATCH
func patchReqBody(body interface{}, p string) (*resty.Response, error) {
	return handleRequestOptions(
		client.R().SetBody(body).SetContentLength(true).Patch(p))
}

// POST
func postReqBody(body interface{}, p string) (*resty.Response, error) {
	return handleRequestOptions(
		client.R().SetBody(body).SetContentLength(true).Post(p))
}

// PUT
func putReq(p string) (*resty.Response, error) {
	return handleRequestOptions(client.R().Put(p))
}

func putReqBody(body interface{}, p string) (*resty.Response, error) {
	return handleRequestOptions(
		client.R().SetBody(body).SetContentLength(true).Put(p))
}

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
		_, err := getReq(fmt.Sprintf("/job/byID/%s/_processed", result.JobID), nil)
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
