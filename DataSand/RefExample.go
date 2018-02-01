package main

import ("reflect"
        "./net"
	"fmt"
	"strings"
)

func main2(){
	structInstance := &net.Packet{}
	sql := getCreateSqlStatementForStruct(structInstance,"dbname")
	fmt.Println(sql)
}

func getCreateSqlStatementForStruct(any interface{}, dbName string) string {
	v := reflect.ValueOf(any).Elem()

	typeName := v.Type().Name()
	sql := "CREATE TABLE IF NOT EXISTS "+dbName+"."+typeName +" (\n"
	sql += addFields(reflect.TypeOf(any), typeName)
	sql+= ");\n"

	return sql
}

func addFields(interfaceType reflect.Type, typeName string) string {

	sql := ""
	typeName = "  "+typeName
	if interfaceType.Kind()==reflect.Ptr {
		interfaceType = interfaceType.Elem()
	}

	for fieldIndex := 0; fieldIndex<interfaceType.NumField(); fieldIndex++ {

		field := interfaceType.Field(fieldIndex)

		fieldName := field.Name
		fieldKind := field.Type.Kind()

		if field.Type.Kind()==reflect.Ptr {
			fieldKind = field.Type.Elem().Kind()
		}

		switch fieldKind {
		case reflect.Bool:
			sql += typeName + fieldName + " bool ,\n"
		case reflect.Uint32:
			sql += typeName + fieldName + " integer ,\n"
		case reflect.Struct:
			sql += addFields(field.Type, fieldName+"_")
		case reflect.String:
			size := "128"
			if strings.Contains(strings.ToLower(fieldName), "uuid") {
				size = "64"
			}
			sql += typeName + fieldName + " VARCHAR[" + size + "] ,\n"
		}
	}
	return sql
}