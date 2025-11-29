package scanner

import (
	"io/fs"
	"path/filepath"
)

func ScanDirectory(root string)([]string, error){
	var files []string
	err:=filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error{
		if err!=nil{
			return err
		}
		files=append(files,path)
		return nil
	})

	if (err!=nil){
		return nil,err
	}

	return files, nil
}

