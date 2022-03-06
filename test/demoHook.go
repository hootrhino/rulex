package test

import "log"

type demoHook struct {
}

func (demo *demoHook) Work(data string) error {
	log.Println(data)
	return nil
}

func (demo *demoHook) Error(err error) {
	log.Println(err.Error())
}

func (demo *demoHook) Name() string {
	return "Demo Hook"
}
