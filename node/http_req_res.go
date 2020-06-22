package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func writeErrRes(w http.ResponseWriter, err error) {
	jsonErrRes, _ := json.Marshal(ErrRes{err.Error()})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write(jsonErrRes)
}

func writeRes(w http.ResponseWriter, content interface{}) {
	contentJson, err := json.Marshal(content)
	if err != nil {
		writeErrRes(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(contentJson)
}

func readReq(r *http.Request, reqBody interface{}) error {
	reqBodyJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("unable to read request body. %s", err.Error())
	}
	defer r.Body.Close()

	err = json.Unmarshal(reqBodyJson, reqBody)
	if err != nil {
		return fmt.Errorf("unable to unmarshal request body. %s", err.Error())
	}

	return nil
}

func readRes(r *http.Response, reqBody interface{}) error {
	reqBodyJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("unable to read response body. %s", err.Error())
	}
	defer r.Body.Close()

	err = json.Unmarshal(reqBodyJson, reqBody)
	if err != nil {
		return fmt.Errorf("unable to unmarshal response body. %s", err.Error())
	}

	return nil
}