package Bencode

import (
	"bufio"
	"errors"
	"io"
	"log"
	"strconv"
	"strings"
)

var (
	invalidEncodedErr error
	invalidInt        error
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
	invalidEncodedErr = errors.New("Invalid Bencode")
	invalidInt = errors.New("InvalidInt")
}

func logErr(e error) {
	if e != nil {
		log.Println("Decode:" + e.Error())
	}
}

func (reader *Reader) getByte() (byte, error) {
	b, e := reader.ReadByte()
	if e != nil {
		logErr(e)
		return 0, invalidEncodedErr
	}
	return b, nil
}

func (reader *Reader) getBytes(num int) ([]byte, error) {
	b := make([]byte, num)
	for i := 0; i < num; i++ {
		temp, e := reader.getByte()
		if e != nil {
			return nil, e
		}
		b[i] = temp
	}
	return b, nil
}

//func (reader *Reader)

func (reader *Reader) getInt() (int, error) {
	b, e := reader.ReadString('e')
	if e != nil {
		logErr(e)
		return 0, invalidEncodedErr
	}
	num, e := strconv.Atoi(b[0 : len(b)-1])
	if e != nil {
		logErr(e)
		return 0, invalidInt
	}
	return num, nil
}

func (reader *Reader) getString() (string, error) {
	b, e := reader.ReadString(':')
	if e != nil {
		logErr(e)
		return "", invalidEncodedErr
	}
	num, e := strconv.Atoi(b[0 : len(b)-1])
	if e != nil {
		logErr(e)
		return "", invalidInt
	}
	data, e := reader.getBytes(num)
	if e != nil {
		return "", invalidEncodedErr
	}
	return string(data), nil
}

func (reader *Reader) dict() (map[string]interface{}, error) {
	m := make(map[string]interface{})
	for {
		key, status := reader.decode()
		if status == end {
			break
		} else if status == fail {
			return nil, invalidEncodedErr
		} else if status != str {
			logErr(invalidEncodedErr)
			return nil, invalidEncodedErr
		}
		value, status := reader.decode()
		if status == fail {
			return nil, invalidEncodedErr
		}
		m[key.(string)] = value
	}
	return m, nil
}

func (reader *Reader) list() ([]interface{}, error) {
	var list []interface{}
	for {
		in, status := reader.decode()
		if status == end {
			break
		} else if status == fail {
			return nil, invalidEncodedErr
		}
		list = append(list, in)
	}
	return list, nil
}

func (reader *Reader) decode() (interface{}, int) {
	temp, e := reader.ReadByte()
	if e != nil {
		logErr(e)
		return nil, fail
	}
	switch temp {
	case 'd':
		dict, e := reader.dict()
		if e != nil {
			return nil, fail
		}
		return dict, success
	case 'l':
		list, e := reader.list()
		if e != nil {
			return nil, fail
		}
		return list, success
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		reader.UnreadByte()
		data, e := reader.getString()
		if e != nil {
			return nil, fail
		}
		return data, str
	case 'i':
		num, e := reader.getInt()
		if e != nil {
			return nil, fail
		}
		return num, integer
	case 'e':
		return nil, end
	default:
		return nil, fail
	}
}

func Decode(handle io.Reader) (map[string]interface{}, error) {
	reader := Reader{bufio.NewReader(handle)}
	firstByte, err := reader.ReadByte()
	if err != nil {
		logErr(err)
		return nil, err
	}
	if firstByte != 'd' {
		return nil, invalidEncodedErr
	}
	x, e := reader.dict()
	return x, e
}

func DecodeString(str string) (map[string]interface{}, error) {
	buf := strings.NewReader(str)
	reader := Reader{bufio.NewReader(buf)}
	firstByte, err := reader.ReadByte()
	if err != nil {
		logErr(err)
		return nil, err
	}
	if firstByte != 'd' {
		return nil, invalidEncodedErr
	}
	x, e := reader.dict()
	return x, e
}

/*
func main() {
	handle, e := os.Open("a")
	logErr(e)
	reader := Reader{bufio.NewReader(handle)}
	x, e := reader.BeginDecode()
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Println(x)
	}
}*/
