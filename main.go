package main

/*
*
* ！！！注意：这个 main 函数仅仅是用来做启动测试用，并非真正的应用，具体的应用需要开发者自己去开发。
* 详情需要关注：http://rulex.pages.dev
*
 */
import (
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"runtime"

	"github.com/urfave/cli/v2"

	"github.com/i4de/rulex/engine"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"
)

//
//go:generate ./gen_info.sh
//
func main() {

	//--------------------------------------
	app := &cli.App{
		Name:  "RULEX FrameWork",
		Usage: "Goto Document: http://rulex.pages.dev",
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
					&cli.StringFlag{
						Name:  "config",
						Usage: "Config of rulex",
						Value: "rulex.ini",
					},
				},
				Action: func(c *cli.Context) error {
					utils.ShowBanner()
					engine.RunRulex(c.String("config"))
					glogger.GLogger.Info("Run rulex successfully.")
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
				Action: func(*cli.Context) error {
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
