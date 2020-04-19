package common

import (
	"reflect"
	"unsafe"
)

type Result struct {
	Code   int       `json:"code"`
	Param  string    `json:"param"`
	Msg    string    `json:"msg"`
	Data   interface{}    `json:"data"`
}
var sizeOfMyStruct = int(unsafe.Sizeof(Result{}))
func ToBytes(re *Result) []byte{
	var x reflect.SliceHeader
	x.Len = sizeOfMyStruct
	x.Cap = sizeOfMyStruct
	x.Data = uintptr(unsafe.Pointer(re))
	return *(*[]byte)(unsafe.Pointer(&x))
}

func (re *Result) SuccessWithData(params string, msg string, data interface{}) {
	re.Code = 200
	re.Param = params
	re.Msg = msg
	re.Data = data
}

func (re *Result) SuccessWithoutData(msg string) {
	re.Code = 200
	re.Msg = msg
}

func (re *Result) ErrorWithMsg(code int, msg string) {
	re.Code = code
	re.Msg = msg
}