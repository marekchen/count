package main

import (
	"log"
	"os"
	"runtime"

	"github.com/aisk/chrysanthemum"
	"github.com/fatih/color"
	"github.com/marekchen/count/commands"
)

func run() {
	// disable the log prefix
	log.SetFlags(0)

	commands.Run(os.Args)
}

func init() {
	if runtime.GOOS == "windows" {
		chrysanthemum.Frames = []string{
			"-",
			"\\",
			"|",
			"/",
		}
		chrysanthemum.Success = ">"
		chrysanthemum.Fail = color.RedString("x")
	}
}

func main() {
	run()
	return
}
