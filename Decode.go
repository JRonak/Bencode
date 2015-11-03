package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

var (
	invalidlEncoded error
	invalidInt      error
)

const (
	success = iota
	end     = iota
	str     = iota
	integer = iota
	fail    = iota
)

type Reader struct {
	*bufio.Reader
}

func init() {
	invalidlEncoded = errors.New("Invalid Bencode")
	invalidInt = errors.New("InvalidInt")
}

func checkErr(e error) {
	if e != nil {
		log.Println("Decode:" + e.Error())
		panic(e)
	}
}

func (reader *Reader) getByte() byte {
	b, e := reader.ReadByte()
	checkErr(e)
	return b
}

func (reader *Reader) getBytes(num int) []byte {
	b := make([]byte, num)
	for i := 0; i < num; i++ {
		temp := reader.getByte()
		b[i] = temp
	}
	return b
}

//func (reader *Reader)

func (reader *Reader) getInt() int {
	b, e := reader.ReadString('e')
	checkErr(e)
	num, e := strconv.Atoi(b[0 : len(b)-1])
	checkErr(e)
	return num
}

func (reader *Reader) getString() string {
	b, e := reader.ReadString(':')
	checkErr(e)
	num, e := strconv.Atoi(b[0 : len(b)-1])
	checkErr(e)
	data := reader.getBytes(num)
	return string(data)
}

func (reader *Reader) dict() map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	for {
		in, status := reader.decode()
		if status == end {
			break
		} else if status != str {
			checkErr(invalidlEncoded)
		}
		val, status := reader.decode()
		if status == fail {
			checkErr(invalidlEncoded)
		}
		m[in] = val
	}
	return m
}

func (reader *Reader) list() []interface{} {
	var list []interface{}
	for {
		in, status := reader.decode()
		if status == end {
			break
		} else if status == fail {
			return nil
		}
		list = append(list, in)
	}
	return list
}

func (reader *Reader) decode() (interface{}, int) {
	temp, _ := reader.ReadByte()
	if temp == 'd' {
		return reader.dict(), success
	} else if temp == 'l' {
		return reader.list(), success
	} else if temp <= '9' && temp >= '0' {
		reader.UnreadByte()
		return reader.getString(), str
	} else if temp == 'i' {
		return reader.getInt(), integer
	} else if temp == 'e' {
		return nil, end
	} else {
		return nil, fail
	}
	return nil, success
}

func (reader *Reader) BeginDecode() (map[interface{}]interface{}, error) {
	firstByte, err := reader.ReadByte()
	checkErr(err)
	if firstByte != 'd' {
		return nil, invalidlEncoded
	}
	//checkErr(reader.UnreadByte())
	x := reader.dict()
	return x, nil
}

func main() {
	defer func() {
		if i := recover(); i != nil {
			fmt.Println("Recovered")
		}
	}()
	handle, e := os.Open("a")
	checkErr(e)
	reader := Reader{bufio.NewReader(handle)}
	x, _ := reader.BeginDecode()
	fmt.Println(x)
}
