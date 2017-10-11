package main

import (
  "github.com/cloudfoundry-community/go-cfenv"
  "html/template"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "regexp"
  "fmt"
  "strings"
  "path/filepath"
  "strconv"
)

const (
  DEFAULT_PORT = "8081"
)

type Page struct {
  Title string
  Body  []byte
}

func (p *Page) save() error {
  filename := p.Title + ".txt"
  return ioutil.WriteFile("documents/"+filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
  appEnv, _ := cfenv.Current()
  if appEnv != nil {
    log.Println("ENV VARIABLES ARE HERE")
  }
  filename := title + ".txt"
  body, err := ioutil.ReadFile("documents/"+filename)
  if err != nil {
    return nil, err
  }
  return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
  p, err := loadPage(title)
  if err != nil {
    http.Redirect(w, r, "/edit/"+title, http.StatusFound)
    return
  }
  renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
  p, err := loadPage(title)
  if err != nil {
    p = &Page{Title: title}
  }
  renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
  body := r.FormValue("body")
  p := &Page{Title: title, Body: []byte(body)}
  err := p.save()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
  err := templates.ExecuteTemplate(w, tmpl+".html", p)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil {
      http.NotFound(w, r)
      return
    }
    fn(w, r, m[2])
  }
}

func listHandler(w http.ResponseWriter, r *http.Request) {
    files, err := ioutil.ReadDir("./documents/")
    if err != nil {
        log.Fatal(err)
    }
    var dd = ""
    for _, file := range files {
        fname := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
        dd = dd + "<div>"+fname+" &emsp;[<a href=/view/"+fname+">edit</a>]</div>"
     }
     str1 := "window.location = '/view/'+document.getElementById('test').value"
      quote1 := strconv.Quote(str1)
      str := "<input type='text' id='test'  value='' /><button onclick="+quote1+" id='testButton'>[new]</button>"
     dd = dd+str
     fmt.Fprintf(w, dd)    

}

func main() {
  http.HandleFunc("/view/", makeHandler(viewHandler))
  http.HandleFunc("/edit/", makeHandler(editHandler))
  http.HandleFunc("/save/", makeHandler(saveHandler))
  http.HandleFunc("/", listHandler)


  var port string
  if port = os.Getenv("PORT"); len(port) == 0 {
    port = DEFAULT_PORT
  }
  log.Println("PORT "+port)
  log.Fatal(http.ListenAndServe(":"+port, nil))
}