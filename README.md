# Bencode
Bencode library supports marshal and unmarshal bencoded data in go. Supports struct tags, similar to the JSON library.

## Installation
  ```
$ go get github.com/JRonak/Bencode
  ```
## Usage
```
package main

import (
	"fmt"
	"github.com/JRonak/Bencode"
)

type sample struct {
	SampleInt    int    `Bencode:"integer"`
	SampleString string `Bencode:"string"`
	SampleList   []int  `Bencode:"list"`
}

func sampleEncode() {
	s := sample{SampleInt: 5, SampleString: "Hello World!",
		SampleList: []int{1, 2, 3, 4}}
	becondeStr, err := Bencode.Encode(s)
	if err != nil {
		panic(err)
	}
	//	d7:integeri5e6:string12:Hello World!4:listli1ei2ei3ei4eee	fmt.Println(beconde)
	fmt.Println(becondeStr)
}

func sampleDecode() {
	decodeStr := "d7:integeri5e6:string12:Hello World!4:listli1ei2ei3ei4eee"
	var s sample
	err := Bencode.Unmarshall(decodeStr, &s)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", s)
}

func main() {
	sampleEncode()
	sampleDecode()
}

```
