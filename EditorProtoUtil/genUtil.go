package EditorProtoUtil

import (
	"bufio"
	"game/util"
	"game/util/fileHelper"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Member struct {
	name         string //字段名称
	desc         string //注释
	typeName     string //字段类型
	labelNumber  string //标签数字
	labelName    string // 标签名 repeated|optional|required
	isMessage    bool   //是否为消息
	isArray      bool   //是数组
	defaultValue string //默认值
	fullName     string //路径
}

type ProtoBuffer struct {
	name    string   //类和枚举名称
	number  int      //协议号
	desc    string   //类注释
	isEnum  bool     //是否为枚举
	members []Member //类成员
}

type ProtoNumber struct {
	structName string
	number     string
	desc       string
}

// 排序使用
type ProtoNumberSlice []ProtoNumber

func (a ProtoNumberSlice) Len() int { // 重写 Len() 方法
	return len(a)
}
func (a ProtoNumberSlice) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}
func (a ProtoNumberSlice) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return a[j].number > a[i].number
}

// 行分割标记
var indexTag int = 0
var protobufPath string
var exportBinary string
var exportType int

func ExportProtocolBinary() {
	util.Log(`
	 ________________________________帮助说明_________________________________
	|	
	|	[协议模板定义]
	|	//注释
	|	message 消息名
	|	{
 	|	   required|repeated 字段类型 字段名称 = 序号;//注释
	|	   ...
	|	}
	|	生成协议类型 1 = go二进制导出, 2 = lua pb导出
	|_________________________________________________________________________
`)
	readConfigFile()
	if !fileHelper.ExistDirectory(protobufPath) {
		return
	}
	protoNumber := map[string]ProtoNumber{}
	dir, _ := os.ReadDir(protobufPath)
	for _, info := range dir {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".proto") {
			fileName := strings.Split(info.Name(), ".")[0]
			protobufList := readFile(fileName, protobufPath+"/"+info.Name())
			if exportType == 1 {
				writeFileBinary(fileName, protobufList, protoNumber)
			}
		}
	}
	if exportType == 1 {
		writeProtoBinary(protoNumber)
	}
	util.Log("完成！")
}

func readConfigFile() {
	configPath := "EditorProtoUtil/ConfigPath.ini"
	lines := fileHelper.ReadLines(configPath)
	for _, line := range lines {
		if len(line) == 0 || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		arr := strings.Split(strings.TrimSpace(line), "=")
		if len(arr) != 2 {
			util.LogError("数据配置有误：", configPath)
			break
		}

		switch strings.TrimSpace(arr[0]) {
		case "ProtobufPath":
			protobufPath = strings.TrimSpace(arr[1])
			break
		case "ExportBinary":
			exportBinary = strings.TrimSpace(arr[1])
			break
		default:
			resultInt, _ := strconv.Atoi(strings.TrimSpace(arr[1]))
			exportType = resultInt
			break
		}
	}
}

// 写入生成的文件
func writeFileBinary(fileName string, pList []ProtoBuffer, protoNumber map[string]ProtoNumber) {
	fileName = exportBinary + "/" + fileName + "_pb.go"
	list := []string{"package protocol"}
	list = append(list, " ")
	list = append(list, "import \"game/server/packet\"")
	list = append(list, " ")
	pLen := len(pList)
	for i := 0; i < pLen; i++ {
		v := pList[i]
		if v.number > 0 {
			protoNumber[v.name] = ProtoNumber{structName: v.name, number: strconv.Itoa(v.number), desc: v.desc}
		}
		//创建结构体
		list = append(list, "//"+v.desc)
		list = append(list, "type "+v.name+" struct {")
		mlen := len(v.members)
		for i := 0; i < mlen; i++ {
			mValue := v.members[i]
			list = append(list, "\t"+mValue.name+"\t\t"+dealTypeArray(mValue)+"\t\t//"+mValue.desc)
		}
		list = append(list, "}")

		//cmd 方法

		list = append(list, "")
		list = append(list, "func (this *"+v.name+") GetCmd() uint16 {")
		if v.number > 0 {
			list = append(list, "\treturn "+strings.ToUpper(v.name))
		} else {
			list = append(list, "\treturn 0")
		}
		list = append(list, "}")

		list = createWriteText(list, mlen, v)
		list = createReadText(list, mlen, v)
	}
	fileHelper.WriteLines(fileName, list)

	util.Log("生成：" + fileName)
}

// 写入协议号
func writeProtoBinary(protoNumber map[string]ProtoNumber) {
	fileName := exportBinary + "/protocol.go"
	list := []string{"package protocol"}
	list = append(list, `
import (
	"game/server/packet"
	"sync"
)`)
	protoList := []ProtoNumber{}
	for _, value := range protoNumber {
		protoList = append(protoList, value)

	}
	sort.Sort(ProtoNumberSlice(protoList))
	tag := ""
	for i := 0; i < len(protoList); i++ {
		value := protoList[i]
		str := "const " + strings.ToUpper(value.structName) + " uint16 = " + value.number + " // " + value.desc
		if value.number[:2] != tag {
			tag = value.number[:2]
			list = append(list, " ")
		}
		list = append(list, str)
	}

	template := `}

var lock sync.Mutex
func GetMsgPB(cmd uint16) packet.IProto {
	lock.Lock()
	defer lock.Unlock()
	function := msgMap[cmd]
	return function()
}
`
	templateStr := ""
	//协议号类型
	for i := 0; i < len(protoList); i++ {
		value := protoList[i]
		templateStr += "\t" + strings.ToUpper(value.structName) + ": func() packet.IProto { return &" + value.structName + "{} },\n"
	}
	list = append(list, "var msgMap = map[uint16]func() packet.IProto{")
	list = append(list, templateStr)
	list = append(list, template)
	fileHelper.WriteLines(fileName, list)
	util.Log("生成：", fileName)
}

// 创建写入方法
func createWriteText(list []string, mlen int, v ProtoBuffer) []string {
	list = append(list, " ")
	list = append(list, "func (this *"+v.name+") WriteProto(p *packet.Packet) {")
	for i := 0; i < mlen; i++ {
		mValue := v.members[i]
		if mValue.isArray {

			list = append(list, "\tif this."+mValue.name+" == nil {")
			list = append(list, "\t\tthis."+mValue.name+" = make([]"+mValue.typeName+", 0)")
			list = append(list, "\t}")

			list = append(list, "\tp.WriteUInt16(uint16(len(this."+mValue.name+")))")
			list = append(list, "\tfor i := 0; i < len(this."+mValue.name+"); i++ {")
			if mValue.isMessage {
				list = append(list, "\t\tthis."+mValue.name+"[i].WriteProto(p)")
			} else {
				list = append(list, "\t\tp."+iOFuncName(mValue.typeName, true)+"(this."+mValue.name+"[i])")
			}
			list = append(list, "\t}")
		} else {
			if mValue.isMessage {
				list = append(list, "\tthis."+mValue.name+".WriteProto(p)")
			} else {
				list = append(list, "\tp."+iOFuncName(mValue.typeName, true)+"(this."+mValue.name+")")
			}
		}
	}
	list = append(list, "}")
	return list
}

// 创建读取方法
func createReadText(list []string, mlen int, v ProtoBuffer) []string {
	list = append(list, " ")
	list = append(list, "func (this *"+v.name+") ReadProto(p *packet.Packet) {")
	for i := 0; i < mlen; i++ {
		mValue := v.members[i]
		if mValue.isArray {
			list = append(list, "\tthis."+mValue.name+" = make([]"+mValue.typeName+", 0)")
			list = append(list, "\t"+mValue.name+"_len := p.ReadUInt16()")
			list = append(list, "\tfor i := 0; i < int("+mValue.name+"_len); i++ {")
			if mValue.isMessage {
				list = append(list, "\t\t"+mValue.name+"_p := &"+mValue.typeName+"{}")
				list = append(list, "\t\t"+mValue.name+"_p.ReadProto(p)")
				list = append(list, "\t\tthis."+mValue.name+" = append(this."+mValue.name+", *"+mValue.name+"_p)")
			} else {
				list = append(list, "\t\tthis."+mValue.name+" = append(this."+mValue.name+", p."+iOFuncName(mValue.typeName, false)+"())")
			}
			list = append(list, "\t}")
		} else {
			if mValue.isMessage {
				list = append(list, "\tthis."+mValue.name+".ReadProto(p)")
			} else {
				list = append(list, "\tthis."+mValue.name+" = p."+iOFuncName(mValue.typeName, false)+"()")
			}
		}
	}
	list = append(list, "}")
	return list
}

func dealTypeArray(mValue Member) string {
	if mValue.isArray {
		return "[]" + mValue.typeName
	}
	return mValue.typeName
}

func iOFuncName(typename string, isWrite bool) string {
	var wordDic = map[string]bool{
		"string":  true,
		"bool":    true,
		"byte":    true,
		"float64": true,
		"float32": true,
		"int16":   true,
		"int32":   true,
		"int64":   true,
		"uint16":  true,
		"uint32":  true,
		"uint64":  true,
	}

	if !wordDic[typename] {
		util.LogError("此类型不存在,需要检查代码:", typename)
	}

	result := strings.ToUpper(typename)[:1] + typename[1:]
	if strings.HasPrefix(typename, "u") {
		result = strings.ToUpper(typename)[:2] + typename[2:]
	}

	if isWrite {
		return "Write" + result
	}
	return "Read" + result
}

func toUpperName(name string) string {
	result := strings.ToUpper(name)[:1] + name[1:]
	return result
}

func readFile(fileName string, filePath string) []ProtoBuffer {
	names := strings.Split(fileName, "_")
	if len(names) != 2 {
		util.LogError("文件名不存在协议模板号，需要重新命名")
		return nil
	}

	inputFile, inputError := os.Open(filePath)
	if inputError != nil {
		util.LogError(inputError)
	}
	defer inputFile.Close()
	inputReader := bufio.NewReader(inputFile)
	lines := make([]string, 0)
	for {
		inputString, readerError := inputReader.ReadString('\n')
		if len(strings.TrimSpace(inputString)) > 0 {
			lines = append(lines, strings.TrimSpace(inputString))
		}
		if readerError == io.EOF {
			break
		}
	}
	desc := ""
	packageName := ""
	// todo lines
	var protobuf *ProtoBuffer = nil
	protobufList := make([]ProtoBuffer, 0)
	for _, line := range lines {
		if strings.HasPrefix(line, "package") {
			packageName = strings.Split(line, " ")[1]
			continue
		}

		if strings.HasPrefix(line, "//") {
			desc = line[2:]
			continue
		}

		if strings.HasPrefix(line, "message") || strings.HasPrefix(line, "enum") {
			protobuf = &ProtoBuffer{}
			protobuf.isEnum = strings.HasPrefix(line, "enum")
			protobuf.desc = desc
			protobuf.name = toUpperName(strings.Split(strings.TrimRight(line, "{"), " ")[1])

			pArray := strings.Split(protobuf.name, "_")
			if len(pArray) >= 2 {
				if pValue, err := strconv.Atoi(pArray[1]); err == nil {
					protobuf.number = pValue
				}
			}
		} else if protobuf != nil {
			indexTag = 0
			array := strings.FieldsFunc(line, splitLine)
			if len(array) > 0 {
				member := Member{}
				member.desc = desc
				if protobuf.isEnum {
					member.name = array[0]
					member.labelNumber = array[1]
				} else {
					index := 0
					for _, v := range array {
						if v == "repeated" || v == "optional" || v == "required" {
							member.labelName = v
							member.isArray = v == "repeated"
							continue
						}

						if v == "default" {
							continue
						}

						switch index {
						case 0:
							typeName, isMessage := getTypeDescriptor(v)
							member.typeName = typeName
							member.isMessage = isMessage
							break
						case 1:
							member.name = toUpperName(v)
							break
						case 2:
							member.labelNumber = v
							break
						default:
							if index == 3 {
								member.defaultValue = v
								break
							}
							member.desc = v
							break
						}
						index++
					}
				}
				if len(strings.TrimSpace(packageName)) > 0 {
					member.fullName = packageName + "." + protobuf.name + "." + member.name
				}
				member.fullName = protobuf.name + "." + member.name
				protobuf.members = append(protobuf.members, member)
			}
		}
		desc = ""

		if strings.HasSuffix(line, "}") && protobuf != nil {
			protobufList = append(protobufList, *protobuf)
		}
	}
	return protobufList
}

func getTypeDescriptor(typeName string) (string, bool) {
	var wordDic = map[string]string{
		"string": "string",
		"bool":   "bool",
		"byte":   "byte",
		"bytes":  "string",
		"double": "float64",
		"float":  "float32",
		"int16":  "int16",
		"int32":  "int32",
		"int64":  "int64",
		"uint16": "uint16",
		"uint32": "uint32",
		"uint64": "uint64",
	}

	if wordDic[strings.ToLower(typeName)] != "" {
		return wordDic[strings.ToLower(typeName)], false
	}
	return toUpperName(typeName), true
}

func splitLine(r rune) bool {
	if r == '/' {
		indexTag++
		return false
	}

	if indexTag > 0 {
		return false
	}
	if r == ' ' || r == ';' || r == '[' || r == ']' || r == '=' || r == '{' || r == '}' {
		return true
	}
	return false
}
