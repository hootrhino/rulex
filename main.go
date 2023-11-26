package main

import (
	"context"
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"

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
					if runtime.GOOS != "linux" {
						fmt.Println("[RULEX UPGRADE] Only Support Linux")
						return nil
					}
					file, err := os.Create("./local-upgrade-log.txt")
					if err != nil {
						fmt.Println(err)
						return nil
					}
					defer file.Close()
					os.Stdout = file
					os.Stderr = file
					if !c.Bool("upgrade") {
						fmt.Println("[RULEX UPGRADE] Nothing todo")
						return nil
					}
					// unzip Firmware
					if err := ossupport.UnzipFirmware(
						"/usr/local/upload/Firmware/Firmware.zip",
						"/usr/local"); err != nil {
						fmt.Println("[RULEX UPGRADE] Unzip error:", err)
						return nil
					}
					if err := ossupport.RestartRulex(); err != nil {
						fmt.Println("[RULEX UPGRADE] Restart rulex error", err)
						return nil
					}
					// Remove old package
					if err := os.Remove("/usr/local/upload/Firmware/Firmware.zip"); err != nil {
						fmt.Println("[RULEX UPGRADE] Restart rulex error", err)
						return nil
					}
					fmt.Println("[RULEX UPGRADE] Restart rulex success, Upgrade Process Exited")
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
					if runtime.GOOS != "linux" {
						fmt.Println("[DATA RECOVER] Only Support Linux")
						return nil
					}
					file, err := os.Create("./rulex-recover-log.txt")
					if err != nil {
						fmt.Println(err)
						return nil
					}
					defer file.Close()
					os.Stdout = file
					os.Stderr = file
					if !c.Bool("recover") {
						fmt.Println("[DATA RECOVER] Nothing todo")
						return nil
					}
					fmt.Println("[DATA RECOVER] Remove Old Db File")
					if err := os.Remove("./rulex.db"); err != nil {
						fmt.Println("[DATA RECOVER] Remove Old Db File error:", err)
						return nil
					}
					fmt.Println("[DATA RECOVER] Remove Old Db File Finished")
					fmt.Println("[DATA RECOVER] Move New Db File")
					recoveryDb := "./upload/Backup/recovery.db"
					if err := ossupport.MoveFile(recoveryDb, "./rulex.db"); err != nil {
						fmt.Println("[DATA RECOVER] Move New Db File error", err)
						return nil
					}
					fmt.Println("[DATA RECOVER] Move New Db File Finished")
					fmt.Println("[DATA RECOVER] Try to Restart rulex")
					// upgrade lock
					if err := os.WriteFile("/var/run/rulex-upgrade.lock", []byte{48}, 0755); err != nil {
						fmt.Println("[DATA RECOVER] Write Recover Lock File error:", err)
						return nil
					}
					if err := ossupport.RestartRulex(); err != nil {
						fmt.Println("[DATA RECOVER] Restart rulex error", err)
					} else {
						fmt.Println("[DATA RECOVER] Restart rulex success, Recover Process Exited")
					}
					// upgrade lock
					if err := os.Remove("/var/run/rulex-upgrade.lock"); err != nil {
						fmt.Println("[DATA RECOVER] Remove Recover Lock File error:", err)
						return nil
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
					fmt.Printf("[LICENCE ACTIVE]: Admin(%s,%s), mac addr:[%s] try to request license from %s\n",
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
					fmt.Println("[*] Rulex Version: " + version)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
