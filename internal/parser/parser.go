package parser

import (
	"context"
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)


type FunctionInfo struct {
	Name string 
	Parameters string
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

func ExtractFunctionNames(rootNode *sitter.Node, file[]byte)([]string,error){
	var (query_string=`
	(function_declaration
		name: (identifier) @function-name
		parameters: (parameter_list) @parameter-list
	)
	`
)
	query,err:=sitter.NewQuery([]byte(query_string),golang.GetLanguage())
	if err!=nil{
		return nil,err
	}
	defer query.Close()

	cursor:=sitter.NewQueryCursor()
	defer cursor.Close()

	cursor.Exec(query,rootNode)
	var functionNames []string 

	for{
		match,ok:=cursor.NextMatch()
		if !ok{
			break
		}
		var currentFuncName string 
		var currentParams string
		for _, capture:=range match.Captures{
			captureName:=query.CaptureNameForId(capture.Index)
			content:=string(file[capture.Node.StartByte():capture.Node.EndByte()])

			switch captureName {
			case "function-name":
				currentFuncName=content 
			case "parameter-list":
				currentParams=content
			}

		}
		fullSignature := fmt.Sprintf("%s%s", currentFuncName, currentParams)
        functionNames = append(functionNames, fullSignature)
		
	}
	return functionNames,nil
}