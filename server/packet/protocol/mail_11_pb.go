package protocol
 
import "game/server/packet"
 
//
type MailData struct {
	Id		int16		//邮件和公告id
	Title		string		//邮件和公告标题，角色昵称
	IsLife		bool		//邮件和公告内容
	Content		string		//邮件和公告内容
	Rewards		string		//邮件和公告奖励：后面用于充值补发json数据
}

func (this *MailData) GetCmd() uint16 {
	return 0
}
 
func (this *MailData) WriteProto(p *packet.Packet) {
	p.WriteInt16(this.Id)
	p.WriteString(this.Title)
	p.WriteBool(this.IsLife)
	p.WriteString(this.Content)
	p.WriteString(this.Rewards)
}
 
func (this *MailData) ReadProto(p *packet.Packet) {
	this.Id = p.ReadInt16()
	this.Title = p.ReadString()
	this.IsLife = p.ReadBool()
	this.Content = p.ReadString()
	this.Rewards = p.ReadString()
}
//---邮件和公告数据操作
type Msg_1101 struct {
	State		byte		//0=接收成功，1=请求数据，2=读取，3=删除，4=领取
	Type		byte		//功能类型
	RoleId		int32		//角色id
	Id		int16		//邮件和公告id 用于已读
	LanCode		string		//公告语言代码
	Mails		[]MailData		//邮件和公告id
}

func (this *Msg_1101) GetCmd() uint16 {
	return MSG_1101
}
 
func (this *Msg_1101) WriteProto(p *packet.Packet) {
	p.WriteByte(this.State)
	p.WriteByte(this.Type)
	p.WriteInt32(this.RoleId)
	p.WriteInt16(this.Id)
	p.WriteString(this.LanCode)
	if this.Mails == nil {
		this.Mails = make([]MailData, 0)
	}
	p.WriteUInt16(uint16(len(this.Mails)))
	for i := 0; i < len(this.Mails); i++ {
		this.Mails[i].WriteProto(p)
	}
}
 
func (this *Msg_1101) ReadProto(p *packet.Packet) {
	this.State = p.ReadByte()
	this.Type = p.ReadByte()
	this.RoleId = p.ReadInt32()
	this.Id = p.ReadInt16()
	this.LanCode = p.ReadString()
	this.Mails = make([]MailData, 0)
	Mails_len := p.ReadUInt16()
	for i := 0; i < int(Mails_len); i++ {
		Mails_p := &MailData{}
		Mails_p.ReadProto(p)
		this.Mails = append(this.Mails, *Mails_p)
	}
}
