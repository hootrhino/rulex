package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func post(data interface{}, url string) (string, error) {
	p, err1 := json.Marshal(data)
	if err1 != nil {
		return "", err1
	}
	r, err2 := http.Post(url, "application/json", bytes.NewBuffer(p))
	if err2 != nil {
		return "", err2
	}
	defer r.Body.Close()

	body, err3 := ioutil.ReadAll(r.Body)
	if err3 != nil {
		return "", err3
	}
	return string(body), nil
}

type HTTPTarget struct {
	XStatus
}

func (ht *HTTPTarget) Register(outEndId string) {
	ht.PointId = outEndId
}
func (ht *HTTPTarget) Start() {
	fmt.Println("OK")
}
func (ht *HTTPTarget) Test(outEndId string) bool {
	return true
}
func (ht *HTTPTarget) Enabled() bool {
	return ht.Enable
}
func (ht *HTTPTarget) Reload() {

}
func (ht *HTTPTarget) Pause() {

}
func (ht *HTTPTarget) Status() State {
	return UP

}
func (ht *HTTPTarget) To(data interface{}) error {
	config := ht.RuleEngine.GetOutEnd(ht.PointId).Config
	_, err := post(data, (*config)["url"].(string))
	return err
}

//
func (ht *HTTPTarget) Stop() {

}
