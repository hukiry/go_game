package packet

import (
	"bytes"
	"encoding/binary"
	"game/util"
	"reflect"
	"sync"
)

const PACKET_LIMIT = 65535

var mut sync.Mutex

const MSG_LEN = 2

type IProto interface {
	WriteProto(p *Packet) //写入
	ReadProto(p *Packet)  //读取
	GetCmd() uint16       //获取协议号
}

type Packet struct {
	pos     int    //偏移位置
	data    []byte //数据包
	cmd     uint16 //cmd
	size    uint16 //数据大小
	version byte   //版本号
}

func NewPack(cmd uint16, data []byte) *Packet {
	return &Packet{data: data, cmd: cmd, size: uint16(len(data))}
}
func NewSpacePack(cmd uint16) *Packet {
	return &Packet{data: []byte{}, cmd: cmd}
}
func (p *Packet) MsgBody() []byte {
	return p.data
}

func (p *Packet) CMD() uint16 {
	return p.cmd
}

// 解析包，读取
func Decode(data []byte) *Packet {
	reader := &Packet{data: data}
	reader.cmd = binary.BigEndian.Uint16(reader.data[0:MSG_LEN])
	resultData := reader.data[MSG_LEN:]
	reader.pos = 0
	reader.size = uint16(len(resultData))
	reader.data = resultData
	return reader
}

// 打整包数据，发送
func (p *Packet) ToBytes() []byte {
	size := uint16(len(p.data))
	writer := &Packet{data: []byte{}}
	//前后端一致，必须写入数据包的总长度：2+cmd+数据包，便于前端计算
	writer.WriteRawBytes(binary.BigEndian.AppendUint16(make([]byte, 0), size+MSG_LEN+2))
	writer.WriteRawBytes(binary.BigEndian.AppendUint16(make([]byte, 0), p.cmd))
	writer.WriteRawBytes(p.data)
	return writer.MsgBody()
}
func (p *Packet) Length() int {
	return len(p.data)
}

func (p *Packet) WriteBool(val bool) {
	if val {
		p.data = append(p.data, byte(1))
	} else {
		p.data = append(p.data, byte(0))
	}
	p.pos += 1
}
func (p *Packet) WriteString(val string) {
	temp := []byte(val)
	pos := len(temp)
	p.WriteUInt16(uint16(pos))       //写长度
	p.data = append(p.data, temp...) //写字符串
	p.pos += pos
}
func (p *Packet) WriteByte(val byte) {
	pos := 1
	p.data = append(p.data, val)
	p.pos += pos
}
func (p *Packet) WriteInt16(val int16) {
	pos := 2
	temp := write(val)
	p.data = append(p.data, temp...)
	p.pos += pos
}
func (p *Packet) WriteInt32(val int32) {
	pos := 4
	temp := write(val)
	p.data = append(p.data, temp...)
	p.pos += pos
}
func (p *Packet) WriteInt64(val int64) {
	pos := 8
	temp := write(val)
	p.data = append(p.data, temp...)
	p.pos += pos
}
func (p *Packet) WriteUInt16(val uint16) {
	pos := 2
	temp := write(val)
	p.data = append(p.data, temp...)
	p.pos += pos
}
func (p *Packet) WriteUInt32(val uint32) {
	pos := 4
	temp := write(val)
	p.data = append(p.data, temp...)
	p.pos += pos
}
func (p *Packet) WriteUInt64(val uint64) {
	pos := 8
	temp := write(val)
	p.data = append(p.data, temp...)
	p.pos += pos
}
func (p *Packet) WriteFloat32(val float32) {
	pos := 4
	temp := write(val)
	p.data = append(p.data, temp...)
	p.pos += pos
}
func (p *Packet) WriteFloat64(val float64) {
	pos := 8
	temp := write(val)
	p.data = append(p.data, temp...)
	p.pos += pos
}
func (p *Packet) WriteRawBytes(data []byte) {
	p.data = append(p.data, data...)
	p.pos += len(data)
}

func write[T any](val T) []byte {
	mut.Lock()
	defer mut.Unlock()
	buf := bytes.NewBuffer([]byte{})
	if err := binary.Write(buf, binary.LittleEndian, val); err != nil {
		util.LogError(err)
	}
	return buf.Bytes()
}
func read[T any](b []byte) T {
	mut.Lock()
	defer mut.Unlock()
	bytesBuffer := bytes.NewBuffer(b)
	var x T
	if err := binary.Read(bytesBuffer, binary.LittleEndian, &x); err != nil {
		util.LogError(err)
	}
	return x
}

func (p *Packet) ReadBool() (ret bool) {
	b := p.ReadByte()
	if b != byte(1) {
		return false
	}

	return true
}
func (p *Packet) ReadString() (ret string) {
	if p.pos+2 > len(p.data) {
		util.LogError("read string header failed")
		return
	}

	size := p.ReadUInt16()
	if p.pos+int(size) > len(p.data) {
		util.LogError("read string data failed")
		return
	}

	bytes := p.data[p.pos : p.pos+int(size)]
	p.pos += int(size)
	ret = string(bytes)
	return
}
func (p *Packet) ReadByte() (ret byte) {
	if p.pos >= len(p.data) {
		util.LogError("read byte failed")
		return
	}

	ret = p.data[p.pos]
	p.pos++
	return
}
func (p *Packet) ReadInt16() (ret int16) {
	if p.pos+2 > len(p.data) {
		util.LogError("read uint16 failed")
		return
	}

	buf := p.data[p.pos : p.pos+2]
	ret = read[int16](buf)
	p.pos += 2
	return
}
func (p *Packet) ReadInt32() (ret int32) {
	if p.pos+4 > len(p.data) {
		util.LogError("read int32 failed")
		return
	}
	buf := p.data[p.pos : p.pos+4]
	ret = read[int32](buf)
	p.pos += 4
	return
}
func (p *Packet) ReadInt64() (ret int64) {
	if p.pos+8 > len(p.data) {
		util.LogError("read int64 failed")
		return
	}

	buf := p.data[p.pos : p.pos+8]
	ret = read[int64](buf)
	p.pos += 8
	return
}
func (p *Packet) ReadUInt16() (ret uint16) {
	if p.pos+2 > len(p.data) {
		util.LogError("read uint16 failed")
		ret = 0
		return
	}

	buf := p.data[p.pos : p.pos+2]
	ret = read[uint16](buf)
	p.pos += 2
	return
}
func (p *Packet) ReadUInt32() (ret uint32) {
	if p.pos+4 > len(p.data) {
		util.LogError("read uint32 failed")
		return
	}

	buf := p.data[p.pos : p.pos+4]
	ret = read[uint32](buf)
	p.pos += 4
	return
}
func (p *Packet) ReadUInt64() (ret uint64) {
	if p.pos+8 > len(p.data) {
		util.LogError("read uint64 failed")
		return
	}

	buf := p.data[p.pos : p.pos+8]
	ret = read[uint64](buf)
	p.pos += 8
	return
}
func (p *Packet) ReadFloat32() (ret float32) {
	if p.pos+4 > len(p.data) {
		util.LogError("read float32 failed")
		return
	}
	buf := p.data[p.pos : p.pos+4]
	ret = read[float32](buf)
	p.pos += 4
	return
}
func (p *Packet) ReadFloat64() (ret float64) {
	if p.pos+8 > len(p.data) {
		util.LogError("read float64 failed")
		return
	}

	buf := p.data[p.pos : p.pos+8]
	ret = read[float64](buf)
	p.pos += 8
	return
}

func (p *Packet) convertBytes(val any) {
	v := reflect.ValueOf(val)
	kind := v.Kind()
	switch kind {
	case reflect.Bool:
		p.WriteBool(v.Bool())
	case reflect.Int8:
	case reflect.Uint8:
		p.WriteByte(byte(v.Uint()))
	case reflect.Uint16:
		p.WriteUInt16(uint16(v.Uint()))
	case reflect.Uint32:
		p.WriteUInt32(uint32(v.Uint()))
	case reflect.Uint64:
		p.WriteUInt64(uint64(v.Uint()))
	case reflect.Int16:
		p.WriteInt16(int16(v.Int()))
	case reflect.Int32:
		p.WriteInt32(int32(v.Int()))
	case reflect.Int64:
		p.WriteInt64(int64(v.Int()))
	case reflect.String:
		p.WriteString(v.String())
	default:
		util.LogError("cannot pack type:", v)
	}
}
func (p *Packet) Write(v reflect.Value) {
	kind := v.Kind()
	switch kind {
	case reflect.Bool:
		p.WriteBool(v.Bool())
	case reflect.Uint8:
	case reflect.Int8:
		p.WriteByte(byte(v.Uint()))
	case reflect.Uint16:
		p.WriteUInt16(uint16(v.Uint()))
	case reflect.Uint32:
		p.WriteUInt32(uint32(v.Uint()))
	case reflect.Uint64:
		p.WriteUInt64(uint64(v.Uint()))
	case reflect.Int16:
		p.WriteInt16(int16(v.Int()))
	case reflect.Int32:
	case reflect.Int:
		p.WriteInt32(int32(v.Int()))
	case reflect.Int64:
		p.WriteInt64(int64(v.Int()))
	case reflect.String:
		p.WriteString(v.String())
	case reflect.Float32:
		p.WriteFloat32(float32(v.Float()))
	case reflect.Float64:
		p.WriteFloat64(v.Float())
	case reflect.Struct:
		numFields := v.NumField()
		for i := 0; i < numFields; i++ {
			p.Write(v.Field(i))
		}
	default:
		util.LogError("cannot pack type:", kind)
	}
}
