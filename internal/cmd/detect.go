// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"fmt"
	"io"

	"github.com/caixw/apidoc/v7/build"
	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

var detectFlagSet *flag.FlagSet

var detectRecursive bool

func initDetect() {
	detectFlagSet = command.New("detect", detect, detectUsage)
	detectFlagSet.BoolVar(&detectRecursive, "r", true, locale.Sprintf(locale.FlagDetectRecursive))
}

func detect(io.Writer) error {
	h := core.NewMessageHandler(messageHandle)
	defer h.Stop()

	uri := getPath(detectFlagSet)

	cfg, err := build.DetectConfig(uri, detectRecursive)
	if err != nil {
		return err
	}

	if err = cfg.Save(uri); err != nil {
		return err
	}

	h.Locale(core.Succ, locale.ConfigWriteSuccess, uri)
	return nil
}

func detectUsage(w io.Writer) error {
	_, err := fmt.Fprintln(w, locale.Sprintf(locale.CmdDetectUsage, getFlagSetUsage(detectFlagSet)))
	return err
}
