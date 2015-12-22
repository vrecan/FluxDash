package merge

import (
	ST "github.com/fatih/structs"
)

//Merge src fields into destination
func Merge(src interface{}, dst interface{}, ignore ...string) interface{} {
	srcStruct := ST.New(src)
	dstStruct := ST.New(dst)
main:
	for _, field := range srcStruct.Fields() {
		_, ok := dstStruct.FieldOk(field.Name())
		if !ok {

			continue
		}
		for _, ign := range ignore {
			//skip field if it's on the ignore list
			if ign == field.Name() {
				continue main
			}
		}
		err := dstStruct.Field(field.Name()).Set(field.Value())
		if nil != err {
			panic(err)
		}
	}
	return dst
}
