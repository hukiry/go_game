package protocol
 
import "game/server/packet"
 
//请求登录
type Msg_1001 struct {
	DeviceId		string		//首次登录的设备id
	Lang		string		//首次登录的语言
	Login_type		byte		//登录方式
	Platform		byte		//登录平台
}

func (this *Msg_1001) GetCmd() uint16 {
	return MSG_1001
}
 
func (this *Msg_1001) WriteProto(p *packet.Packet) {
	p.WriteString(this.DeviceId)
	p.WriteString(this.Lang)
	p.WriteByte(this.Login_type)
	p.WriteByte(this.Platform)
}
 
func (this *Msg_1001) ReadProto(p *packet.Packet) {
	this.DeviceId = p.ReadString()
	this.Lang = p.ReadString()
	this.Login_type = p.ReadByte()
	this.Platform = p.ReadByte()
}
//货币资源
type ItemsResource struct {
	Type		byte		//货币类型
	Number		int32		//货币数量
}

func (this *ItemsResource) GetCmd() uint16 {
	return 0
}
 
func (this *ItemsResource) WriteProto(p *packet.Packet) {
	p.WriteByte(this.Type)
	p.WriteInt32(this.Number)
}
 
func (this *ItemsResource) ReadProto(p *packet.Packet) {
	this.Type = p.ReadByte()
	this.Number = p.ReadInt32()
}
//
type ItemsPack struct {
	Type		byte		//货币类型
	Value		string		//货币数量
}

func (this *ItemsPack) GetCmd() uint16 {
	return 0
}
 
func (this *ItemsPack) WriteProto(p *packet.Packet) {
	p.WriteByte(this.Type)
	p.WriteString(this.Value)
}
 
func (this *ItemsPack) ReadProto(p *packet.Packet) {
	this.Type = p.ReadByte()
	this.Value = p.ReadString()
}
//货币资源
type ItemsMeta struct {
	NumberId		byte		//货币类型
	State		byte		//货币数量
}

func (this *ItemsMeta) GetCmd() uint16 {
	return 0
}
 
func (this *ItemsMeta) WriteProto(p *packet.Packet) {
	p.WriteByte(this.NumberId)
	p.WriteByte(this.State)
}
 
func (this *ItemsMeta) ReadProto(p *packet.Packet) {
	this.NumberId = p.ReadByte()
	this.State = p.ReadByte()
}
//上传游戏数据
type Msg_1002 struct {
	RoleId		uint32		//
	DeviceId		string		//绑定的设备id
	Token		string		//是绑定登录
	TimeStamp		uint32		//时间搓
	State		byte		//用于切换设备
	Items		[]ItemsResource		//货币数据集合
	ItemsPacks		[]ItemsPack		//其他属性集合
}

func (this *Msg_1002) GetCmd() uint16 {
	return MSG_1002
}
 
func (this *Msg_1002) WriteProto(p *packet.Packet) {
	p.WriteUInt32(this.RoleId)
	p.WriteString(this.DeviceId)
	p.WriteString(this.Token)
	p.WriteUInt32(this.TimeStamp)
	p.WriteByte(this.State)
	if this.Items == nil {
		this.Items = make([]ItemsResource, 0)
	}
	p.WriteUInt16(uint16(len(this.Items)))
	for i := 0; i < len(this.Items); i++ {
		this.Items[i].WriteProto(p)
	}
	if this.ItemsPacks == nil {
		this.ItemsPacks = make([]ItemsPack, 0)
	}
	p.WriteUInt16(uint16(len(this.ItemsPacks)))
	for i := 0; i < len(this.ItemsPacks); i++ {
		this.ItemsPacks[i].WriteProto(p)
	}
}
 
func (this *Msg_1002) ReadProto(p *packet.Packet) {
	this.RoleId = p.ReadUInt32()
	this.DeviceId = p.ReadString()
	this.Token = p.ReadString()
	this.TimeStamp = p.ReadUInt32()
	this.State = p.ReadByte()
	this.Items = make([]ItemsResource, 0)
	Items_len := p.ReadUInt16()
	for i := 0; i < int(Items_len); i++ {
		Items_p := &ItemsResource{}
		Items_p.ReadProto(p)
		this.Items = append(this.Items, *Items_p)
	}
	this.ItemsPacks = make([]ItemsPack, 0)
	ItemsPacks_len := p.ReadUInt16()
	for i := 0; i < int(ItemsPacks_len); i++ {
		ItemsPacks_p := &ItemsPack{}
		ItemsPacks_p.ReadProto(p)
		this.ItemsPacks = append(this.ItemsPacks, *ItemsPacks_p)
	}
}
//注销或绑定
type Msg_1003 struct {
	State		byte		//0=接收成功，1=注销，2=绑定, 3=切换设备, 4=修改昵称
	Token		string		//上传绑定Token
	DeviceId		string		//绑定的设备id
	RoleNick		string		//昵称
	BindAccount		byte		//绑定的账号类型
}

func (this *Msg_1003) GetCmd() uint16 {
	return MSG_1003
}
 
func (this *Msg_1003) WriteProto(p *packet.Packet) {
	p.WriteByte(this.State)
	p.WriteString(this.Token)
	p.WriteString(this.DeviceId)
	p.WriteString(this.RoleNick)
	p.WriteByte(this.BindAccount)
}
 
func (this *Msg_1003) ReadProto(p *packet.Packet) {
	this.State = p.ReadByte()
	this.Token = p.ReadString()
	this.DeviceId = p.ReadString()
	this.RoleNick = p.ReadString()
	this.BindAccount = p.ReadByte()
}
//心跳时间
type Msg_1004 struct {
	TimeStamp		uint32		//时间戳
}

func (this *Msg_1004) GetCmd() uint16 {
	return MSG_1004
}
 
func (this *Msg_1004) WriteProto(p *packet.Packet) {
	p.WriteUInt32(this.TimeStamp)
}
 
func (this *Msg_1004) ReadProto(p *packet.Packet) {
	this.TimeStamp = p.ReadUInt32()
}
//
type Msg_1005 struct {
	Type		byte		//货币类型
	Number		byte		//货币数量
}

func (this *Msg_1005) GetCmd() uint16 {
	return MSG_1005
}
 
func (this *Msg_1005) WriteProto(p *packet.Packet) {
	p.WriteByte(this.Type)
	p.WriteByte(this.Number)
}
 
func (this *Msg_1005) ReadProto(p *packet.Packet) {
	this.Type = p.ReadByte()
	this.Number = p.ReadByte()
}
//
type Msg_1006 struct {
	HeadID		int16		//修改头像id
}

func (this *Msg_1006) GetCmd() uint16 {
	return MSG_1006
}
 
func (this *Msg_1006) WriteProto(p *packet.Packet) {
	p.WriteInt16(this.HeadID)
}
 
func (this *Msg_1006) ReadProto(p *packet.Packet) {
	this.HeadID = p.ReadInt16()
}
