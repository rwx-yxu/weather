package weather

import (
	"fmt"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
	"github.com/rwxrob/vars"
)

func init() {
	Z.Vars.SoftInit()
}

var Cmd = &Z.Cmd{

	Name:      `weather`,
	Summary:   `a command that prints out current day routine schedule`,
	Version:   `v0.4.1`,
	Copyright: `Copyright 2023 Yongle Xu`,
	License:   `Apache-2.0`,
	Site:      `yonglexu.dev`,
	Source:    `git@github.com:rwx-yxu/routine.git`,
	Issues:    `github.com/rwx-yxu/routine/issues`,

	Commands: []*Z.Cmd{
		nowCmd,
		// standard external branch imports (see rwxrob/{help,conf,vars})
		help.Cmd, vars.Cmd,
	},

	// Add custom BonzaiMark template extensions (or overwrite existing ones).

	Description: `
		{{cmd .Name}} is a tool that queries the Open weather map API for current weather information.
			`,
}

var nowCmd = &Z.Cmd{
	Name:     `now`,
	Summary:  `print current weather conditions to standard output (default)`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(x *Z.Cmd, _ ...string) error {
		fmt.Println("Hello world")
		return nil
	},
}
