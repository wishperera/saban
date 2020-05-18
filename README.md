# saban
#### Golang soap encoder
![include](./doc/img/soap.png)

Saban is a naive soap encoder using golang. Saban supports nested fields like structs,pointers,slices etc.

##### Define Payload
```golang
//define payload in golang
type SampleSoapPayload struct{
    RecordOne int `saban: "recordOne"`
    RecordTwo string `saban: "recordTwo"`
    RecordThree struct{
        SubRecordOne Float64 `saban:"subRecordOne"`
    } `saban: "recordThree"'`
    RecordFour *string `saban: "recordFour"`
    RecordFive  []int  `saban: "recordFive"`
}
```
`"saban"` tag can be used similar to json tags to instruct the encoder to use a given tag as key. if not provided the 
Name of the structField will be used as the key. 


##### Configuration

use the `MarshallerConfig` type to provide additional information like headers, XML Namespaces etc. The following fields
are supported at the moment.

| Option                    | Description                              | Default Value                          | Mandatory |
|---------------------------|------------------------------------------|----------------------------------------|-----------|
|conf.Envelope.EncodingStyle|envelope encoding style                   |http://www.w3.org/2003/05/soap-encoding | `yes`     | 
|conf.Header.Actor          |actor endpoint                            |""                                      | `no`      |
|conf.Header.MustUnderstand |instruct receiver to understand headers   |false                                   | `no`      |
|conf.Header.NameSpace      |XML Namespace for header                  |""                                      | `no`      |
|conf.Header.Type           |Header type (supported shown below)       |"Trans"                                 | `no`      |
|conf.Body.NameSpace        |XML Namespace for body                    | `none`                                 | `yes`     | 
|conf.Log                   |[tryfix](https://github.com/tryfix/log)log| a noop logger is implemented           | `no`      |

for detailed description of parameters please refer [here](https://www.w3schools.com/xml/xml_soap.asp)

```golang
//configuration parameters
conf := saban.MarshallerConfig{}
conf.Envelope.EncodingStyle = "http://www.w3.org/2003/05/soap-encoding"
conf.Header.Actor = "https://www.w3schools.com/code/"
conf.Header.MustUnderstand = true
conf.Header.NameSpace = "https://www.w3schools.com/transaction/"
conf.Header.Type = saban.HeaderTypeTransaction
conf.Body.NameSpace = "http://www.example.org/stock"

// default log level for the encoder is "DEBUG" 
// to change use following syntax
// tryfix log : "https://github.com/tryfix/log"

lg := log.NewLog(log.WithColors(true), log.WithFilePath(true), log.WithLevel("INFO"))
conf.Log = lg.Log(log.Prefixed("name-for-prefix"))
```

##### Marshaller

```go
marsh := saban.NewEncoder(conf)
var sb = "random_string"
data := SampleSoapPayload{}
data.RecordOne =  1000
data.RecordTwo = "string_two"
data.RecordThree.SubRecordOne = 10.04
data.RecordFour = &sb
data.RecordFive = []int{1,2,3,4,5}

byt,err := marsh.Encode(data)
if err != nil{
    log.Error("error marshalling:",err)
}

log.Info("success marshalling:",string(byt))
```

##### complete Example

```go
package main

import (
	"github.com/tryfix/log"
	"github.com/wishperera/saban"
)

func main() {

	type SampleSoapPayload struct{
    		RecordOne int `saban: "recordOne"`
    		RecordTwo string `saban: "recordTwo"`
    		RecordThree struct{
        		SubRecordOne Float64 `saban:"subRecordOne"`
    		} `saban: "recordThree"'`
    		RecordFour *string `saban: "recordFour"`
    		RecordFive  []int  `saban: "recordFive"`
	}

	conf := saban.MarshallerConfig{}
	conf.Envelope.EncodingStyle = "http://www.w3.org/2003/05/soap-encoding"
	conf.Header.Actor = "https://www.w3schools.com/code/"
	conf.Header.MustUnderstand = true
	conf.Header.NameSpace = "https://www.w3schools.com/transaction/"
	conf.Header.Type = saban.HeaderTypeTransaction
	conf.Body.NameSpace = "http://www.example.org/stock

	marsh := saban.NewEncoder(conf)
	var sb = "random_string"
	data := SampleSoapPayload{}
	data.RecordOne =  1000
	data.RecordTwo = "string_two"
	data.RecordThree.SubRecordOne = 10.04
	data.RecordFour = &sb
	data.RecordFive = []int{1,2,3,4,5}

	enc := saban.NewEncoder(conf)
	byt, err := enc.Encode(data)
	if err != nil {
		log.Error("error marshalling:",err)

	}

	log.Info("success marshalling:",string(byt))
}
```
try on playground [click-here](https://play.golang.org/p/lq8fAVd1hG5)