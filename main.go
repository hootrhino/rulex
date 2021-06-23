package main

import (
	"io/ioutil"
	"os"
	"os/signal"
	"rulenginex/x"
	"syscall"

	"github.com/ngaut/log"
	lua "github.com/yuin/gopher-lua"
)

func main() {
	c := make(chan os.Signal, 1)
	// 监听信号
	signal.Notify(c, syscall.SIGQUIT)
	ruleEngine := x.RuleEngine{}
	binds := make(map[string]x.Rule)
	config := map[string]interface{}{
		"server":   "127.0.0.1",
		"port":     1883,
		"username": "test",
		"password": "test",
		"clientId": "test",
	}
	in1 := x.InEnd{
		Id:          x.MakeUUID("INEND"),
		Type:        "MQTT",
		Name:        "MQTT Stream",
		Description: "MQTT Input Stream",
		Binds:       &binds,
		Config:      &config,
	}
	ruleEngine.LoadInEnds(&in1)
	actions := `
		Actions = {
			function(data)
			    print("[LUA Actions Callback]: Mqtt payload:", data)
			    return true, data
		    end
	}`
	from := []string{in1.Id}
	failed := `function Failed(error) print("[LUA Callback] call failed from lua:", error) end`
	success := `function Success() print("[LUA Callback] call success from lua") end`
	rule1 := x.Rule{
		Id:          x.MakeUUID("RULE"),
		Name:        "just_a_test_rule",
		Description: "just_a_test_rule",
		From:        from,
		Actions:     actions,
		Success:     success,
		Failed:      failed,
		VM:          lua.NewState(),
	}

	//
	if ruleEngine.LoadRules(&rule1) != nil {
		log.Fatal("Rule load failed")
	}
	defaultBanner :=
		`
	--------------------------------------------------------------------
	____  _   _ _     _____ _   _  ____ ___ _   _ _____          __  __
	|  _ \| | | | |   | ____| \ | |/ ___|_ _| \ | | ____|         \ \/ /
	| |_) | | | | |   |  _| |  \| | |  _ | ||  \| |  _|    _____   \  / 
	|  _ <| |_| | |___| |___| |\  | |_| || || |\  | |___  |_____|  /  \ 
	|_| \_\\___/|_____|_____|_| \_|\____|___|_| \_|_____|         /_/\_\
	---------------------------------------------------------------------
	`
	ruleEngine.Start(func() {
		file, err := os.Open("conf/banner.txt")
		if err != nil {
			log.Warn("No banner found, print default banner")
			log.Info(defaultBanner)
		} else {
			data, err := ioutil.ReadAll(file)
			if err != nil {
				log.Warn("No banner found, print default banner")
				log.Info(defaultBanner)
			} else {
				log.Info("\n", string(data))
			}
		}
		log.Info("RulengineX start successfully")
		file.Close()
	})
	<-c
	os.Exit(0)
}
