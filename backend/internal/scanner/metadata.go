package scanner

import (
	"os"
	"path/filepath"
	"time"
)

type MetaData struct{
	Path string
	Size int64 
	Extension string 
	LastModified time.Time
}

func ExtractMetaData(path string)(MetaData,error){
	
	fileInfo,err:=os.Stat(path)
	if err!=nil{
		return MetaData{},err
	}

	metadata := MetaData{
		Path:         path,
		Size:         fileInfo.Size(),
		Extension:    filepath.Ext(path),
		LastModified: fileInfo.ModTime(),
	}


	return metadata, nil

}