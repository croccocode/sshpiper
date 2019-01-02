package main

import (
	"fmt"
	"github.com/tg123/sshpiper/sshpiperd/upstream"
)

func createPipeMgr(driver *string) interface{} {

	load := func() (upstream.Provider, error) {
		if *driver == "" {
			return nil, fmt.Errorf("must provider upstream driver")
		}

		return upstream.Get(*driver).(upstream.Provider), nil
	}

	// pipe management
	pipeMgrCmd := struct {
		List struct {
			subCommand
		} `command:"list" description:"list all pipes"`
		Add struct {
			subCommand

			PiperUserName string `long:"piper-username" description:"" required:"true" no-ini:"true"`
			// PiperAuthorizedKeysFile flags.Filename

			UpstreamUserName string `long:"upstream-username" description:"mapped user name" no-ini:"true"`
			UpstreamHost     string `long:"host" description:"upstream sshd host" required:"true" no-ini:"true"`
			UpstreamPort     uint   `long:"port" description:"upstream sshd port" default:"22" no-ini:"true"`
			// UpstreamKeyFile  flags.Filename

			// UpstreamHostKey
			// MapType

		} `command:"add" description:"add a pipe to current upstream"`
		Remove struct {
			subCommand

			Name string `long:"name" required:"true" no-ini:"true"`
		} `command:"remove" description:"remove a pipe from current upstream"`
	}{}

	pipeMgrCmd.List.callback = func(args []string) error {
		return nil
	}

	pipeMgrCmd.Add.callback = func(args []string) error {
		p, err := load()
		if err != nil {
			return err
		}

		opt := pipeMgrCmd.Add

		return p.CreatePipe(upstream.CreatePipeOption{
			Username:         opt.PiperUserName,
			UpstreamUsername: opt.UpstreamUserName,
			Host:             opt.UpstreamHost,
			Port:             opt.UpstreamPort,
		})
	}

	pipeMgrCmd.Remove.callback = func(args []string) error {

		name := pipeMgrCmd.Remove.Name
		p, err := load()
		if err != nil {
			return err
		}

		return p.RemovePipe(name)
	}

	return &pipeMgrCmd
}