package greetings

import (
  "testing"
  "regexp"
)

func TestHelloName(t *testing.T) {
  name := "Jefferson"
  want := regexp.MustCompile(`\b`+name+`\b`)

  msg, err := Hello("Jefferson")
  if !want.MatchString(msg) || err != nil {
    t.Fatalf(`Hello("Jefferson") = %q, %v, want match for %#q nil`, msg, err, want)
  }
}

func TestHelloEmpty(t *testing.T) {
  msg, err := Hello("")
  if msg != "" || err == nil {
    t.Fatalf(`Hello("") = %q, %v, want "", error`, msg, err)
  }
}

