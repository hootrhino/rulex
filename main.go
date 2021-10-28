package main

import (
	"fmt"

	"github.com/ngaut/log"
	"github.com/urfave/cli/v2"

	_ "net/http/pprof"
	"os"
	"rulex/engine"
	"rulex/typex"
	"rulex/utils"
)

//
//go:generate ./gen_version.sh
//go:generate ./gen_proto.sh
//
func main() {
	//--------------------------------------
	app := &cli.App{
		Name:  "RULEX, a lightweight iot data rule gateway",
		Usage: "http://rulex.ezlinker.cn",
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "rulex run [path of 'rulex.db']",
				Action: func(c *cli.Context) error {
					utils.ShowBanner()
					if c.Args().Len() > 0 {
						log.Info("Use config db:", c.Args().Get(0))
						engine.RunRulex(c.Args().Get(0))
					} else {
						engine.RunRulex("rulex.db")
					}
					log.Debug("Run rulex successfully.")
					return nil
				},
			},
			// version
			{
				Name:  "version",
				Usage: "rulex version",
				Action: func(c *cli.Context) error {
					fmt.Println("Current Version is: " + typex.DefaultVersion.Version)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
