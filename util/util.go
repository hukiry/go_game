package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

var HUKIRY_TEXT = `
✿       ✿    ✿        ✿    ✿    ✿     ✿✿✿✿    ✿  ✿✿     ✿        ✿    
❀       ❀    ❀        ❀    ❀   ❀        ❀       ❀     ❀    ❀       ❀   
❀       ❀    ❀        ❀    ❀  ❀         ❀       ❀     ❀     ❀     ❀  
❀❀❀ ❀❀❀    ❀        ❀    ❀❀           ❀       ❀  ❀❀        ❀  ❀ 
❀       ❀    ❀        ❀    ❀  ❀         ❀       ❀ ❀             ❀ 
❀       ❀     ❀      ❀     ❀   ❀        ❀       ❀  ❀          ❀ 
✿       ✿       ✿  ✿       ✿    ✿     ✿✿✿✿    ✿     ✿      ✿
`

func GetHUKIRY_TEXT() string {
	return HUKIRY_TEXT
}

func Log(params ...any) {
	str := strings.Fields(time.Now().String())[:2]
	params = append([]any{str}, params...) //插入和追加数据
	_, err := fmt.Fprintln(os.Stdout, params...)
	if err != nil {
		return
	}
}

func LogInput(params ...any) {
	str := strings.Fields(time.Now().String())[:2]
	params = append([]any{str}, params...) //插入和追加数据
	fmt.Print(params...)
}

func LogError(params ...any) {
	//时间戳秒 time.Now().Unix()
	//毫秒 time.Now().UnixMilli()
	str := strings.Fields(time.Now().String())[:2]
	params = append([]any{str}, params...) //插入和追加数据

	params = append(params, tryDebug())

	_, err := fmt.Fprintln(os.Stderr, params...)
	if err != nil {
		return
	}
}

func tryDebug() string {
	result := ""
	for i := 2; i <= 8; i++ {
		_, filename, line, ok := runtime.Caller(i)
		if ok {
			fileName := path.Base(filename)
			if fileName == "main.go" {
				break
			}
			result += fmt.Sprintf("\n	%s:%d", filename, line)
		}
	}
	return result
}

// Handle 闭包函数
func Handle(method func(...any), params ...any) func() {
	return func() {
		method(params...)
	}
}

// CountChineseChars 中文字符计数：一个中文字符2个字节长度，len() 按照字节个数计算
func CountChineseChars(s string) int {
	count := 0
	for _, runeValue := range s {
		if runeValue >= 0x4E00 && runeValue <= 0x9FFF {
			count++
		} else {
			count++
		}
	}
	return count
}

// CountChineseRune 统计字符数个数二
func CountChineseRune(str string) int {
	rt := []rune(str)
	return len(rt)
}

func Md5String(s string) string {
	// 创建 MD5 哈希对象
	hash := md5.New()
	// 写入要计算哈希的数据
	hash.Write([]byte(s))
	// 计算哈希值，返回 []byte
	sum := hash.Sum(nil)
	// 转换为十六进制字符串
	return hex.EncodeToString(sum)
}

func ToString(v interface{}) string {
	if v == nil {
		return "nil"
	}
	// 其他类型使用默认转换
	return fmt.Sprintf("%v", v)
}
