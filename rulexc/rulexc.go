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
			log.Debug(c.String("host"))
			log.Debug(c.String("username"))
			log.Debug(c.String("password"))
			return nil
		},
		Commands: []*cli.Command{
			// Auth
			{
				Name:  "auth",
				Usage: "auth",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "host",
						Usage: "Host of rulex",
						Value: "127.0.0.1",
					},
					&cli.StringFlag{
						Name:     "username",
						Required: true,
						Usage:    "Username of rulex",
					},
					&cli.StringFlag{
						Name:     "password",
						Required: true,
						Usage:    "Password of rulex",
					},
				},
				Action: func(c *cli.Context) error {
					host := c.String("host")
					username := c.String("username")
					password := c.String("password")
					_, result := post(map[string]interface{}{
						"username": username,
						"password": password,
					}, host, "auth")
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
			// Create InEnd
			//go run .\rulexc.go inend-create --config  '{\"name\":\"test\",\"type\":\"MQTT\",\"config\":{\"server\":\"127.0.0.1\",\"port\":1883,\"username\":\"test\",\"password\":\"test\",\"clientId\":\"test\"},\"description\":\"Description\"}'
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
						log.Error(err)
					} else {
						_, result := post(maps, host, "inends")
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
	}
}
