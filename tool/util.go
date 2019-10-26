package main

import (
	"go/ast"
	"io/ioutil"
	"log"
	"os"
)

// 遍历文件处理;
func getServiceNames(cloudName string) []string {
	var serviceNames []string
	dir, err := ioutil.ReadDir(RootPath + ServiceRoot + cloudName)
	if err != nil {
		log.Println(err)
	}
	for _, file := range dir {
		if file.IsDir() {
			path := RootPath + ServiceRoot + cloudName + "\\" + file.Name() + "\\domain"
			// 通过判断是否存在domain文件夹,来判断是否为服务文件夹;
			if isDirExist(path) {
				serviceNames = append(serviceNames, file.Name())
			}
		}
	}

	return serviceNames
}

func isDirExist(path string) bool {
	// 通过判断是否存在domain文件夹,来判断是否为服务文件夹;
	if _, err := os.Stat(path); err != nil {
		if os.IsExist(err) {
			return true
		}
	} else {
		return true
	}
	return false
}

/*
go 类型转换成proto类型
*/
func getProtoType(fieldType string) string {
	switch fieldType {
	case "float64":
		return "double"
	case "float32":
		return "float"
	case "int32", "int":
		return "int32"
	case "int64":
		return "int64"
	case "uint32":
		return "uint32"
	case "uint64":
		return "uint64"
	case "bool":
		return "bool"
	case "string":
		return "string"
	case "int8", "int16":
		log.Fatal("int8, int16 is nonsupport type")
		return fieldType
	default:
		// 如果类型是引用其他结构体，则直接返回名称
		return fieldType
	}
}

// 获取指针的原始类型;
func getPointedType(expr ast.Expr) ast.Expr {
	// 表示指针,则取原始类型;
	if fieldType, ok := expr.(*ast.StarExpr); ok {
		return fieldType.X
	} else {
		return expr
	}
}
