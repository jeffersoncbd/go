package main

import (
  "fmt"
  "io/ioutil"
)

type Page struct {
  Title string
  Body []byte
}

func main() {
  p1 := &Page{Title: "TestPage", Body: []byte("This is a simple page.")}
  p1.save()

  p2, _ := loadPage("TestPage")
  fmt.Println(string(p2.Body))
}

func (p *Page) save() error {
  fileName := p.Title + ".txt"
  return ioutil.WriteFile(fileName, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
  fileName := title + ".txt"
  body, err := ioutil.ReadFile(fileName)
  if err != nil {
    return nil, err
  }
  return &Page{Title: title, Body: body}, nil
}
