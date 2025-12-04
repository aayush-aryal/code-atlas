package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OllamaRequest struct {
	Model string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool `json:"stream"`

}

type OllamaResponse struct{
	Response string `json:"response"`
}

func Ask(prompt string)(string,error){
	request:=OllamaRequest{
		Model: "phi3.5:latest",
		Prompt: prompt,
		Stream: false,
	}

	jsonData,err:=json.Marshal(request)
	if err!=nil{
		return "",err
	}

	resp,err:=http.Post(`http://localhost:11434/api/generate`,"application/json",bytes.NewBuffer(jsonData))

	if err!=nil{
		return "",fmt.Errorf("ollama connected failed: %v",err)
	}

	defer resp.Body.Close()
	body,_:=io.ReadAll(resp.Body)
	var ollamaResp OllamaResponse
	if err:=json.Unmarshal(body,&ollamaResp);err!=nil{
		return "", fmt.Errorf("failed to decode response: %v",err)
	}

	return ollamaResp.Response,nil
}