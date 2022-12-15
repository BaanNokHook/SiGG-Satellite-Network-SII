// SiGG-Satellite-Network-SII  //

package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	_ "go.uber.org/automaxprocs"
)

// version will be initialized when building
var version = "latest"

func main() {
	app := cli.NewApp()
	app.Name = "SkyWalking-Satellite"
	app.Version = version
	app.Compiled = time.Now()
	app.Usage = "Satellite is for collecting APM data."
	app.Description = "A lightweight collector/sidecar could be deployed closing to the target monitored system, to collect metrics, traces, and logs."
	app.Commands = []*cli.Command{
		&cmdStart,
		&cmdDocs,
	}
	app.Action = cli.ShowAppHelp
	if err := app.Run(os.Args); err != nil {
		log.Fatalln("start SkyWalking Satellite fail", err)
	}
}
