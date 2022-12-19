package gofolder

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Copy("../dir1", "../dir2")
func Copy(srcPath, genPath string) error {
	err := filepath.Walk(srcPath, func(path string, file os.FileInfo, err error) error {
		if file == nil {
			return err
		}
		relPath, err := filepath.Rel(srcPath, path) // relative path
		if err != nil {
			return err
		}
		newPath := filepath.Join(genPath, relPath)
		if file.IsDir() { // create null directory
			err = os.MkdirAll(newPath, 0777)
			if err != nil {
				return err
			}
		} else {
			err = os.MkdirAll(filepath.Dir(newPath), 0777) // create directory for file
			if err != nil {
				return err
			}
			fileData, err := ioutil.ReadFile(path) // read file
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(newPath, fileData, 0666) // write file
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func Remove(path string) error {
	return os.RemoveAll(path)
}

// Zip("../dir2/xxxx", "../dir1/xxxx.zip")
func Zip(srcPath, genPath string) error {
	zos, err := os.Create(genPath)
	if err != nil {
		return err
	}
	zipWriter := zip.NewWriter(zos)
	defer zipWriter.Close()

	absPath := filepath.Dir(srcPath)
	err = filepath.Walk(srcPath, func(path string, file os.FileInfo, err error) error {
		if file == nil {
			return err
		}
		filepath.Join()
		newPath, err := filepath.Rel(absPath, path)
		if err != nil {
			return err
		}
		if file.IsDir() { // create null directory
			_, err = zipWriter.Create(newPath + "/")
			if err != nil {
				return err
			}
		} else {
			fileData, err := ioutil.ReadFile(path) //create file
			if err != nil {
				return err
			}
			newFile, err := zipWriter.Create(newPath)
			if err != nil {
				return err
			}
			_, err = newFile.Write(fileData)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// UnZip("../dir1/xxxx.zip", "../dir3/xxxx")
func UnZip(srcPath, genPath string) (err error) {
	zipReader, err := zip.OpenReader(srcPath)
	if err != nil {
		return err
	}
	for _, file := range zipReader.File {
		newFilePath := filepath.Join(genPath, file.Name)
		if file.Mode().IsDir() { //create null directory
			err = os.MkdirAll(newFilePath, 0777)
			if err != nil {
				return err
			}
			continue
		} else {
			err = os.MkdirAll(filepath.Dir(newFilePath), 0777) //create directory
			if err != nil {
				return err
			}
			fileReader, err := file.Open() //create file
			if err != nil {
				return err
			}
			newFile, err := os.OpenFile(newFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			_, err = io.Copy(newFile, fileReader)
			fileReader.Close()
			newFile.Close()
			if err != nil {
				return err
			}
		}
	}
	err = zipReader.Close()
	return err
}
