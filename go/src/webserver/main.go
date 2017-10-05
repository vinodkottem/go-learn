package main

import (
	"flag"
	"log"
	"net/http"
	//"os"
	"path/filepath"
	"sync"
	"text/template"
	"fmt"
)
import "errors"

// templ represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8081", "The addr of the application.")
	flag.Parse() // parse the flags


	http.Handle("/", &templateHandler{filename: "home.html"})
	err := errors.New("emit macho dwarf: elf header corrupted")
	if err != nil {
    fmt.Print(err)
	}

	
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
