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
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/mattn/go-shellwords"
	"github.com/urfave/cli/v2"
)

func FormatOut(c *cli.Context, data []byte, cmd string) error {
	if string(data) == `` {
		return nil
	}

	if c.Bool(`json`) {
		return printJSON(data)
	}

	switch cmd {
	// TODO
	default:
		// hardwire JSON output for now
		return printJSON(data)
	}
}

func printJSON(data []byte) error {
	var outputDevice io.WriteCloser
	var proc *exec.Cmd
	var err error
	var processorARGS []string

	// check package variable, setup output
	switch postProcessor {
	case ``:
		outputDevice = wrapWNopCloser(os.Stdout)
	default:
		if processorARGS, err = shellwords.Parse(postProcessor); err != nil {
			return err
		}

		proc = exec.Command(
			processorARGS[0],
			processorARGS[1:]...,
		)
		proc.Stdout = os.Stdout
		proc.Stderr = ioutil.Discard
		if outputDevice, err = proc.StdinPipe(); err != nil {
			return err
		}
		if err = proc.Start(); err != nil {
			return err
		}
	}

	// print JSON
	fmt.Fprintln(outputDevice, string(data))

	// close postprocessor if required
	switch proc {
	case nil:
		return nil
	default:
		if err = outputDevice.Close(); err != nil {
			return err
		}
		return proc.Wait()
	}
}

type wNopCloser struct {
	io.Writer
}

func (wNopCloser) Close() error { return nil }

func wrapWNopCloser(w io.Writer) io.WriteCloser {
	return wNopCloser{w}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
