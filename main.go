package main

import (
	"fmt"
	"runtime"

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
//go:generate ./gen_banner.sh
//
func main() {

	//--------------------------------------
	app := &cli.App{
		Name:  "RULEX, a lightweight iot data rule gateway",
		Usage: "http://rulex.ezlinker.cn",
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Start rulex",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "db",
						Usage: "Database of rulex",
						Value: "rulex.db",
					},
				},
				Action: func(c *cli.Context) error {
					utils.ShowBanner()
					log.Info("Load config db:", c.String("db"))
					engine.RunRulex(c.String("db"))
					log.Debug("Run rulex successfully.")
					return nil
				},
			},
			// version
			{
				Name:  "version",
				Usage: "Rulex version",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "version",
						Usage: "Rulex version",
					},
				},
				Action: func(c *cli.Context) error {
					version := fmt.Sprintf("[%v-%v-%v]", runtime.GOOS, runtime.GOARCH, typex.DefaultVersion.Version)
					fmt.Println("|> Current Version is: " + version)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
