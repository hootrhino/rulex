package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/ngaut/log"
)

//
// MakeUUID
//
func MakeUUID(prefix string) string {
	return prefix + "_" + uuid.NewString()
}

//
//
//

func post(data map[string]interface{}, api string) string {
	p, errs1 := json.Marshal(data)
	if errs1 != nil {
		log.Error(errs1)
	}
	r, errs2 := http.Post(api, "application/json",
		bytes.NewBuffer(p))
	if errs2 != nil {
		log.Error(errs2)
	}
	defer r.Body.Close()

	body, errs5 := ioutil.ReadAll(r.Body)
	if errs5 != nil {
		log.Error(errs5)
	}
	return string(body)
}
