package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aayush-aryal/code-atlas/internal/parser"
	"github.com/aayush-aryal/code-atlas/internal/scanner"
)

type resp struct{
	Status int 
	Message string 
} 

func hello(w http.ResponseWriter, req *http.Request){
	fmt.Fprint(w,"hello\n")
}

func headers(w http.ResponseWriter, req *http.Request){
	for name, headers:=range req.Header {
		for _,h:=range headers {
			fmt.Fprintf(w,"%v: %v\n", name,h)
		}
	}
}

func health(w http.ResponseWriter, req *http.Request){
	response:= resp{
		Status: 200,
		Message:"good health nener ne",
	}
	res,err:=json.Marshal(response)

	if (err!=nil){
		panic(err)
	}
	jsonString:=string(res)
	fmt.Fprint(w,jsonString)

}


func main(){
	files,err:=scanner.ScanDirectory(".")
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println(files)
	for _,value:=range files{
		metadata,err:=scanner.ExtractMetaData(value)

		if err!=nil{
			fmt.Println("Iono")
		}
		file_data,err:=scanner.ReadFile(value)
		if err!=nil{
			fmt.Println("Iono")
		}
		// use this to print root
		root_node,err:=parser.ParseFile(file_data)
		if err!=nil{
			fmt.Println("How do i fix a billion errors")
		}
		names,err:=parser.ExtractFunctionNames(root_node,file_data)
		if err!=nil{
			fmt.Println("Error")
		}
		fmt.Println(names)
		fmt.Println(metadata)
	}
	http.HandleFunc("/hello",hello)
	http.HandleFunc("/headers",headers)
	http.HandleFunc("/health",health)
	http.ListenAndServe(":8090",nil)
}

