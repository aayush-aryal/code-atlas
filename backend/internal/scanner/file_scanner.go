package scanner

import "os"

func ReadFile(path string)([]byte,error){
	data,err:=os.ReadFile(path)
	if err!=nil{
		return nil,err
	}
	return data,nil

}