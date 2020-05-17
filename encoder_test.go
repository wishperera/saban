package saban

import "testing"

func TestEncoder(t *testing.T) {
	type ABC struct {
		FieldOne   int    `saban:"one"`
		FieldTwo   string `saban:"two"`
		FieldThree struct {
			SubFieldOne int `saban:"threeSubOne"`
		} `saban: "three"`
		FieldFour *string `saban:"four"`
		FieldFive []int   `saban:"five"`
	}

	conf := MarshallerConfig{}
	conf.Envelope.EncodingStyle = "http://www.w3.org/2003/05/soap-encoding"
	conf.Header.Actor = "https://www.w3schools.com/code/"
	conf.Header.MustUnderstand = true
	conf.Header.NameSpace = "https://www.w3schools.com/transaction/"
	conf.Header.Type = HeaderTypeTransaction
	conf.Body.NameSpace = "http://www.example.org/stock"

	var rs = "random string"
	object := ABC{
		FieldOne: 1000,
		FieldTwo: "value one",
		FieldThree: struct {
			SubFieldOne int `saban:"threeSubOne"`
		}{SubFieldOne: 300},
		FieldFour: &rs,
		FieldFive: []int{1, 2, 3, 4, 5},
	}

	enc := NewEncoder(conf)
	out, err := enc.Marshal(object)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	t.Log(string(out))

}
