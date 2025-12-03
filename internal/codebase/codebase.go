package codebase

import (
	"log"
	"sync"

	parser "github.com/aayush-aryal/code-atlas/internal/parser"
	"github.com/aayush-aryal/code-atlas/internal/scanner"
)

type FileNode struct {
	Path string `json:"name"`
	MetaData scanner.MetaData `json:"metadata"`
	Functions []parser.FunctionInfo `json:"functions"`
	Imports []parser.ImportDetail `json:"imports"`
}


type Graph map[string]FileNode

type Project struct{
    Graph Graph `json:"graph"`
    FunctionTable FunctionTable `json:"functionTable"`
}

func Analyze(dir string)(*Project,error){
   // scan codebase for files
   graph:=make(Graph)


   files,err:=scanner.ScanDirectory(dir)
   if err!=nil{
       return nil,err
   }


   jobs:=make(chan string, len(files))
   results:=make(chan FileNode, len(files))


   var wg sync.WaitGroup


   // there are 8 workers who can perform the func worker using go
   // go causes the use of multithreading
   numWorkers:=8
   for i:=0; i<numWorkers;i++{
       wg.Add(1)
       go worker(jobs,results,&wg)
   }


   // after that you add all file path to the list of jobs you need to do
   // ie all file path needs a filenode so add it
   for _,file:=range files{
       jobs<- file
   }


   // close the jobs channel after the function runs
   // ie no more jobs to send
   close(jobs)


   // closer routine
   // wait for workgroup until results are done
   // after that close result
   go func(){
       wg.Wait()
       close(results)
   }()
   // take your results channel where each job gives a file path
   // use that and add to your dict


   for node:=range results{
       graph[node.Path]=node
   }
   functionTable:=MapFunctionToImport(graph)
   project:=Project{
    Graph: graph,
    FunctionTable: functionTable,
   }
   return &project,nil


}


func worker(jobs <-chan string, results chan<-FileNode, wg *sync.WaitGroup){
   defer wg.Done()
   for filePath:=range jobs{
       metadata,err:=scanner.ExtractMetaData(filePath)
       if err!=nil{
           log.Printf("Could not get metadata for %s",filePath)
       }

       // read file
       readFile,err:=scanner.ReadFile(filePath)
       if err!=nil{
           log.Printf("Could not read file %s",filePath)
           continue
       }
	   	   	//parse the file and get root node first
		rootNode,err:=parser.ParseFile(readFile)
		if err!=nil{continue}	
       // extract functions
       functions,err:=parser.ExtractFunctionNames(rootNode,readFile)
	   if err!=nil{continue}
	   imports,err:=parser.ExtractImports(rootNode,readFile)
	   if err!=nil{
		log.Printf("Error while parsing imports %v",err)
	   }
       if err!=nil{
           log.Printf("Could not extract functions for %s", filePath)
           continue
       }


       results<- FileNode{
           Path:filePath,
           MetaData: metadata,
           Functions: functions, 
		   Imports: imports, 
       }
   }
}


      

		

	
