package m2s

import (
	"log"
	"reflect"
	"strconv"
	"strings"
)

func Map2Struct(m map[string]string, targ interface{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()
	v := reflect.ValueOf(targ)
	if v.Kind() != reflect.Ptr {
		log.Fatal("Please input a Point")
	}
	// if targ is ptr, change to elem
	v = getElem(v)
	for mk, mv := range m {
		item := v.FieldByName(strings.Title(mk))
		if !item.CanSet() {
			continue
		} else {
			// 将map中的内容set到对应的value中
			switch item.Type().Name() {
			case "int":
				fallthrough
			case "int64":
				data, err := strconv.ParseInt(mv, 10, 64)
				if err != nil {
					panic("Config wrong type, should be int")
				}
				item.SetInt(data)
			case "float":
				fallthrough
			case "float64":
				data, err := strconv.ParseFloat(mv, 64)
				if err != nil {
					panic("Config wrong type, should be float")
				}
				item.SetFloat(data)
			case "bool":
				data, err := strconv.ParseBool(mv)
				if err != nil {
					panic("Config wrong type, should be bool")
				}
				item.SetBool(data)
			case "string":
				item.SetString(mv)
			}
		}
	}
}

func getElem(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}
