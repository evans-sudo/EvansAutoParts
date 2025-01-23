package main

import (
	"EvansAutoParts/src/bw/ctrl"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
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


	f, _ := os.Open("templates")

	fileInfos, _ := f.Readdir(-1)
	fileNames := make([]string, len(fileInfos))
	for idx, fi := range fileInfos {
		if fi.ModTime().After(lastModTime) {
			lastModTime = fi.ModTime()
			needUpdate = true
		}
		fileNames[idx] = "templates/" + fi.Name()
	}

	var tc *template.Template
	if needUpdate {
		log.Print("Template change detected, updating...")
		tc = template.Must(template.ParseFiles(fileNames...))
		log.Println("template update complete")
	}

	return tc, needUpdate
}