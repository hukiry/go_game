package protocol
 
import "game/server/packet"
 
//服务器错误处理
type Msg_1901 struct {
	Cmd		uint16		//请求的cmd
	Code		byte		//错误码
}

func (this *Msg_1901) GetCmd() uint16 {
	return MSG_1901
}
 
func (this *Msg_1901) WriteProto(p *packet.Packet) {
	p.WriteUInt16(this.Cmd)
	p.WriteByte(this.Code)
}
 
func (this *Msg_1901) ReadProto(p *packet.Packet) {
	this.Cmd = p.ReadUInt16()
	this.Code = p.ReadByte()
}
