package protocol
 
import "game/server/packet"
 
//所有活动：定时请求，领取操作
type Msg_1201 struct {
	State		byte		////0=接收成功，1=请求数据，2=活动完成，3=活动进度值
	Type		byte		//功能类型 SystemFunctionType
	RoleId		int32		//角色id
	ConfigId		int16		//活动ID
	IsFinshed		bool		//
	CreatTime		uint32		//活动开始时间
	ExpirateTime		uint32		//活动结束时间
	StrValue		string		//活动进度字符串
	ParamsValue		string		//活动id和参数值:后端传值使用
}

func (this *Msg_1201) GetCmd() uint16 {
	return MSG_1201
}
 
func (this *Msg_1201) WriteProto(p *packet.Packet) {
	p.WriteByte(this.State)
	p.WriteByte(this.Type)
	p.WriteInt32(this.RoleId)
	p.WriteInt16(this.ConfigId)
	p.WriteBool(this.IsFinshed)
	p.WriteUInt32(this.CreatTime)
	p.WriteUInt32(this.ExpirateTime)
	p.WriteString(this.StrValue)
	p.WriteString(this.ParamsValue)
}
 
func (this *Msg_1201) ReadProto(p *packet.Packet) {
	this.State = p.ReadByte()
	this.Type = p.ReadByte()
	this.RoleId = p.ReadInt32()
	this.ConfigId = p.ReadInt16()
	this.IsFinshed = p.ReadBool()
	this.CreatTime = p.ReadUInt32()
	this.ExpirateTime = p.ReadUInt32()
	this.StrValue = p.ReadString()
	this.ParamsValue = p.ReadString()
}
