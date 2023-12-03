package main

import (
	"context"
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"

	_ "github.com/hootrhino/rulex/component/cron_task/docs"
	"github.com/hootrhino/rulex/engine"
	"github.com/hootrhino/rulex/ossupport"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	"github.com/urfave/cli/v2"
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
	env := os.Getenv("ARCHSUPPORT")
	if env == "EEKITT507" {
		typex.DefaultVersion.Product = env
	}
	if env == "EEKITH3" {
		typex.DefaultVersion.Product = env
	}
	if env == "WKYS805" {
		typex.DefaultVersion.Product = env
	}
	if env == "RPI4B" {
		typex.DefaultVersion.Product = env
	}
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
				Usage: "Start rulex, Must with config: -config=path/rulex.ini",
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
					utils.CLog(typex.Banner)
					engine.RunRulex(c.String("config"))
					fmt.Printf("[RULEX UPGRADE] Run rulex successfully.")
					return nil
				},
			},
			{
				Name:   "upgrade",
				Hidden: true,
				Usage:  "! JUST FOR Upgrade FirmWare",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "upgrade",
						Usage: "! THIS PARAMENT IS JUST FOR Upgrade FirmWare",
						Value: false,
					},
				},
				Action: func(c *cli.Context) error {
					file, err := os.Create(ossupport.UpgradeLogPath)
					if err != nil {
						utils.CLog(err.Error())
						return nil
					}
					defer file.Close()
					os.Stdout = file
					os.Stderr = file
					// upgrade lock
					if err := os.WriteFile(ossupport.UpgradeLockPath, []byte{48}, 0755); err != nil {
						utils.CLog("[RULEX UPGRADE] Write Upgrade Lock File error:%s", err.Error())
						return nil
					}
					defer func() {
						// upgrade lock
						if err := os.Remove(ossupport.UpgradeLockPath); err != nil {
							utils.CLog("[RULEX UPGRADE] Remove Upgrade Lock File error:%s", err.Error())
							return
						}
						utils.CLog("[RULEX UPGRADE] Remove Upgrade Lock File Finished")
					}()
					if runtime.GOOS != "linux" {
						utils.CLog("[RULEX UPGRADE] Only Support Linux")
						return nil
					}
					if !c.Bool("upgrade") {
						utils.CLog("[RULEX UPGRADE] Nothing todo")
						return nil
					}
					// unzip Firmware
					utils.CLog("[RULEX UPGRADE] Unzip Firmware")
					if err := ossupport.UnzipFirmware(
						ossupport.FirmwarePath, ossupport.MainWorkDir); err != nil {
						utils.CLog("[RULEX UPGRADE] Unzip error:%s", err.Error())
						return nil
					}
					utils.CLog("[RULEX UPGRADE] Unzip Firmware finished")
					// Remove old package
					utils.CLog("[RULEX UPGRADE] Remove Firmware")
					if err := os.Remove(ossupport.FirmwarePath); err != nil {
						utils.CLog("[RULEX UPGRADE] Remove Firmware error:%s", err.Error())
						return nil
					}
					utils.CLog("[RULEX UPGRADE] Remove Firmware finished")
					//
					utils.CLog("[RULEX UPGRADE] Restart rulex")
					if err := ossupport.RestartRulex(); err != nil {
						utils.CLog("[RULEX UPGRADE] Restart rulex error:%s", err.Error())
						return nil
					}
					utils.CLog("[RULEX UPGRADE] Restart rulex finished, Upgrade Process Exited")
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
					file, err := os.Create(ossupport.RecoverLogPath)
					if err != nil {
						utils.CLog(err.Error())
						return nil
					}
					defer file.Close()
					os.Stdout = file
					os.Stderr = file
					// upgrade lock
					if err := os.WriteFile(ossupport.UpgradeLockPath, []byte{48}, 0755); err != nil {
						utils.CLog("[DATA RECOVER] Write Recover Lock File error:%s", err.Error())
						return nil
					}
					defer func() {
						// upgrade lock
						if err := os.Remove(ossupport.UpgradeLockPath); err != nil {
							utils.CLog("[DATA RECOVER] Remove Recover Lock File error:%s", err.Error())
							return
						}
						utils.CLog("[DATA RECOVER] Remove Recover Lock File Finished")
					}()
					if runtime.GOOS != "linux" {
						utils.CLog("[DATA RECOVER] Only Support Linux")
						return nil
					}

					if !c.Bool("recover") {
						utils.CLog("[DATA RECOVER] Nothing todo")
						return nil
					}
					utils.CLog("[DATA RECOVER] Remove Old Db File")
					if err := os.Remove(ossupport.RunDbPath); err != nil {
						utils.CLog("[DATA RECOVER] Remove Old Db File error:%s", err.Error())
						return nil
					}
					utils.CLog("[DATA RECOVER] Remove Old Db File Finished")
					utils.CLog("[DATA RECOVER] Move New Db File")
					if err := ossupport.MoveFile(ossupport.RecoveryDbPath,
						ossupport.RunDbPath); err != nil {
						utils.CLog("[DATA RECOVER] Move New Db File error:%s", err.Error())
						return nil
					}
					utils.CLog("[DATA RECOVER] Move New Db File Finished")
					utils.CLog("[DATA RECOVER] Try to Restart rulex")
					if err := ossupport.RestartRulex(); err != nil {
						utils.CLog("[DATA RECOVER] Restart rulex error:%s", err.Error())
					} else {
						utils.CLog("[DATA RECOVER] Restart rulex success, Recover Process Exited")
					}
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
					utils.CLog("[LICENCE ACTIVE]: Admin(%s,%s), mac addr:[%s] try to request license from %s\n",
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
					utils.CLog("[*] Rulex Version: " + version)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
