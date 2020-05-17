package saban

import (
	"bytes"
	"github.com/tryfix/log"
	"html/template"
)

var envelopeTemplate *template.Template

const (
	prefixEnvelope       = "saban.envelope"
	defaultEncodingStyle = "http://www.w3.org/2003/05/soap-encoding"
)

func init() {
	var err error
	envelopeTemplate, err = template.ParseFiles("templates/envelope.txt")
	if err != nil {
		log.Fatal(log.WithPrefix(prefixEnvelope, "error parsing template"), err)
	}
}

type envelope struct {
	EncodingStyle string `validate:"url"`
	Header        string
	Body          string
}

func (e envelope) validate() error {
	return validator.Struct(e)
}

func (e envelope) getContent() (out string, err error) {
	doc := &bytes.Buffer{}
	err = envelopeTemplate.Execute(doc, e)
	return doc.String(), nil
}
