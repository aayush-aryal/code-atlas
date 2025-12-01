package scanner

import (
	"io/fs"
	"path/filepath"
)

var allowedExtensions=map[string]bool{
	".go":true,
	".c":true,
	".js":true,
	".ts":true,
	".jsx":true,
	".tsx":true,
}

var skipDir=map[string]bool{
	".git":true,
	"node_modules":true,
	"dist":true,
	"build":true,
}

func ScanDirectory(root string)([]string, error){
	var files []string
	err:=filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error{
		if err!=nil{
			return err
		}

		// skip directories 
		if d.IsDir()&& skipDir[d.Name()]{
			return filepath.SkipDir
		}

		if d.IsDir(){
			return nil 
		}

		ext:=filepath.Ext(path)
		if !allowedExtensions[ext]{
			return nil
		}
		files=append(files, path)
		return nil
	})

	if err!=nil{
		return nil,err
	}

	return files, nil
}

