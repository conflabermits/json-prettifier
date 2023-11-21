package common

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
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

type ResultDetails struct {
	Success  bool
	URL      string
	Response string
}

//go:embed static
var content embed.FS

func Web(port string) {
	filesys := fs.FS(content)
	tmpl := template.Must(template.ParseFS(filesys, "static/index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		reqURL := r.FormValue("url")
		fmt.Println("Request URL: " + reqURL)
		httpResponse := Http_req(reqURL)
		response := Parse_json(httpResponse)
		fmt.Println("Response: " + response)

		result := ResultDetails{
			Success:  true,
			URL:      reqURL,
			Response: response,
		}
		tmpl.Execute(w, result)
	})

	http.ListenAndServe(":"+port, nil)
}
