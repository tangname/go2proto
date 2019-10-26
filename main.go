package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
)

var RootPath = "E:\\gm-workspace\\go2proto-master\\"
var ServiceRoot = "t1\\"
var CloudName = "stocking"
var ServiceName = "test1"

var cloud Cloud

/*func init() {
	flag.StringVar(&filePath, "f", "", "source file path")
	// flag.StringVar(&dir, "d", "", "source file dir path")
	flag.StringVar(&target, "t", "proto", "proto file target path")
	flag.Usage = usage
}*/

func main() {
	//flag.Parse()
	//if filePath == "" {
	//	flag.Usage()
	//	return
	//}

	fmt.Println("开始生成proto,Go RootPath:", RootPath)

	cloud = Cloud{}
	cloud.Name = CloudName
	cloud.ImportDtos = map[string]*Dto{}

	// 处理poes;
	processPo()

	var serviceNames []string
	// 表示指定服务;
	if ServiceName != "" {
		serviceNames = append(serviceNames, ServiceName)
	} else {
		// 如果不指定服务名称,则将整个Cloud的服务全部生成;
		serviceNames = getServiceNames(cloud.Name)
	}
	for _, serviceName := range serviceNames {
		service := &Service{}
		service.RootDto = &Dto{IsRoot: true}
		service.RootDto.Create()
		cloud.Services = append(cloud.Services, service)

		basePath := RootPath + ServiceRoot + cloud.Name + "\\" + serviceName
		processDto(basePath+"\\dto\\d.go", service, service.RootDto)
		processDomain(basePath+"\\domain\\domain.go", service)
	}

	for _, service := range cloud.Services {
		saveToFile(service, service.Name+".proto")
	}
	fmt.Println("proto生成完成!")
}

// 解析领域模型,获取方法和引用的实体;
func processDomain(path string, service *Service) {
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	ast.Inspect(f, func(node ast.Node) bool {
		if node == nil {
			return true
		}
		// domain
		if typeSpecNode, ok := node.(*ast.TypeSpec); ok {
			if _, f := typeSpecNode.Type.(*ast.StructType); f {
				structName := typeSpecNode.Name.Name
				if !strings.HasSuffix(structName, "Domain") {
					return true
				}
				lastIndex := strings.LastIndex(structName, "Domain")
				service.Name = structName[0:lastIndex] + "Service"
				fmt.Println("接口名称：", service.Name)
			}
		}

		// 处理domain的func
		if funcDeclNode, ok := node.(*ast.FuncDecl); ok {
			// 只处理长度为1的receiver
			if len(funcDeclNode.Recv.List) == 1 {
				recv := funcDeclNode.Recv.List[0]
				if len(recv.Names) == 1 {

					serviceFunction := &ServiceFunction{}
					serviceFunction.Name = funcDeclNode.Name.Name
					service.ServiceFunctions = append(service.ServiceFunctions, serviceFunction)

					if funcDeclNode.Doc.List != nil {
						for i, doc := range funcDeclNode.Doc.List {
							serviceFunction.Comment += doc.Text
							if i < len(funcDeclNode.Doc.List)-1 {
								serviceFunction.Comment += "\n"
							}
						}
					}
					funcType := funcDeclNode.Type
					// 解析参数列表
					for _, param := range funcType.Params.List {
						////获取参数名称
						//for _, paramName := range param.Names {
						//
						//}
						ptype := getPointedType(param.Type)

						// 获取参数类型,如果是简单类型直接输出;
						if paramType, ok := ptype.(*ast.Ident); ok {
							serviceFunction.ParamTypes = append(serviceFunction.ParamTypes, paramType.Name)
						}
						// 如果是类;
						if paramType, ok := ptype.(*ast.SelectorExpr); ok {
							serviceFunction.ParamTypes = append(serviceFunction.ParamTypes, paramType.Sel.Name)
						}
					}
					// 解析返回值
					for _, result := range funcType.Results.List {
						rtype := getPointedType(result.Type)
						// 获取参数类型,如果是简单类型直接输出;
						if resultType, ok := rtype.(*ast.Ident); ok {
							serviceFunction.ResultTypes = append(serviceFunction.ResultTypes, resultType.Name)
						}
						// 如果是外部类;
						if resultType, ok := rtype.(*ast.SelectorExpr); ok {
							// 获取前缀;
							selectorX := fmt.Sprintf("%v", resultType.X)
							// 表示类型是Poes;
							if selectorX == "poes" {
								messageName := resultType.Sel.Name
								// 如果当前po没有加入到消息列表中，则添加;
								if message, ok := cloud.Poes.Messages[messageName]; !ok {
									log.Fatalf("service[%v]返回参数[%v]并未在poes文件中定义", service.Name, messageName)
								} else {
									service.RootDto.Messages[messageName] = message
								}
							}
							serviceFunction.ResultTypes = append(serviceFunction.ResultTypes, resultType.Sel.Name)
						}
					}
				}
			}
		}

		return true
	})
}

// 解析领域模型,获取引用的数据库实体;
func processPo() {
	dirPath := RootPath + ServiceRoot + cloud.Name + "\\poes"
	if !isDirExist(dirPath) {
		log.Println("读取无法%v库的poes, 文件夹不存在.", cloud.Name)
		return
	}
	fmt.Println("开始读取poes,路径:", ServiceRoot+cloud.Name+"\\poes")

	cloud.Poes = &Dto{}
	cloud.Poes.Create()
	cloud.Poes.Name = "POES"

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, dirPath+"\\p.go", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	ast.Inspect(f, func(node ast.Node) bool {
		if node == nil {
			return true
		}
		if typeSpecNode, ok := node.(*ast.TypeSpec); ok {
			// 处理结构体
			if structNode, f := typeSpecNode.Type.(*ast.StructType); f {
				structName := typeSpecNode.Name.Name
				log.Println("struct名称：", structName)

				messageFields := structParser(cloud.Poes, structName, structNode)
				message := &Message{}
				message.Name = structName
				message.Type = POES
				message.MessageFields = messageFields
				cloud.Poes.Messages[message.Name] = message
			}
		}
		return true
	})

	fmt.Println("poes生成完成!")
}

func processDto(path string, service *Service, dto *Dto) {
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	ast.Inspect(f, func(node ast.Node) bool {
		if node == nil {
			return true
		}
		// 处理包;
		if importSpec, ok := node.(*ast.ImportSpec); ok {
			// 获取路径时,删除前后引号,并转换为小写;
			packagePath := strings.ToLower(importSpec.Path.Value[1 : len(importSpec.Path.Value)-1])
			// 包别名,若是没有,则取当前包的名称(最后的一个/名称)
			packageAlias := ""
			if importSpec.Name != nil {
				packageAlias = importSpec.Name.Name
			} else {
				pathSplit := strings.Split(packagePath, "/")
				packageAlias = pathSplit[len(pathSplit)-1]
			}
			dto.Imports[packageAlias] = &Import{
				PackageName: packageAlias,
				PackagePath: packagePath,
			}

			// 判断当前包是否已经被解析过了;
			if _, ok := cloud.ImportDtos[packagePath]; !ok {
				importDto := &Dto{}
				importDto.Create()
				cloud.ImportDtos[packagePath] = importDto

				// 递归,处理嵌入的dto包;
				importDtoPath := RootPath + strings.Replace(packagePath, "/", "\\", -1) + "\\d.go"
				processDto(importDtoPath, service, importDto)
			}
		}

		if typeSpecNode, ok := node.(*ast.TypeSpec); ok {
			// 处理结构体
			if structNode, f := typeSpecNode.Type.(*ast.StructType); f {
				structName := typeSpecNode.Name.Name
				log.Println("struct名称：", structName)

				messageFields := structParser(dto, structName, structNode)
				dto.MessageIndex++
				message := &Message{}
				message.Name = structName
				if strings.HasSuffix(structName, "Request") {
					message.Type = REQ
				} else if strings.HasSuffix(structName, "Response") {
					message.Type = RESP
				} else {
					message.Type = ELSE
				}
				message.Index = dto.MessageIndex
				message.MessageFields = messageFields
				dto.Messages[message.Name] = message
			}
		}

		return true
	})
}

/*
解析结构体,获取其字段;
可解析嵌入体,例如:
type User struct{
	common.UserInfo // 嵌入common包中的UserInfo
	Audit //嵌入本包中的Audit struct
	UserId int64
}
*/
func structParser(dto *Dto, structName string, structNode *ast.StructType) []MessageField {
	messageFields := []MessageField{}
	index := 0
	for _, field := range structNode.Fields.List {
		// 如果是嵌入体;
		if field.Names == nil {
			var fields []MessageField
			// 嵌入一个实体;
			// 表示嵌入当前文件夹的实体;
			if fieldType, ok := field.Type.(*ast.Ident); ok {
				if message, ok := dto.Messages[fieldType.Name]; ok {
					fields = message.MessageFields
				} else {
					log.Fatalf("在dto[%v]中,解析%v结构时,未找到%v结构的定义.嵌入结构必须在使用前声明.", dto.Name, structName, fieldType.Name)
				}
			}

			// 表示嵌入外部包的类;
			if fieldType, ok := field.Type.(*ast.SelectorExpr); ok {
				// 通过名称或别名获取路径;
				imp := dto.Imports[fmt.Sprintf("%v", fieldType.X)]
				// 通过包路径获取dto中的消息集合;
				messages := cloud.ImportDtos[imp.PackagePath].Messages
				// 通过实体名称,获取字段;
				fields = messages[fieldType.Sel.Name].MessageFields
			}

			for _, field := range fields {
				index++
				messageField := MessageField{}
				messageField.Index = index
				messageField.FieldName = field.FieldName
				messageField.FieldType = field.FieldType
				messageField.Comment = field.Comment

				messageFields = append(messageFields, messageField)
			}
			continue
		}

		// 如果是类似name,address string 这样的定义则报错
		if len(field.Names) > 1 {
			log.Fatalf("struct %v error,字段不允许`name,address string`这样的定义", structName)
		}

		index++
		messageField := MessageField{}
		messageField.Index = index
		messageField.FieldName = field.Names[0].Name

		// 若是指针,获取指向的类型;
		ft := getPointedType(field.Type)

		// 基本类型处理
		if fieldType, ok := ft.(*ast.Ident); ok {
			messageField.FieldType = getProtoType(fieldType.Name)
		}

		// map类型处理
		if fieldType, ok := ft.(*ast.MapType); ok {
			key, value := "", ""
			if keyType, ok := fieldType.Key.(*ast.Ident); ok {
				key = getProtoType(keyType.Name)
			}
			if valueType, ok := fieldType.Value.(*ast.Ident); ok {
				value = getProtoType(valueType.Name)
			}
			messageField.FieldType = fmt.Sprintf("map<%v,%v>", key, value)
		}

		// 处理引用类型
		if fieldType, ok := ft.(*ast.SelectorExpr); ok {
			messageField.FieldType = fieldType.Sel.Name
		}

		// 处理参数是数组的情况
		if fieldType, ok := ft.(*ast.ArrayType); ok {
			elt := getPointedType(fieldType.Elt)

			if fieldTypeElt, ok := elt.(*ast.Ident); ok {
				// byte 数组特殊处理 转换成bytes
				if fieldTypeElt.Name == "byte" {
					messageField.FieldType = "bytes"
				} else {
					messageField.FieldType = "repeated " + getProtoType(fieldTypeElt.Name)
				}
			}
		}

		// 获取注释
		if field.Comment != nil {
			messageField.Comment = field.Comment.List[0].Text
		}
		messageFields = append(messageFields, messageField)
	}
	return messageFields
}

func usage() {
	fmt.Fprintf(os.Stderr, `go2proto version: go2proto/1.0.0
Usage: go2proto [-f] [-t]

Options:
`)
	flag.PrintDefaults()
}
