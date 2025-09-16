package protocol
 
import "game/server/packet"
 
//排行榜
type Msg_1301 struct {
	Type		byte		//功能类型 1=周排行，2=月排行，3=年排行
	JsonDatas		string		//json数据：角色id，角色昵称，角色等级，角色任务积分
}

func (this *Msg_1301) GetCmd() uint16 {
	return MSG_1301
}
 
func (this *Msg_1301) WriteProto(p *packet.Packet) {
	p.WriteByte(this.Type)
	p.WriteString(this.JsonDatas)
}
 
func (this *Msg_1301) ReadProto(p *packet.Packet) {
	this.Type = p.ReadByte()
	this.JsonDatas = p.ReadString()
}
//反馈数据：存入到数据库，3个工作日内回复
type Msg_1302 struct {
	RoleId		int32		//角色id
	Type		byte		//反馈问题类型：1=游戏问题，2=崩溃问题，3=充值问题，4=其他问题
	Content		string		//反馈描述：30个字
	E_mail		string		//电子邮件
}

func (this *Msg_1302) GetCmd() uint16 {
	return MSG_1302
}
 
func (this *Msg_1302) WriteProto(p *packet.Packet) {
	p.WriteInt32(this.RoleId)
	p.WriteByte(this.Type)
	p.WriteString(this.Content)
	p.WriteString(this.E_mail)
}
 
func (this *Msg_1302) ReadProto(p *packet.Packet) {
	this.RoleId = p.ReadInt32()
	this.Type = p.ReadByte()
	this.Content = p.ReadString()
	this.E_mail = p.ReadString()
}
