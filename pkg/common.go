package common

import (
	//"embed"
	"encoding/json"
	"fmt"

	//"html/template"
	"io"
	//"io/fs"
	"net/http"
	"time"
)

func Http_req(url string) string {
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err.Error()
	} else {
		req.Header.Set("user-agent", "json-prettifier")
		resp, err := client.Do(req)
		if err != nil {
			return err.Error()
		} else {
			body, err := io.ReadAll(resp.Body)
			defer resp.Body.Close()
			if err != nil {
				return err.Error()
			} else {
				return string(body)
			}
		}
	}
}

func Parse_json(jsonString string) string {

	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(jsonString), &jsonMap)
	var response string

	full_json, err := json.MarshalIndent(jsonMap, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	response = string(full_json)
	return response
}
