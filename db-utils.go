package dbx

import (
	"reflect"
)

func isSlicePtr(res interface{}) (ok bool) {
	v := reflect.ValueOf(res)
	if v.Kind() != reflect.Ptr {
		return false
	}
	ev := v.Elem()
	ok = (ev.Kind() == reflect.Slice)
	return
}

func mk1ElemSlicePtr(res interface{}) interface{} {
	ev := reflect.ValueOf(res).Elem()
	et := ev.Type()
	st := reflect.SliceOf(et)
	sv := reflect.New(st)
	return sv.Interface()
}

func copySliceElem(ptrSl interface{}, res interface{}) {
	s0 := reflect.ValueOf(ptrSl).Elem().Index(0)
	rv := reflect.ValueOf(res).Elem()
	rv.Set(s0)
}

func sliceLen(ptrSl interface{}) int {
	sv := reflect.ValueOf(ptrSl).Elem()
	return sv.Len()
}
