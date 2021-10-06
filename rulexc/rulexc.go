package main

//
// Cli client of rulex
//
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ngaut/log"
	"github.com/urfave/cli/v2"
)

//HTTP Post
func post(data map[string]interface{}, host string, api string) (int, string) {
	p, errs1 := json.Marshal(data)
	if errs1 != nil {
		log.Fatal(errs1)
	}
	r, errs2 := http.Post("http://"+host+":2580/api/v1/"+api,
		"application/json",
		bytes.NewBuffer(p))
	if errs2 != nil {
		log.Fatal(errs2)
	}
	defer r.Body.Close()

	body, errs5 := ioutil.ReadAll(r.Body)
	if errs5 != nil {
		log.Fatal(errs5)
	}
	return r.StatusCode, string(body)
}

// HTTP Get
func get(host string, api string) string {
	// Get list
	r, errs := http.Get(("http://" + host + ":2580/api/v1/" + api))
	if errs != nil {
		log.Fatal(errs)
	}
	defer r.Body.Close()
	body, errs2 := ioutil.ReadAll(r.Body)
	if errs2 != nil {
		log.Fatal(errs2)
	}
	return string(body)
}

// main
func main() {
	app := &cli.App{
		Action: func(c *cli.Context) error {
			return nil
		},
		Commands: []*cli.Command{
			// SystemInfo
			{
				Name:  "system-info",
				Usage: "system-info",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "host",
						Usage: "Host of rulex",
						Value: "127.0.0.1",
					},
				},
				Action: func(c *cli.Context) error {
					host := c.String("host")
					result := get(host, "system")
					fmt.Println(result)
					return nil
				},
			},
			// List all inends
			{
				Name:  "inend-list",
				Usage: "inend-list",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "host",
						Usage: "Host of rulex",
						Value: "127.0.0.1",
					},
				},
				Action: func(c *cli.Context) error {
					host := c.String("host")
					result := get(host, "inends")
					fmt.Println(result)
					return nil
				},
			},
			// List all outends
			{
				Name:  "outend-list",
				Usage: "outend-list",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "host",
						Usage: "Host of rulex",
						Value: "127.0.0.1",
					},
				},
				Action: func(c *cli.Context) error {
					host := c.String("host")
					result := get(host, "outends")
					fmt.Println(result)
					return nil
				},
			},
			// List all rules
			{
				Name:  "rules-list",
				Usage: "rules-list",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "host",
						Usage: "Host of rulex",
						Value: "127.0.0.1",
					},
				},
				Action: func(c *cli.Context) error {
					host := c.String("host")
					result := get(host, "rules")
					fmt.Println(result)
					return nil
				},
			},
			// Create InEnd
			{
				Name:  "inend-create",
				Usage: "inend-create",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "host",
						Usage: "Host of rulex",
						Value: "127.0.0.1",
					},
					&cli.StringFlag{
						Name:     "config",
						Usage:    "Config of rulex",
						Value:    "",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					host := c.String("host")
					config := c.String("config")
					maps := map[string]interface{}{}
					err := json.Unmarshal([]byte(config), &maps)
					if err != nil {
						log.Error(config, err)
					} else {
						_, result := post(maps, host, "inends")
						fmt.Println(result)
					}
					return nil
				},
			},
			// create outend
			{
				Name:  "outend-create",
				Usage: "outend-create",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "host",
						Usage: "Host of rulex",
						Value: "127.0.0.1",
					},
					&cli.StringFlag{
						Name:     "config",
						Usage:    "Config of rulex",
						Value:    "",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					host := c.String("host")
					config := c.String("config")
					maps := map[string]interface{}{}
					err := json.Unmarshal([]byte(config), &maps)
					if err != nil {
						log.Error(config, err)
					} else {
						_, result := post(maps, host, "outends")
						fmt.Println(result)
					}
					return nil
				},
			},
			// create rule
			{
				Name:  "rule-create",
				Usage: "rule-create",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "host",
						Usage: "Host of rulex",
						Value: "127.0.0.1",
					},
					&cli.StringFlag{
						Name:     "config",
						Usage:    "Config of rulex",
						Value:    "",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					host := c.String("host")
					config := c.String("config")
					maps := map[string]interface{}{}
					err := json.Unmarshal([]byte(config), &maps)
					if err != nil {
						log.Error(config, err)
					} else {
						_, result := post(maps, host, "rules")
						fmt.Println(result)
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Info("Finished.")
	}
}
