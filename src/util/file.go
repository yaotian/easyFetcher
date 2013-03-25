//所有文件的操作
package util

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//将数据读出来，返回数组
func Read(file string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	var lines []string
	for {
		line, _ := r.ReadString('\n')
		if line == "" {
			break
		}
		lines = append(lines, strings.Trim(line, "\r\n"))
	}
	return lines, nil
}

/*
filter 表示哪些内容不写入文件，即删除掉
*/
func Write(file string, lines []string, filter string) error {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, line := range lines {
		if !strings.Contains(line, filter) {
			fmt.Fprintf(f, "%s\n", line)
		}
	}
	return nil
}

func CookieToFile(data []*http.Cookie, fileName string) error {
	fmt.Print("MapToFile:", data)
	file, _ := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0)
	defer file.Close()
	enc := gob.NewEncoder(file)
	err := enc.Encode(data)
	if err != nil {
		fmt.Println("Error in encoding gob")
		return err
	}
	return nil

}

func FileToCookie(fileName string) ([]*http.Cookie, error) {
	file, e := os.OpenFile(fileName, os.O_RDONLY, 0)
	if e != nil {
		return nil, e
	}
	defer file.Close()
	var data []*http.Cookie
	dec := gob.NewDecoder(file)
	err := dec.Decode(&data)
	if err != nil {
		fmt.Println("Error in decoding gob")
	} else {
		return data, nil
	}
	return nil, nil

}

func ConfigToFile(data []string, fileName string) error {
	fmt.Print("ConfigToFile:", data)
	file, _ := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0)
	defer file.Close()
	enc := gob.NewEncoder(file)
	err := enc.Encode(data)
	if err != nil {
		fmt.Println("ConfigToFile,Error in encoding gob")
		return err
	}
	return nil

}

func FileToConfig(fileName string) ([]string, error) {
	file, e := os.OpenFile(fileName, os.O_RDONLY, 0)
	if e != nil {
		return nil, e
	}
	defer file.Close()
	var data []string
	dec := gob.NewDecoder(file)
	err := dec.Decode(&data)
	if err != nil {
		fmt.Println("FileToConfig, Error in decoding gob")
	} else {
		return data, nil
	}
	return nil, nil

}

func GetOneLineAndRemoveIt(file string) string {
	lines, error := Read(file)

	if error != nil || len(lines) == 0 {
		return ""
	}

	if result := lines[0]; result != "" {
		Write(file, lines, result)
	}
	Trace("Got one msg from file:", lines[0])
	return lines[0]
}

func MakeFileWriteChan(filename string) chan string {
	file_chan := make(chan string)
	go func() {
		for {
			content := <-file_chan
			if content != "" {
				Trace("will write :", content)
				file, _ := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
				defer file.Close()
				fmt.Fprintf(file, "%s\n", content)
			}
		}
	}()
	return file_chan
}

func GetRandomMsgFromFolder(folder string) string {
	files := GetFileListFromFolder(folder)
	howManyFile := len(files)
	if howManyFile == 0 {
		return ""
	}
	whichOne := rand.Intn(howManyFile)
	whichfileName := files[whichOne]
	return GetOneLineAndRemoveIt(whichfileName)
}

func GetFileListFromFolder(folder string) []string {
	result := []string{}
	self := F{
		files: make([]*sysFile, 0),
	}
	err := filepath.Walk(folder, func(path string, f os.FileInfo, err error) error {
		return self.visit(path, f, err)
	})

	if err != nil {
		Report_Error("can't get information from folder", err)
		return []string{""}
	}

	for _, v := range self.files {
		if v.fType == IsRegular {
			result = append(result, v.fName)
		}
	}
	return result
}

//==========================================for list file
const (
	IsDirectory = iota
	IsRegular
	IsSymlink
)

type sysFile struct {
	fType  int
	fName  string
	fLink  string
	fSize  int64
	fMtime time.Time
	fPerm  os.FileMode
}

type F struct {
	files []*sysFile
}

func (self *F) visit(path string, f os.FileInfo, err error) error {
	if f == nil {
		return err
	}
	var tp int
	if f.IsDir() {
		tp = IsDirectory
	} else if (f.Mode() & os.ModeSymlink) > 0 {
		tp = IsSymlink
	} else {
		tp = IsRegular
	}
	inoFile := &sysFile{
		fName:  path,
		fType:  tp,
		fPerm:  f.Mode(),
		fMtime: f.ModTime(),
		fSize:  f.Size(),
	}
	self.files = append(self.files, inoFile)
	return nil
}

//======================================End for list file
