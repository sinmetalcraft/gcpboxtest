package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func StoragePubSubNotifyHandler(w http.ResponseWriter, r *http.Request) {
	h, err := json.Marshal(r.Header)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(h))

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}
