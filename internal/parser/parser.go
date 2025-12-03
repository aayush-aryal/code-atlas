package parser

import (
	"context"
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)


type FunctionInfo struct {
	Name string `json:"name"`
	Parameters string `json:"parameters"`
	Calls []string `json:"calls"`
	StartLine int `json:"startLine"`
	EndLine int `json:"endLine"`
}


type ImportDetail struct{
	Name string `json:"name"`
	Path string `json:"path"`
}

func ParseFile(file[]byte)(*sitter.Node,error){
	parser:=sitter.NewParser()

	// get the language object
	lang:=golang.GetLanguage()
	parser.SetLanguage(lang)

	tree,err:=parser.ParseCtx(context.Background(),nil,file)
	if err!=nil{
		return nil,err
	}
	return tree.RootNode(),nil

}

func ExtractFunctionNames(rootNode *sitter.Node,file []byte) ([]FunctionInfo, error) {
	queryString := `
	(function_declaration
		name: (identifier) @function-name
		parameters: (parameter_list) @parameter-list
		body: (block) @function-body
	)@function-root
	`

	query, err := sitter.NewQuery([]byte(queryString), golang.GetLanguage())
	if err != nil {
		return nil, err
	}
	defer query.Close()

	cursor := sitter.NewQueryCursor()
	defer cursor.Close()
	cursor.Exec(query, rootNode)

	var functions []FunctionInfo

	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		var currentFunc FunctionInfo
		var currentBodyNode *sitter.Node


		for _, capture := range match.Captures {
			captureName := query.CaptureNameForId(capture.Index)
			raw := file[capture.Node.StartByte():capture.Node.EndByte()]
			content := string(raw)

			switch captureName {
			case "function-name":
				currentFunc.Name = content
			case "parameter-list":
				currentFunc.Parameters = content
			case "function-body":
				currentBodyNode = capture.Node
			case "function-root":
				currentFunc.StartLine= int(capture.Node.StartPoint().Row)
				currentFunc.EndLine=int(capture.Node.EndPoint().Row)

			}
		}

		if currentBodyNode != nil {
			currentFunc.Calls = ExtractCalls(currentBodyNode, file)
		}

		functions = append(functions, currentFunc)
	}

	return functions, nil
}


func ExtractCalls(node *sitter.Node, file[]byte)([]string){
	queryStr:=`(call_expression function:[
		(identifier) @call 
		(selector_expression field: (field_identifier) @call)
	])
	`

	query,err:=sitter.NewQuery([]byte(queryStr),golang.GetLanguage())
	if err != nil {
		fmt.Printf("Query compilation failed: %v\n", err)
		return nil
	}
	
	defer query.Close()
	cursor:=sitter.NewQueryCursor()
	defer cursor.Close()
	cursor.Exec(query,node)

	var calls []string 
	for{
		match,ok:=cursor.NextMatch()
		if !ok{break}
		for _,capture:=range match.Captures{
			callName:=string(file[capture.Node.StartByte():capture.Node.EndByte()])
			calls=append(calls,callName)
		}
	}
	return  calls
}

func ExtractImports(node *sitter.Node,file []byte)([]ImportDetail,error){
	queryStr:=`
		(import_spec
			name: (package_identifier)? @import-name 
			path:(interpreted_string_literal) @import-path	
)
	`
	query,err:=sitter.NewQuery([]byte(queryStr),golang.GetLanguage())
	if err!=nil{
		fmt.Printf(`Error occured while parsing query %v`,err)
		return nil,err
	}
	defer query.Close()
	cursor:=sitter.NewQueryCursor()
	defer cursor.Close()
	cursor.Exec(query,node)
	
	var imports []ImportDetail
	for{
		match,ok:=cursor.NextMatch()
		if !ok{
			break
		}
		var importPath string 
		var importName string
		for _,capture:=range match.Captures{
			captureName:=query.CaptureNameForId(capture.Index)
			raw:=file[capture.Node.StartByte():capture.Node.EndByte()]
			content:=string(raw)
			switch captureName{
			case "import-path":
				importPath=content 
			case "import-name":
				importName=content
			}
		}
		currImpDetail:=ImportDetail{
			Name: importName,
			Path: importPath,
		}
		imports=append(imports, currImpDetail)
	}
	return imports,nil
}