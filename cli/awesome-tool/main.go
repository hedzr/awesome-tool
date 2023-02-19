/*
 * Copyright © 2019 Hedzr Yeh.
 */

package main

import (
	"fmt"

	awesome_tool "github.com/hedzr/awesome-tool"
	"github.com/hedzr/awesome-tool/ags"
	"github.com/hedzr/cmdr"
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

		// trace.WithTraceEnable(true),
		cmdr.WithXrefBuildingHooks(func(root *cmdr.RootCommand, args []string) {
			cmdr.NewBool(false).
				Titles("trace", "tr").
				Description("enable trace mode for tcp/mqtt send/recv data dump", "").
				// Action(func(cmd *cmdr.Command, args []string) (err error) {
				// 	println("trace mode on")
				// 	cmdr.SetTraceMode(true)
				// 	return
				// }).
				Group(cmdr.SysMgmtGroup).
				EnvKeys("TRACE").
				AttachToRoot(root)
		}, nil),

		cmdr.WithWatchMainConfigFileToo(true),
	); err != nil {
		cmdr.Logger.Errorf("Error: %v", err)
	}

}

func buildRootCmd() (rootCmd *cmdr.RootCommand) {

	// root

	root := cmdr.Root(awesome_tool.AppName, awesome_tool.Version).
		// Header("cmdr-http2 - An HTTP2 server - no version - hedzr").
		Copyright(copyright, "Hedzr").
		Description(desc, longDesc).
		Examples(examples)
	rootCmd = root.RootCommand()

	// info

	cmdr.NewSubCmd().
		Titles("info", "i").
		Description("debug information...", ``).
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			fmt.Printf("'nothing' in config file is: %q\n", cmdr.GetStringR("nothing"))
			fmt.Printf("'z-flag' in config file is: %v\n", cmdr.GetBoolR("z-flag"))
			fmt.Printf("'z-mergable' in config file is: %v\n", cmdr.GetStringSliceR("z-mergeable"))
			fmt.Printf("'z-string' in config file is: %q\n", cmdr.GetStringR("z-string"))
			s1 := cmdr.GetStringSliceR("z-mergeable")
			s1 = append(s1, "e")
			cmdr.Set("z-mergeable", s1)
			return
		}).
		AttachTo(root)

	// build

	buildCmd := cmdr.NewSubCmd().
		Titles("build", "b").
		Description("building operations...", ``).
		AttachTo(root)

	// attachConsulConnectFlags(buildCmd)

	boCmd := cmdr.NewSubCmd().
		Titles("one", "1", "single").
		Description("Build a repos' stars page for an awesome-list", ``).
		Examples(examplesBuildOne).
		Action(buildOne).
		AttachTo(buildCmd)
	cmdr.NewString(`./output`).Placeholder(`DIR`).
		Titles("work-dir", "w", "wdir", "wd").
		Description("working directory", ``).
		AttachTo(boCmd)
	cmdr.NewString(`https://github.com/avelino/awesome-go`).Placeholder(`URL`).
		Titles("source", "s", "src", "src-url").
		Description("the source awesomeness repo url", ``).
		AttachTo(buildCmd)
	cmdr.NewString().Placeholder(`NAME`).
		Titles("name", "n").
		Description("main name", ``).
		AttachTo(buildCmd)

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
