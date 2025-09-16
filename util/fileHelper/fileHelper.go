package fileHelper

import (
	"bufio"
	"encoding/xml"
	"game/util"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func exists(path string) (bool, bool) {
	info, err := os.Stat(path) //os.Stat获取文件信息
	if os.IsNotExist(err) {
		return false, false
	}
	return true, info.IsDir()
}

// 文件是否存在
func ExistFile(path string) bool {
	exist, isDir := exists(path)
	if exist && isDir == false {
		return true
	}
	return false
}

// 目录是否存在
func ExistDirectory(path string) bool {
	exist, isDir := exists(path)
	if exist && isDir {
		return true
	}
	return false
}

type DirectoryInfo struct {
	Directories []string //当前目录下目录路径集合
	Files       []string //当前目录下文件路径集合
}

// 获取目录下的所有文件和目录路径
func GetDirectoryInfo(dirPath string) (DirectoryInfo, error) {
	var currentDirs []string
	var currentFiles []string
	abs, err := filepath.Abs(dirPath)

	if err != nil {
		return DirectoryInfo{}, nil
	}
	err = filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		absPath := filepath.Join(abs, path)
		if info.IsDir() {
			currentDirs = append(currentDirs, absPath)
		} else {
			currentFiles = append(currentFiles, absPath)
		}
		return err
	})
	if err != nil {
		return DirectoryInfo{}, nil
	}
	return DirectoryInfo{Directories: currentDirs, Files: currentFiles}, nil
}

// 获取目录下的所有文件
func GetFiles(dirPath string) []string {
	var files []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return err
	})
	if err != nil {
		util.LogError(err)
	}
	return files
}

// 目录不存在，则创建
func CreateDirectory(dirPath string) {
	if !ExistDirectory(dirPath) {
		err := os.Mkdir(dirPath, 0755) //create a directory and give it required permissions
		if err != nil {
			util.LogError(err) //print the error on the console
			return
		}
	}
}

// 获取父目录完整路径
func GetParentDirectoryPath(filePath string) string {
	abs, _ := filepath.Abs(filePath)
	index := strings.LastIndex(abs, "\\")
	return abs[0:index]
}

// 写入文件
func WriteLines(fileName string, lines []string) {
	if !ExistFile(fileName) {
		pDirName := GetParentDirectoryPath(fileName)
		CreateDirectory(pDirName)
	}

	file, err := os.Create(fileName)
	if err != nil {
		util.LogError("Error creating file:", err)
	}
	defer file.Close()

	var newline string
	if runtime.GOOS == "windows" {
		newline = "\r\n"
	} else {
		newline = "\n"
	}

	result := ""
	for i := 0; i < len(lines); i++ {
		result += lines[i] + newline
	}
	_, err = file.WriteString(result)
	if err != nil {
		util.LogError("Error writing to file:", err)
		return
	}
}

// 读取文件
func ReadLines(fileName string) []string {
	if !ExistFile(fileName) {
		util.LogError("文件路径不存在:", fileName)
		return make([]string, 0)
	}
	f, err := os.Open(fileName)
	if err != nil {
		util.LogError("err:", err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	lines := make([]string, 0)
	// 按行读取
	for {
		line, _, err := reader.ReadLine()
		lines = append(lines, string(line))
		if err == io.EOF {
			break
		} else if err != nil {
			util.LogError("err:", err)
			break
		}
	}
	return lines
}

func ReadXML(inputxml string) {
	inputReader := strings.NewReader(inputxml)
	p := xml.NewDecoder(inputReader)

	for t, err := p.Token(); err == nil; t, err = p.Token() {
		switch token := t.(type) {
		case xml.StartElement:
			name := token.Name.Local
			util.Log("Token name: %s\n", name)
			for _, attr := range token.Attr {
				attrName := attr.Name.Local
				attrValue := attr.Value
				util.Log("An attribute is: %s %s\n", attrName, attrValue)
			}
		case xml.EndElement:
			util.Log("End of token")
		case xml.CharData:
			content := string([]byte(token))
			util.Log("This is the content: %v\n", content)
		default:
		}
	}
}
func ReadFile(fileName string) string {
	inputFile, inputError := os.Open(fileName)
	if inputError != nil {
		util.LogError(inputError)
		return ""
	}
	defer inputFile.Close()

	inputReader := bufio.NewReader(inputFile)
	result := ""
	for {
		inputString, readerError := inputReader.ReadString('\n')
		result += inputString
		if readerError == io.EOF {
			return result
		}
	}
}
func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}
