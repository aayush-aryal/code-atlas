package codebase

import (
	"fmt"
	"strings"

	"github.com/aayush-aryal/code-atlas/internal/scanner"
)

// go through graph

type FunctionRef struct{
	Path string `json:"path"`
	StartLine int `json:"startLine"`
	EndLine int `json:"endLine"`
}
type FunctionTable map[string][]FunctionRef

func MapFunctionToImport(graph Graph) FunctionTable{
	functionTable:=make(FunctionTable)

	for filePath,fileNode:=range graph{
		// a filenode consists of functions map out each function name to its file path(name param)
		for _,function:=range fileNode.Functions{
			ref:=FunctionRef{
				Path:filePath,
				StartLine: function.StartLine,
				EndLine: function.EndLine,
			}
			functionTable[function.Name]=append(functionTable[function.Name],ref)
		}
	}
	return functionTable
}

// reciever is used to connect stucts and data
// since the GetContext needs access to project we pass it as reciever

func (p *Project) GetContext(functionName string,depth int)(string,error){
	visited:=make(map[string]bool)
	var sb strings.Builder
	
	err:=p.RecursiveContext(functionName , 0 ,depth , visited , &sb)
	if err!=nil{
		return "",err
	}
	return sb.String(),nil
}


func (p*Project) RecursiveContext(funcName string, currDepth int,maxDepth int, visited map[string]bool,sb *strings.Builder)(error){
	if visited[funcName]{
		return nil
	}

	if currDepth>maxDepth{
		return nil
	}

	visited[funcName]=true

	refs,ok:=p.FunctionTable[funcName]
	if (!ok || len(refs)==0){
		return nil
	}
	targetFunc:=refs[0]
	code,err:=ReadFunction(targetFunc)
	if err!=nil{
		return  fmt.Errorf("failed to read %s: %v",funcName,err)
	}
	// E. Format Output
	indent := strings.Repeat("  ", currDepth) // Visual indentation
	sb.WriteString(fmt.Sprintf("%s--- DEPTH %d: %s ---\n%s\n\n", indent, currDepth, funcName, code))
	fileNode,ok:=p.Graph[targetFunc.Path]
	if !ok{return nil}
	// go through all functions in the file
	for _,function:=range fileNode.Functions{
		// check the func we have called on
		if function.Name==funcName{
			// call recursive context on its calls
			for _,call:=range function.Calls{
				if !visited[call]{
					
					err:=p.RecursiveContext(call,currDepth+1,maxDepth,visited,sb)
					if err!=nil{
						return err
					}
				}
			}
			break
		}
		

	}
	return nil
}

func ReadFunction(ref FunctionRef)(string,error){
	content,err:=scanner.ReadFile(ref.Path)
	if err!=nil{return "",err}
	lines:=strings.Split(string(content),"\n")
	if ref.StartLine>=len(lines)|| ref.EndLine>len(lines){
		return "",fmt.Errorf("lines out of range")
	}
	return strings.Join(lines[ref.StartLine:ref.EndLine+1],"\n"),nil
}


