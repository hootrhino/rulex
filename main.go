package main

import (
	"context"
	"fmt"
	"github.com/hootrhino/rulex/docs"
	"log"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/hootrhino/rulex/engine"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

func init() {
	info := docs.SwaggerInfo
	fmt.Println(info)
	if runtime.GOOS == "linux" {
		uid := os.Getuid()
		if uid != 0 {
			fmt.Println("Run failed, Rulex MUST run with ROOT permission!")
			os.Exit(1)
			return
		}
	}

	go func() {
		for {
			select {
			case <-context.Background().Done():
				return
			default:
				time.Sleep(30 * time.Second)
				runtime.GC()
			}
		}
	}()
	dist, err := utils.GetOSDistribution()
	if err != nil {
		panic(err)
	}
	typex.DefaultVersion.Dist = dist
	arch := fmt.Sprintf("%s-%s", typex.DefaultVersion.Dist, runtime.GOARCH)
	typex.DefaultVersion.Arch = arch
}

/*
*
* ！！！注意：这个 main 函数仅仅是用来做启动测试用，并非真正的应用，具体的应用需要开发者自己去开发。
* 详情需要关注：http://www.hootrhino.com
*
 */
// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {

	//--------------------------------------
	app := &cli.App{
		Name:  "RULEX FrameWork",
		Usage: "Goto Document: http://www.hootrhino.com",
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
					fmt.Println(typex.Banner)
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
