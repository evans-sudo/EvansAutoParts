package main

import (
	"EvansAutoParts/src/bw/ctrl"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)




func main() {
	templateCache, _ := buildTemplateCache()
	ctrl.Setup(templateCache)

	go http.ListenAndServe(":3000", nil)


	go func() {
		for range time.Tick(300 * time.Millisecond) {
			tc, isupdated := buildTemplateCache()
			if isupdated {
				ctrl.SetTemplateCache(tc)
			}
		}
	}()

	log.Println("Server Started Press Enter to Exit")
	fmt.Scanln()
}

var lastModTime time.Time = time.Unix(0, 0)

func buildTemplateCache() (*template.Template, bool) {
	needUpdate := false


	f, err := os.Open("template")
	if err != nil {
		log.Fatalf("Failed to open templates directory: %v", err)
	}
	defer f.Close()

	fileInfos, _ := f.Readdir(-1)
	if err != nil {
		log.Fatalf("Failed to read templates directory: %v", err)
	}

	fileNames := []string{}
	for _, fi := range fileInfos {
		if fi.Mode().IsRegular() && strings.HasSuffix(fi.Name(), ".html") { // Ensure only HTML files
			filePath := "template/" + fi.Name()
			fileNames = append(fileNames, filePath)

			if fi.ModTime().After(lastModTime) {
				lastModTime = fi.ModTime()
				needUpdate = true
			}
		}
	}

	if len(fileNames) == 0 {
		log.Fatal("No templates found in the templates directory")
	}

	// Always parse templates initially
	tc, err := template.ParseFiles(fileNames...)
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	if needUpdate {
		log.Print("Template change detected, updating...")
		log.Println("Template update complete")
	}

	return tc, needUpdate
}