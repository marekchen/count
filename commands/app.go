package commands

import (
	"github.com/codegangsta/cli"
	"github.com/marekchen/count/logo"
)

// Run the command line
func Run(args []string) {
	// add banner text to help text
	cli.AppHelpTemplate = logo.Logo() + cli.AppHelpTemplate
	cli.SubcommandHelpTemplate = logo.Logo() + cli.SubcommandHelpTemplate

	app := cli.NewApp()
	app.Name = "count"
	app.Version = "1.0.001"
	app.Action = countAction
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "eachfolder,e",
			Usage: "是否每个文件夹单独计算",
		},
		cli.StringFlag{
			Name:  "path,p",
			Usage: "路径",
			Value: "./",
		},
		cli.StringFlag{
			Name:  "suffix,s",
			Usage: "指定后缀",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "help",
			Aliases:   []string{"h"},
			Usage:     "显示全部命令或者某个子命令的帮助",
			ArgsUsage: "[command]",
			Action: func(c *cli.Context) error {
				args := c.Args()
				if args.Present() {
					return cli.ShowCommandHelp(c, args.First())
				}

				cli.ShowAppHelp(c)
				return nil
			},
		},
	}

	//app.Before = func(c *cli.Context) error {
	//	args := []string{"--_collect-stats"}
	//	args = append(args, c.Args()...)
	//	err := exec.Command(os.Args[0], args...).Start()
	//	_ = err
	//	return nil
	//}

	app.Run(args)
}
