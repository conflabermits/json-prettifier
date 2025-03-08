package common

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"reflect"
	"strings"
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
		log.Println(err)
	}
	response = string(full_json)
	return response
}

func get_wheel_count(jsonString string) string {

	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(jsonString), &jsonMap)

	childarraylength := len(jsonMap["result"].(map[string]interface{})["data"].(map[string]interface{})["json"].([]interface{}))
	response := fmt.Sprintf("%d", childarraylength)
	//fmt.Println(response)
	return response
}

type ResultDetails struct {
	Type     string
	Input    string
	Schema   string
	Response string
}

//go:embed static
var content embed.FS

func Web(port string) {
	filesys := fs.FS(content)
	tmpl := template.Must(template.ParseFS(filesys, "static/index.html"))

	http.HandleFunc("/json-prettifier", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		reqURL := r.FormValue("url")
		jsonInput := r.FormValue("json")
		uglyJson := ""
		inputType := ""
		inputValue := ""
		if reqURL == "" && jsonInput == "" {
			http.Error(w, "No input provided", http.StatusBadRequest)
			return
		}
		if reqURL != "" && jsonInput != "" {
			http.Error(w, "Provide only one input", http.StatusBadRequest)
			return
		}
		if reqURL != "" && jsonInput == "" {
			inputType = "URL"
			inputValue = reqURL
			uglyJson = Http_req(reqURL)
		}
		if reqURL == "" && jsonInput != "" {
			inputType = "JSON"
			inputValue = jsonInput
			uglyJson = jsonInput
		}
		prettyJson := Parse_json(uglyJson)
		schema, err := analyzeJSON([]byte(prettyJson))
		if err != nil {
			log.Println("Error:", err)
			return
		}

		result := ResultDetails{
			Type:     inputType,
			Input:    inputValue,
			Schema:   schema,
			Response: prettyJson,
		}
		tmpl.Execute(w, result)
	})

	http.HandleFunc("/json-prettifier/wheelcount", func(w http.ResponseWriter, r *http.Request) {
		reqURL := "https://officedrummerwearswigs.com/api/trpc/songRequest.getLatest"
		httpResponse := Http_req(reqURL)
		response := get_wheel_count(httpResponse)
		fmt.Fprint(w, response)
	})

	http.ListenAndServe(":"+port, nil)
}

func analyzeJSON(jsonData []byte) (string, error) {
	var data interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	schema := analyzeValue(data, "")
	return schema, nil
}

func analyzeValue(value interface{}, path string) string {
	switch v := value.(type) {
	case map[string]interface{}:
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("(dict, length: %d)\n", len(v)))
		for key, val := range v {
			newPath1 := fmt.Sprintf("    %s", path)
			keyString1 := fmt.Sprintf("\"%s\": ", key)
			sb.WriteString(newPath1 + keyString1)
			sb.WriteString(analyzeValue(val, newPath1))
		}
		return sb.String()
	case []interface{}:
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("(array, length: %d)\n", len(v)))

		if len(v) > 0 {
			// Analyze the first element to determine the structure of the array items
			firstElement := v[0]
			switch element := firstElement.(type) {
			case map[string]interface{}:
				for key, val := range element {
					newPath2 := fmt.Sprintf("    %s", path)
					keyString2 := fmt.Sprintf("\"%s\": ", key)
					sb.WriteString(newPath2 + keyString2)
					sb.WriteString(analyzeValue(val, path))
				}
			default: //If the first element is not a map, just print the type
				sb.WriteString(fmt.Sprintf("    %s\n", getType(element)))
			}
		}
		return sb.String()

	default:
		return fmt.Sprintf("%s\n", getType(v))
	}
}

func getType(value interface{}) string {
	if value == nil {
		return "null"
	}
	t := reflect.TypeOf(value)
	if t == nil {
		return "unknown"
	}

	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "int"
	case reflect.Float32, reflect.Float64:
		return "float"
	case reflect.Bool:
		return "bool"
	default:
		// Attempt type assertion for time.Time
		if _, ok := value.(string); ok {
			return "date" // Or a more specific date/time check if needed
		}
		return t.String()
	}
}
