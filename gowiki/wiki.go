package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
  "os"
)

type Page struct {
  Title string
  Body []byte
}

var templates = template.Must(template.ParseFiles("templates/list.html", "templates/create.html", "templates/view.html", "templates/edit.html"))
var validPath = regexp.MustCompile("^/(edit|update|view|delete)/([a-zA-Z0-9]+)$")

func main() {
  http.HandleFunc("/", listHandler)
  http.HandleFunc("/create", createHandler)
  http.HandleFunc("/save", saveHandler)
  http.HandleFunc("/view/", makeHandler(viewHandler))
  http.HandleFunc("/edit/", makeHandler(editHandler))
  http.HandleFunc("/update/", makeHandler(updateHandler))
  http.HandleFunc("/delete/", makeHandler(deleteHandler))

  log.Fatal(http.ListenAndServe(":8080", nil))
}

func (p *Page) save() error {
  fileName := "pages/"+p.Title
  return ioutil.WriteFile(fileName, p.Body, 0600)
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
  m := validPath.FindStringSubmatch(r.URL.Path)
  if m == nil {
    http.NotFound(w, r)
    return "", errors.New("Invalid Page Title")
  }
  return m[2], nil
}

func loadPage(title string) (*Page, error) {
  fileName := "pages/"+title
  body, err := ioutil.ReadFile(fileName)
  if err != nil {
    return nil, err
  }
  return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
  err := templates.ExecuteTemplate(w, tmpl+".html", p)
  if err != nil {
    http.Error(w, "renderTemplate: "+err.Error(), http.StatusInternalServerError)
  }
}

func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
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
  files, err := ioutil.ReadDir("pages/")
  if err != nil {
    fmt.Fprintf(w, err.Error())
    return
  }

  err = templates.ExecuteTemplate(w, "list.html", files)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
  p, err := loadPage(title)
  if err != nil {
    http.Redirect(w, r, "/edit/"+title, http.StatusFound)
    return
  }

  renderTemplate(w, "view", p)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
  renderTemplate(w, "create", nil)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
  title := r.FormValue("title")
  body := r.FormValue("body")

  validTitle := regexp.MustCompile("^([a-zA-Z0-9]+)$")
  validationErr := validTitle.FindStringSubmatch(title)
  if validationErr == nil {
    fmt.Fprintf(w, "Caracteres inválidos no título: "+title)
    return
  }

  p := &Page{Title: title, Body: []byte(body)}
  err := p.save()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func updateHandler(w http.ResponseWriter, r *http.Request, title string) {
  body := r.FormValue("body")
  p := &Page{Title: title, Body: []byte(body)}
  err := p.save()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
  p, err := loadPage(title)
  if err != nil {
    p = &Page{Title: title}
  }

  renderTemplate(w, "edit", p)
}

func deleteHandler(w http.ResponseWriter, r *http.Request, title string) {
  _, err := loadPage(title)
  if err != nil {
    http.NotFound(w, r)
  }

  os.Remove("pages/"+title)

  http.Redirect(w, r, "/", http.StatusFound)
}

