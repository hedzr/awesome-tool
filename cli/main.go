/*
 * Copyright © 2019 Hedzr Yeh.
 */

package main

import (
	awesome_tool "github.com/hedzr/awesome-tool"
	"github.com/hedzr/awesome-tool/ags"
	"github.com/hedzr/cmdr"
	"github.com/sirupsen/logrus"
)

func main() {
	Entry()
}

// Entry is app main entry
func Entry() {

	// logrus.SetLevel(logrus.DebugLevel)
	// logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	// logex.EnableWith(logrus.DebugLevel)

	if err := cmdr.Exec(buildRootCmd(),
		// To disable internal commands and flags, uncomment the following codes
		// cmdr.WithBuiltinCommands(false, false, false, false, false),
		// daemon.WithDaemon(svr.NewDaemon(), nil, nil, nil),
		// cmdr.WithHelpTabStop(40),
		cmdr.WithLogex(logrus.DebugLevel),
		cmdr.WithWatchMainConfigFileToo(true),
	); err != nil {
		logrus.Errorf("Error: %v", err)
	}

}

func buildRootCmd() (rootCmd *cmdr.RootCommand) {

	// To disable internal commands and flags, uncomment the following codes
	// cmdr.EnableVersionCommands = false
	// cmdr.EnableVerboseCommands = false
	// cmdr.EnableCmdrCommands = false
	// cmdr.EnableHelpCommands = false
	// cmdr.EnableGenerateCommands = false

	// daemon.Enable(server.NewDaemon(), modifier, onAppStart, onAppExit)

	// root

	root := cmdr.Root(awesome_tool.AppName, awesome_tool.Version).
		// Header("cmdr-http2 - An HTTP2 server - no version - hedzr").
		Copyright(copyright, "Hedzr").
		Description(desc, longDesc).
		Examples(examples)
	rootCmd = root.RootCommand()

	// build

	buildCmd := root.NewSubCommand().
		Titles("b", "build").
		Description("building operations...", ``)

	// attachConsulConnectFlags(buildCmd)

	boCmd := buildCmd.NewSubCommand().
		Titles("1", "one").
		Description("Build a repos' stars page for an awesome-list", ``).
		Examples(examplesBuildOne).
		Action(buildOne)
	boCmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("w", "work-dir", "wdir", "wd").
		Description("working directory", ``).
		DefaultValue("./output", "DIR")
	boCmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("s", "source", "src", "src-url").
		Description("the source awesomeness repo url", ``).
		DefaultValue("https://github.com/avelino/awesome-go", "URL")
	boCmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("n", "name").
		Description("main name", ``).
		DefaultValue("", "NAME")

	boCmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("1", "first-loop", "1st").
		Description("stop at end of first loop", ``).
		DefaultValue(false, "")
	boCmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("2", "2nd-loop", "2nd").
		Description("stop at end of second loop", ``).
		DefaultValue(false, "")
	boCmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("5", "5th-loop", "5th").
		Description("stop at end of fifth loop", ``).
		DefaultValue(false, "")

	return
}

func buildOne(cmd *cmdr.Command, args []string) (err error) {
	err = ags.Main()
	return
}

func msTagsModify(cmd *cmdr.Command, args []string) (err error) {
	// err = consul.Tags()
	return
}

const (
	// appName   = "awesome-tool"
	copyright = "awesome-tool is a tool for github awesome topics"
	desc      = "awesome-tool is a tool. It builds the page to trace the repo statistics."
	longDesc  = "awesome-tool is a tool. It builds the page to trace the repo statistics."
	examples  = `
$ {{.AppName}} gen shell [--bash|--zsh|--auto]
  generate bash/shell completion scripts
$ {{.AppName}} gen man
  generate linux man page 1
$ {{.AppName}} --help
  show help screen.
`
	examplesBuildOne = `
$ {{.AppName}} build one --name=awesome-go --source=https://github.com/avelino/awesome-go --work-dir=./output
  build the repo stars list for awesome-go list. 'name' is optional.
`
	overview = ``
)
