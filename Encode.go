package Bencode

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

type sex struct {
	A int
	B int
}

type Writer struct {
	io.Writer
}

var (
	dictError        error
	unknownTypeError error
)

func init() {
	dictError = errors.New("Dictionary key is not a string type")
	unknownTypeError = errors.New("Unknown type error")
}

func (writer *Writer) writeString(data interface{}) error {
	str := data.(string)
	length := len(str)
	_, e := writer.Write([]byte(strconv.Itoa(length) + ":"))
	if e != nil {
		return e
	}
	_, e = writer.Write([]byte(str))
	if e != nil {
		return e
	}
	return nil
}

func (writer *Writer) writeInt(data interface{}) error {
	var e error
	r := reflect.ValueOf(data)
	if r.Kind() == reflect.Int {

		_, e = fmt.Fprintf(writer, "i%de", r.Int())
	} else {

		_, e = fmt.Fprintf(writer, "i%de", r.Uint())
	}
	if e != nil {
		return e
	}
	return nil
}

func (writer *Writer) writeList(data interface{}) error {
	_, e := writer.Write([]byte("l"))
	if e != nil {
		return e
	}
	r := reflect.ValueOf(data)
	for i := 0; i < r.Len(); i++ {
		e = writer.detectType(r.Index(i).Interface())
		if e != nil {
			return e
		}
	}
	_, e = writer.Write([]byte("e"))
	if e != nil {
		return e
	}
	return nil
}

func (writer *Writer) writeDict(data interface{}) error {
	_, e := writer.Write([]byte("d"))
	if e != nil {
		return e
	}
	r := reflect.ValueOf(data)
	keys := r.MapKeys()
	for _, key := range keys {
		if key.Kind() != reflect.String {
			return dictError
		}
		if e = writer.detectType(key.Interface()); e != nil {
			return e
		}
		if e = writer.detectType(r.MapIndex(key).Interface()); e != nil {
			return e
		}
	}
	_, e = writer.Write([]byte("e"))
	if e != nil {
		return e
	}
	return nil
}

func (writer *Writer) writeStruct(data interface{}) error {
	_, e := writer.Write([]byte("d"))
	if e != nil {
		return e
	}
	r := reflect.ValueOf(data)
	t := r.Type()
	for i := 0; i < t.NumField(); i++ {
		if (r.FieldByIndex(t.Field(i).Index).CanInterface()) == false {
			continue
		}
		if e = writer.detectType(t.Field(i).Name); e != nil {
			return e
		}
		if e = writer.detectType(r.FieldByIndex(t.Field(i).Index).Interface()); e != nil {
			return e
		}

	}
	_, e = writer.Write([]byte("e"))
	if e != nil {
		return e
	}
	return nil
}

func (writer *Writer) detectType(data interface{}) error {
	r := reflect.ValueOf(data)

	switch r.Kind() {
	case reflect.Ptr:
		writer.detectType(r.Elem().Interface())
	case reflect.Array, reflect.Slice:
		writer.writeList(data)
	case reflect.Map:
		if e := writer.writeDict(data); e != nil {
			return e
		}
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Int8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		if e := writer.writeInt(data); e != nil {
			return e
		}
	case reflect.String:
		if e := writer.writeString(data); e != nil {
			return e
		}
	case reflect.Struct:
		if e := writer.writeStruct(data); e != nil {
			return e
		}
	default:
		fmt.Println(r.Kind())

	}
	return nil
}

func Encode(data interface{}) string {
	buf := new(bytes.Buffer)
	writer := Writer{buf}
	writer.detectType(data)
	return buf.String()
}

/*
func main() {

	y := sex{}
	y.A = 1
	y.B = 2
	fmt.Println(Encode(y))
}
*/
