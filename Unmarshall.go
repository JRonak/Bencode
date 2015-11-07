package Bencode

import (
	"errors"
	"fmt"
	"reflect"
)

func makeError(method, err string) error {
	return errors.New(method + "-> " + err)
}

func checkInt(r reflect.Kind) bool {
	switch r {
	case reflect.Int, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Int8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		return true
	default:
		return false
	}
}

func checkUint(r reflect.Kind) bool {
	switch r {
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		return true
	default:
		return false
	}
}

func checkStr(r reflect.Kind) bool {
	if r != reflect.String {
		return false
	}
	return true
}

func setInt(r reflect.Value, data interface{}) error {
	if !r.CanSet() {
		return makeError("setInt", "Value cannot be set")
	}
	rkind := r.Kind()
	dreflect := reflect.ValueOf(data)
	dkind := dreflect.Kind()
	if (!checkInt(rkind) && !checkUint(rkind)) || (!checkInt(dkind) && !checkUint(dkind)) {
		return makeError("setInt", "Value not Int type R:"+rkind.String()+" data:"+dkind.String())
	}

	if checkUint(r.Kind()) {
		if checkUint(dreflect.Kind()) {
			r.SetUint(dreflect.Uint())
		} else {
			r.SetUint(uint64(dreflect.Int()))
		}
	} else {
		if checkUint(dreflect.Kind()) {
			r.SetInt(int64(dreflect.Uint()))
		} else {
			r.SetInt(dreflect.Int())
		}
	}
	return nil
}

func setString(r reflect.Value, data interface{}) error {
	if !r.CanSet() {
		return makeError("setString", "Value cannot be set")
	}
	rkind := r.Kind()
	dreflect := reflect.ValueOf(data)
	dkind := dreflect.Kind()
	if !checkStr(rkind) || !checkStr(dkind) {
		return makeError("setString", "Value not string tpye R:"+rkind.String()+" data:"+dkind.String())
	}
	r.SetString(dreflect.String())
	return nil
}

func setArray(r reflect.Value, data interface{}) error {
	if !r.CanSet() {
		return makeError("setArray", "Value cannot be set")
	}
	rkind := r.Kind()
	dreflect := reflect.ValueOf(data)
	dkind := dreflect.Kind()
	if rkind != reflect.Array || dkind != reflect.Slice {
		return makeError("setArray", "Value not array tpye R:"+rkind.String()+" data:"+dkind.String())
	}
	rtype := r.Type()
	j := 0
	for i := 0; i < dreflect.Len() && j < r.Len(); i++ {
		ptr := reflect.Indirect(reflect.New(rtype.Elem()))
		e := detectType(ptr, dreflect.Index(i).Interface())
		if e != nil {
			return e
		}
		fmt.Println(ptr.Interface())
		r.Index(j).Set(ptr)
		j += 1
	}
	return nil
}

func setlist(r reflect.Value, data interface{}) error {
	if !r.CanSet() {
		return makeError("setList", "Value cannot be set")
	}
	rkind := r.Kind()
	dreflect := reflect.ValueOf(data)
	dkind := dreflect.Kind()
	if rkind != reflect.Slice || dkind != reflect.Slice {
		return makeError("setList", "Value not slice type R:"+rkind.String()+" data:"+dkind.String())
	}
	rtype := r.Type()
	for i := 0; i < dreflect.Len(); i++ {
		ptr := reflect.Indirect(reflect.New(rtype.Elem()))
		e := detectType(ptr, dreflect.Index(i).Interface())
		if e != nil {
			return e
		}
		r.Set(reflect.Append(r, ptr))
	}
	return nil
}

func setStruct(r reflect.Value, data interface{}) error {
	if !r.CanSet() {
		return makeError("setList", "Value cannot be set")
	}
	rkind := r.Kind()
	rtype := r.Type()
	dreflect := reflect.ValueOf(data)
	dkind := dreflect.Kind()
	if rkind != reflect.Struct || dkind != reflect.Map {
		return makeError("setStruct", "Value not Struct type R:"+rkind.String()+" data:"+dkind.String())
	}
	keys := dreflect.MapKeys()
	for i := 0; i < r.NumField(); i++ {
		name := getTag(rtype.Field(i))
		if isIgnore(rtype.Field(i)) {
			continue
		}
		index := searchKey(keys, name)
		if index == -1 {
			continue
		}
		e := detectType(r.Field(i), dreflect.MapIndex(keys[index]).Interface())
		if e != nil {
			return e
		}

	}
	return nil
}

func searchKey(rs []reflect.Value, key string) int {
	for i := 0; i < len(rs); i++ {
		if rs[i].Kind() != reflect.String {
			continue
		}
		if rs[i].String() == key {
			return i
		}
	}
	return -1
}

func detectType(r reflect.Value, data interface{}) error {
	switch r.Kind() {
	case reflect.Struct:
		return setStruct(r, data)
	case reflect.String:
		return setString(r, data)
	case reflect.Int, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Int8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		return setInt(r, data)
	case reflect.Array:
		return setArray(r, data)
	case reflect.Slice:
		return setlist(r, data)
	case reflect.Ptr:
		return makeError("detectType", "Pointer not supported")
	default:
		return makeError("detectType", "Unsupported type "+r.Kind().String())
	}
}

func Check(r reflect.Value, data interface{}) error {
	return detectType(r, data)
}

func Unmarshall(str string, into interface{}) error {
	if reflect.ValueOf(into).CanSet() {
		return makeError("Unmarshall", "Not a pointer tpye")
	}
	data, e := DecodeString(str)
	if e != nil {
		return e
	}
	r := reflect.ValueOf(into)
	if r.Kind() == reflect.Ptr && !r.IsNil() {
		r = r.Elem()
	}
	e = detectType(r, data)
	if e != nil {
		return e
	}
	return nil
}
