package main

import (
	"context"
	"fmt"
	"github.com/hootrhino/rulex/engine"
	"github.com/hootrhino/rulex/ossupport"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	"github.com/urfave/cli/v2"
	"log"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"
)

func init() {
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

func main() {
	app := &cli.App{
		Name:  "RULEX FrameWork",
		Usage: "Homepage: http://www.hootrhino.com",
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
					log.Printf("[Prepare Stage] Run rulex successfully.")
					return nil
				},
			},
			{
				Name:  "upgrade",
				Usage: "! JUST FOR Upgrade FirmWare",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "oldpid",
						Usage: "! THIS PARAMENT IS JUST FOR Upgrade FirmWare",
						Value: -1,
					},
				},
				Action: func(c *cli.Context) error {
					OldPid := c.Int("oldpid")
					log.Println("[Prepare Stage] Updater Pid=",
						os.Getpid(), "Gid=", os.Getegid(), " OldPid:", OldPid)
					if OldPid < 0 {
						log.Printf("[Prepare Stage] Invalid OldPid:%d", OldPid)
						return nil
					}
					// Try 5 times
					killOld := true
					log.Println("[Prepare Stage] Try to kill Old Process:", OldPid)
					if killOld {
						// EEKITH3 Use SystemCtl manage RULEX
						env := os.Getenv("ARCHSUPPORT")
						if runtime.GOOS == "linux" {
							log.Println("[Prepare Stage] Ready to Upgrade on product:", env)
							if err := ossupport.UnzipFirmware(
								"/usr/local/upload/Firmware/Firmware.zip",
								"/usr/local"); err != nil {
								log.Println("[Prepare Stage] Unzip error:", err)
								return err
							}
							if err := ossupport.Restart(); err != nil {
								log.Println("[Prepare Stage] Restart rulex error", err)
								return nil
							}
							log.Println("[Prepare Stage] Restart rulex success, Upgrade Process Exited")
							os.Exit(0)
							return nil
						}
					}
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
					version := fmt.Sprintf("[%v-%v-%v]",
						runtime.GOOS, runtime.GOARCH, typex.DefaultVersion.Version)
					fmt.Println("[*] Current Version: " + version)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
