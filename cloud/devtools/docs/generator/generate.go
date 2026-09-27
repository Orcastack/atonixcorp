package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Metadata struct {
	Services map[string]struct {
		Methods []string `json:"methods"`
		Path    string   `json:"path"`
	} `json:"services"`
}

func main() {
	data, err := os.ReadFile("generator/metadata.json")
	if err != nil {
		panic(err)
	}

	var meta Metadata
	json.Unmarshal(data, &meta)

	generateSDKDocs("go", meta)
	generateSDKDocs("python", meta)
	generateSDKDocs("java", meta)
	generateSDKDocs("javascript", meta)
	generateAPIDocs(meta)

	fmt.Println("Documentation generated.")
}

func generateSDKDocs(lang string, meta Metadata) {
	out := fmt.Sprintf("# %s SDK Documentation\n\n", capitalize(lang))

	for name, svc := range meta.Services {
		out += fmt.Sprintf("## %s\n", capitalize(name))
		out += fmt.Sprintf("Path: `%s`\n\n", svc.Path)

		out += "### Methods\n"
		for _, m := range svc.Methods {
			out += fmt.Sprintf("- `%s`\n", m)
		}
		out += "\n"
	}

	os.WriteFile(fmt.Sprintf("sdk/%s.md", lang), []byte(out), 0644)
}

func generateAPIDocs(meta Metadata) {
	out := "# API Endpoints\n\n"

	for name, svc := range meta.Services {
		out += fmt.Sprintf("## %s\n", capitalize(name))
		out += fmt.Sprintf("Base Path: `%s`\n\n", svc.Path)

		out += "### Endpoints\n"
		for _, m := range svc.Methods {
			out += fmt.Sprintf("- `%s %s/%s`\n", httpVerb(m), svc.Path, m)
		}
		out += "\n"
	}

	os.WriteFile("api/endpoints.md", []byte(out), 0644)
}

func httpVerb(method string) string {
	switch method {
	case "create", "uploadObject":
		return "POST"
	case "delete", "deleteContainer", "deleteObject":
		return "DELETE"
	default:
		return "GET"
	}
}

func capitalize(s string) string {
	return string(s[0]-32) + s[1:]
}
