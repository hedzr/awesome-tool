/*
 * Copyright © 2019 Hedzr Yeh.
 */

package main

import (
	"github.com/hedzr/awesome-tool"
	"github.com/hedzr/awesome-tool/ags"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr-addons/pkg/plugins/trace"
	"github.com/hedzr/logex/build"
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
		//cmdr.WithLogex(logrus.DebugLevel),
		cmdr.WithLogx(build.New(cmdr.NewLoggerConfigWith(true, "logrus", "debug"))),
		trace.WithTraceEnable(true),
		cmdr.WithWatchMainConfigFileToo(true),
	); err != nil {
		cmdr.Logger.Errorf("Error: %v", err)
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
		Titles("build", "b").
		Description("building operations...", ``)

	// attachConsulConnectFlags(buildCmd)

	boCmd := buildCmd.NewSubCommand().
		Titles("one", "1", "single").
		Description("Build a repos' stars page for an awesome-list", ``).
		Examples(examplesBuildOne).
		Action(buildOne)
	boCmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("work-dir", "w", "wdir", "wd").
		Description("working directory", ``).
		DefaultValue("./output", "DIR")
	boCmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("source", "s", "src", "src-url").
		Description("the source awesomeness repo url", ``).
		DefaultValue("https://github.com/avelino/awesome-go", "URL")
	boCmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("name", "n").
		Description("main name", ``).
		DefaultValue("", "NAME")

	cmdr.NewInt(-1).
		HeadLike(true, -1, 65535).
		Titles("loops", "l").
		Description("the looping times before stop").
		Group("Debug").
		AttachTo(boCmd)

	//boCmd.NewFlag(cmdr.OptFlagTypeBool).
	//	Titles("first-loop", "1", "1st").
	//	Description("stop at end of first loop", ``).
	//	DefaultValue(false, "").
	//	ToggleGroup("Debug")
	//boCmd.NewFlag(cmdr.OptFlagTypeBool).
	//	Titles("2nd-loop", "2", "2nd").
	//	Description("stop at end of second loop", ``).
	//	DefaultValue(false, "").
	//	ToggleGroup("Debug")
	//boCmd.NewFlag(cmdr.OptFlagTypeBool).
	//	Titles("5th-loop", "5", "5th").
	//	Description("stop at end of fifth loop", ``).
	//	DefaultValue(false, "").
	//	ToggleGroup("Debug")

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
