package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/aayush-aryal/code-atlas/internal/ai"
	"github.com/aayush-aryal/code-atlas/internal/codebase"
	"github.com/aayush-aryal/code-atlas/internal/scanner"
)

type resp struct{
	Status int 
	Message string 
} 

type AskRequest struct{
	Question string `json:"question"`
	Function string `json:"function"`
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

var currentProject *codebase.Project


func main(){
	
	fmt.Println("CodeAtlas in running")
	var err error 
	
	currentProject, err=codebase.Analyze(".")
	if err!=nil{
		log.Fatal("Failed to analyze codebase: ",err)
	}
	fmt.Printf("Analysis complete! Found %d files.\n", len(currentProject.Graph))


	http.HandleFunc("/hello",hello)
	http.HandleFunc("/headers",headers)
	http.HandleFunc("/health",health)
	http.HandleFunc("/file",handleFileContent)
	http.HandleFunc("/codebase_analyze",handleAnalyze)
	http.HandleFunc("/context", handleGetContext)
	http.HandleFunc("/ask", handleAsk)
	http.HandleFunc("/visual",handleGetGraph)
	http.HandleFunc("/func", handleGetFunc)
	http.ListenAndServe(":8090",nil)
}


func handleAnalyze(w http.ResponseWriter, req *http.Request){
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Content-Type","application/json")
	
	if err:=json.NewEncoder(w).Encode(currentProject);err!=nil{
		http.Error(w,"Failed to encode JSON", http.StatusInternalServerError)

	}
}

func handleFileContent(w http.ResponseWriter, req *http.Request){
	w.Header().Set("Accesss-Control-Allow-Origin","*")
	path:=req.URL.Query().Get("path")
	if path==""{
		http.Error(w,"Missing path", http.StatusBadRequest)
		return 
	}

	content,err:=scanner.ReadFile(path)
	if err!=nil{
		http.Error(w,"File not found", http.StatusNotFound)
		return
	}
	w.Write(content)
}

func handleGetContext(w http.ResponseWriter, req *http.Request){
	w.Header().Set("Accesss-Control-Allow-Origin","*")
	functionName:=req.URL.Query().Get("func")
	depth:=req.URL.Query().Get("depth")
	contextDepth:=0
	if functionName==""{
		http.Error(w,"Missing function name",http.StatusBadRequest)
		return
	}
	if depth==""{
		contextDepth=1
	}else{
		var err error
		contextDepth, err = strconv.Atoi(depth)
		if err != nil {
            http.Error(w, "Invalid depth (must be a number)", http.StatusBadRequest)
            return
        }
	}
	if currentProject==nil{
		http.Error(w, "Server not ready (Analysis missing)", http.StatusServiceUnavailable)
		return
	}
	context,err:=currentProject.GetContext(functionName,contextDepth)
	if err!=nil{
		http.Error(w, fmt.Sprintf("Error generating context: %v", err), http.StatusInternalServerError)
        return
	}
	w.Header().Set("Content-Type", "text/plain")
    w.Write([]byte(context))
}

func handleAsk(w http.ResponseWriter, req *http.Request){
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.Header().Set("Content-Type", "application/json")

	if req.Method == "OPTIONS" {
        return
    }

	reqBody:=req.Body
	if reqBody==nil{
		http.Error(w,"Invalid request",http.StatusBadRequest)
	}
	var r AskRequest
	if err:=json.NewDecoder(req.Body).Decode(&r);err!=nil{
		http.Error(w,"invalid JSON body", http.StatusBadRequest)
		return
	}
	if currentProject==nil{
		http.Error(w,"Invalid request", http.StatusServiceUnavailable)
		return
	}
	context,err:=currentProject.GetContext(r.Function,2)
	if err!=nil{
		http.Error(w,"Could not find the funcition", http.StatusBadRequest)
		return
	}
	systemPrompt := fmt.Sprintf(`
	You are an expert Senior Go Engineer. 
	Answer the user's question based ONLY on the code context provided below.

	--- BEGIN CODE CONTEXT ---
	%s
	--- END CODE CONTEXT ---

	Question: %s
	`, context, r.Question)
	answer,err:=ai.Ask(systemPrompt)
	if err!=nil{
		http.Error(w,"Ollama Error"+err.Error(), http.StatusInternalServerError)
		return 
	}

	response:=map[string]string{
		"answer":answer,
	}

	json.NewEncoder(w).Encode(response)
}

func handleGetGraph(w http.ResponseWriter, req *http.Request){
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Content-Type", "application/json")
	if currentProject==nil{
		http.Error(w,"No project", http.StatusServiceUnavailable)
	}
	visualGraph:=currentProject.ComputeVisualGraph();
	if err:=json.NewEncoder(w).Encode(visualGraph);err!=nil{
		http.Error(w,"Something went wrong",http.StatusInternalServerError)
	}
}

func handleGetFunc(w http.ResponseWriter, req *http.Request){
	w.Header().Set("Content-Type","text/plain")
	if currentProject==nil{
		http.Error(w,"No project", http.StatusServiceUnavailable)
	}
	funcName:=req.URL.Query().Get("functionName")
	funcPath:=currentProject.FunctionTable[funcName]
	if (len(funcPath)>0){
		funcText,err:=codebase.ReadFunction(funcPath[0])
	if err!=nil{
		http.Error(w,"Could not find the function",http.StatusNotFound)
	}
	fmt.Fprint(w,funcText)
	}
	
}