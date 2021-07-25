package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func post(data interface{}, url string) (string, error) {
	p, errs1 := json.Marshal(data)
	if errs1 != nil {
		return "", errs1
	}
	r, errs2 := http.Post(url, "application/json", bytes.NewBuffer(p))
	if errs2 != nil {
		return "", errs2
	}
	defer r.Body.Close()

	body, errs3 := ioutil.ReadAll(r.Body)
	if errs3 != nil {
		return "", errs3
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
	config := ht.ruleEngine.GetOutEnd(ht.PointId).Config
	_, err := post(data, (*config)["url"].(string))
	return err
}
func (ht *HTTPTarget) Stop() {

}
