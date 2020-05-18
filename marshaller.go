package saban

import (
	"encoding/xml"
	"fmt"
	"github.com/tryfix/log"
	"github.com/wishperera/saban/logger"
	"reflect"
	"strings"
)

type headerType string

const (
	HeaderTypeTransaction headerType = "Trans"
)

func (h headerType) String() string {
	return string(h)
}

type marshaller struct {
	config *MarshallerConfig
	output string
}

type MarshallerConfig struct {
	Envelope configEnvelope
	Header   configHeader
	Body     configBody
	Log      log.Logger
}

type configBody struct {
	NameSpace string `validate:"url,required"`
}

type configEnvelope struct {
	EncodingStyle string `validate:"url,required"`
}

type configHeader struct {
	MustUnderstand bool
	Type           headerType
	NameSpace      string `validate:"url,required"`
	Actor          string `validate:"url,omitempty"`
}

// override empty configurations with defaults
func (m marshaller) validate() {
	if m.config == nil || reflect.DeepEqual(m.config, MarshallerConfig{}) {
		cf := MarshallerConfig{}
		cf.Envelope.EncodingStyle = defaultEncodingStyle
		cf.Header.Actor = ""
		cf.Header.MustUnderstand = false
		cf.Header.Type = HeaderTypeTransaction
		cf.Log = logger.NoopLogger
		m.config = &cf
	}
	if m.config.Envelope.EncodingStyle == "" {
		m.config.Envelope.EncodingStyle = defaultEncodingStyle
	}
	if m.config.Log == nil {
		m.config.Log = logger.NoopLogger
	}

	if err := validator.Struct(m.config); err != nil {
		m.config.Log.Fatal("error validating config", err)
	}

}

func NewEncoder(config MarshallerConfig) *marshaller {
	m := new(marshaller)
	m.config = &config
	m.validate()
	return m
}

func (m marshaller) Encode(in interface{}) (out []byte, err error) {

	defer func() {
		if err != nil {
			m.config.Log.Error("error marshalling input", "input:", in, "error:", err)
		} else {
			m.config.Log.Debug("success marshalling input", "input:", in, "out:", string(out))
		}

	}()

	m.parse(in, true)
	body := fmt.Sprintf(`<soap:Body xmlns:m="%v">%v</soap:Body>`, m.config.Body.NameSpace, m.output)

	header := header{}
	if m.config.Header.MustUnderstand {
		header.MustUnderstand = 1
	} else {
		header.MustUnderstand = 0
	}
	header.Actor = m.config.Header.Actor
	header.NameSpace = m.config.Header.NameSpace
	header.Type = m.config.Header.Type.String()

	h, err := header.getContent()
	if err != nil {
		return nil, err
	}

	envelope := envelope{}
	envelope.EncodingStyle = m.config.Envelope.EncodingStyle
	envelope.Body = body
	envelope.Header = h

	env, err := envelope.getContent()
	if err != nil {
		return nil, err
	}

	out, err = xml.Marshal(env)
	return out, err

}

func (m *marshaller) remove(s string) {
	m.output = strings.TrimSuffix(m.output, s)
}

func (m *marshaller) append(s string) {
	m.output = m.output + s
}

func (m *marshaller) parse(in interface{}, base bool) {
	var iter int
	t := reflect.TypeOf(in)
	k := t.Kind()
	switch k {
	case reflect.Struct:
		value := reflect.ValueOf(in)
		for iter = 0; iter < value.NumField(); iter++ {

			var n string
			f := t.Field(iter)
			tag, ok := f.Tag.Lookup("saban")
			if !ok {
				n = f.Name
			} else {
				n = tag
			}
			//fmt.Printf("field: %v \n",value.Field(iter))
			m.append(fmt.Sprintf("<m:%v>", n))
			m.parse(value.Field(iter).Interface(), false)
			m.append(fmt.Sprintf("</m:%v>", n))

		}

	case reflect.Ptr:
		m.parse(reflect.ValueOf(in).Elem().Interface(), false)

	case reflect.Slice:
		m.append("[")
		value := reflect.ValueOf(in)
		for mapIter := 0; mapIter < value.Len(); mapIter++ {
			m.parse(value.Index(mapIter).Interface(), false)
			m.append(",")
		}
		m.remove(",")
		m.append("]")

	default:
		m.append(fmt.Sprintf("%v", reflect.ValueOf(in)))
		//fmt.Printf("value: %v \n",reflect.ValueOf(in))
	}
}
