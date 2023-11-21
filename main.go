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

//go:generate bash ./gen_info.sh
func main() {
	app := &cli.App{
		Name:  "RULEX Gateway FrameWork",
		Usage: "Homepage: http://rulex.hootrhino.com",
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Start rulex, Must with config: -config path/rulex.ini",
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
					log.Printf("[RULEX UPGRADE] Run rulex successfully.")
					return nil
				},
			},
			{
				Name:   "upgrade",
				Hidden: true,
				Usage:  "! JUST FOR Upgrade FirmWare",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "oldpid",
						Usage: "! THIS PARAMENT IS JUST FOR Upgrade FirmWare",
						Value: -1,
					},
				},
				Action: func(c *cli.Context) error {
					OldPid := c.Int("oldpid")
					log.Println("[RULEX UPGRADE] Updater Pid=",
						os.Getpid(), "Gid=", os.Getegid(), " OldPid:", OldPid)
					if OldPid < 0 {
						log.Printf("[RULEX UPGRADE] Invalid OldPid:%d", OldPid)
						return nil
					}
					// Try 5 times
					killOld := true
					log.Println("[RULEX UPGRADE] Try to kill Old Process:", OldPid)
					if killOld {
						// EEKITH3 Use SystemCtl manage RULEX
						env := os.Getenv("ARCHSUPPORT")
						if runtime.GOOS == "linux" {
							log.Println("[RULEX UPGRADE] Ready to Upgrade on product:", env)
							if err := ossupport.UnzipFirmware(
								"/usr/local/upload/Firmware/Firmware.zip",
								"/usr/local"); err != nil {
								log.Println("[RULEX UPGRADE] Unzip error:", err)
								return err
							}
							if err := ossupport.Restart(); err != nil {
								log.Println("[RULEX UPGRADE] Restart rulex error", err)
								return err
							}
							log.Println("[RULEX UPGRADE] Restart rulex success, Upgrade Process Exited")
						}
					}
					os.Exit(0)
					return nil
				},
			},
			// 数据恢复
			{
				Name:   "recover",
				Usage:  "! JUST FOR Recover Data",
				Hidden: true,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "recover",
						Usage: "! THIS PARAMENT IS JUST FOR Recover Data",
						Value: false,
					},
				},
				Action: func(c *cli.Context) error {
					if !c.Bool("recover") {
						return fmt.Errorf("[DATA RECOVER] Nothing todo")
					}
					if err := ossupport.StopRulex(); err != nil {
						log.Println("[DATA RECOVER] Stop rulex error", err)
						return err
					}
					recoveryDb := "/usr/local/upload/Backup/recovery.db"
					log.Println("[DATA RECOVER] Move Db File")
					if err := ossupport.MoveFile(recoveryDb, "/usr/local/rulex.db"); err != nil {
						log.Println("[DATA RECOVER] Move Db File error", err)
						return err
					}
					log.Println("[DATA RECOVER] Move Db File Finished")
					log.Println("[DATA RECOVER] Try to Restart rulex")
					if err := ossupport.Restart(); err != nil {
						log.Println("[DATA RECOVER] Restart rulex error", err)
						return err
					}
					log.Println("[DATA RECOVER] Restart rulex success, Recover Process Exited")
					os.Exit(0)
					return nil
				},
			},
			{
				Name:   "active",
				Usage:  "active -H host -U rhino -P hoot",
				Hidden: true,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "H",
						Usage: "active server ip",
					},
					&cli.StringFlag{
						Name:  "U",
						Usage: "active admin username",
					},
					&cli.StringFlag{
						Name:  "P",
						Usage: "active admin password",
					},
				},

				Action: func(c *cli.Context) error {
					host := c.String("H")
					if host == "" {
						return fmt.Errorf("[LICENCE ACTIVE]: missing host")
					}
					username := c.String("U")
					if username == "" {
						return fmt.Errorf("[LICENCE ACTIVE]: missing admin username")
					}
					password := c.String("P")
					if password == "" {
						return fmt.Errorf("[LICENCE ACTIVE]: missing admin password")
					}
					macAddr, err := ossupport.ReadIfaceMacAddr("eth0")
					if err != nil {
						return err
					}
					// commercial version will implement it
					log.Printf("[LICENCE ACTIVE]: Admin(%s,%s), mac addr:[%s] try to request license from %s\n",
						username, password, macAddr, host)
					return nil
				},
			},
			// version
			{
				Name:  "version",
				Usage: "Show Rulex Current Version",
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
