package saban

import (
	"bytes"
	"github.com/tryfix/log"
	"html/template"
)

var headerTemplate *template.Template

const (
	prefixHeader = "saban.header"
)

func init() {
	var err error
	headerTemplate, err = template.ParseFiles("templates/header.txt")
	if err != nil {
		log.Fatal(log.WithPrefix(prefixHeader, "error parsing template"), err)
	}
}

type header struct {
	MustUnderstand int
	Actor          string
	NameSpace      string
	Type           string
}

func (h header) getContent() (out string, err error) {
	doc := &bytes.Buffer{}
	err = headerTemplate.Execute(doc, h)
	return doc.String(), nil
}
