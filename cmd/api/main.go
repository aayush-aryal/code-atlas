package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	files,err:=scanner.ScanDirectory()
	http.HandleFunc("/hello",hello)
	http.HandleFunc("/headers",headers)
	http.HandleFunc("/health",health)
	http.ListenAndServe(":8090",nil)
}

