package protocol
 
import "game/server/packet"
 
//好友消息
type GFriendInfo struct {
	RoleId		uint32		//用户id
	State		byte		//0=帮助一点生命，1=申请好友，2=被申请好友，3=邀请，4=被邀请挑战，5=删除好友，6，取消申请，7，取消挑战
	Level		uint16		//等级
	Nick		string		//昵称
	HeadId		byte		//头像
	MapId		byte		//好友邀请地图
}

func (this *GFriendInfo) GetCmd() uint16 {
	return 0
}
 
func (this *GFriendInfo) WriteProto(p *packet.Packet) {
	p.WriteUInt32(this.RoleId)
	p.WriteByte(this.State)
	p.WriteUInt16(this.Level)
	p.WriteString(this.Nick)
	p.WriteByte(this.HeadId)
	p.WriteByte(this.MapId)
}
 
func (this *GFriendInfo) ReadProto(p *packet.Packet) {
	this.RoleId = p.ReadUInt32()
	this.State = p.ReadByte()
	this.Level = p.ReadUInt16()
	this.Nick = p.ReadString()
	this.HeadId = p.ReadByte()
	this.MapId = p.ReadByte()
}
//社团成员信息
type GMemeberInfo struct {
	RoleId		uint32		//用户id
	Level		uint16		//等级
	Nick		string		//昵称
	HeadId		byte		//头像
	IsMasser		bool		//是团主
}

func (this *GMemeberInfo) GetCmd() uint16 {
	return 0
}
 
func (this *GMemeberInfo) WriteProto(p *packet.Packet) {
	p.WriteUInt32(this.RoleId)
	p.WriteUInt16(this.Level)
	p.WriteString(this.Nick)
	p.WriteByte(this.HeadId)
	p.WriteBool(this.IsMasser)
}
 
func (this *GMemeberInfo) ReadProto(p *packet.Packet) {
	this.RoleId = p.ReadUInt32()
	this.Level = p.ReadUInt16()
	this.Nick = p.ReadString()
	this.HeadId = p.ReadByte()
	this.IsMasser = p.ReadBool()
}
//生命信息
type GLifeInfo struct {
	Id		uint32		//聊天角色id
	Time		uint32		//寻求帮助时间
	State		byte		//状态：请求帮助，帮助其他
	Nick		string		//昵称
	Count		byte		//帮助进度
}

func (this *GLifeInfo) GetCmd() uint16 {
	return 0
}
 
func (this *GLifeInfo) WriteProto(p *packet.Packet) {
	p.WriteUInt32(this.Id)
	p.WriteUInt32(this.Time)
	p.WriteByte(this.State)
	p.WriteString(this.Nick)
	p.WriteByte(this.Count)
}
 
func (this *GLifeInfo) ReadProto(p *packet.Packet) {
	this.Id = p.ReadUInt32()
	this.Time = p.ReadUInt32()
	this.State = p.ReadByte()
	this.Nick = p.ReadString()
	this.Count = p.ReadByte()
}
//聊天消息
type GChatInfo struct {
	Id		uint32		//聊天角色id
	Time		uint32		//聊天时间
	Content		string		//聊天内容
}

func (this *GChatInfo) GetCmd() uint16 {
	return 0
}
 
func (this *GChatInfo) WriteProto(p *packet.Packet) {
	p.WriteUInt32(this.Id)
	p.WriteUInt32(this.Time)
	p.WriteString(this.Content)
}
 
func (this *GChatInfo) ReadProto(p *packet.Packet) {
	this.Id = p.ReadUInt32()
	this.Time = p.ReadUInt32()
	this.Content = p.ReadString()
}
//好友系统30个
type Msg_1501 struct {
	RoleId		uint32		//玩家id
	State		byte		//1=请求消息列表，2=我的好友，3=陌生人列表
	FriendId		uint32		//好友的id操作
	FriendInfos		[]GFriendInfo		//好友消息
}

func (this *Msg_1501) GetCmd() uint16 {
	return MSG_1501
}
 
func (this *Msg_1501) WriteProto(p *packet.Packet) {
	p.WriteUInt32(this.RoleId)
	p.WriteByte(this.State)
	p.WriteUInt32(this.FriendId)
	if this.FriendInfos == nil {
		this.FriendInfos = make([]GFriendInfo, 0)
	}
	p.WriteUInt16(uint16(len(this.FriendInfos)))
	for i := 0; i < len(this.FriendInfos); i++ {
		this.FriendInfos[i].WriteProto(p)
	}
}
 
func (this *Msg_1501) ReadProto(p *packet.Packet) {
	this.RoleId = p.ReadUInt32()
	this.State = p.ReadByte()
	this.FriendId = p.ReadUInt32()
	this.FriendInfos = make([]GFriendInfo, 0)
	FriendInfos_len := p.ReadUInt16()
	for i := 0; i < int(FriendInfos_len); i++ {
		FriendInfos_p := &GFriendInfo{}
		FriendInfos_p.ReadProto(p)
		this.FriendInfos = append(this.FriendInfos, *FriendInfos_p)
	}
}
//
type Msg_1502 struct {
	RoleId		uint32		//玩家id
	Type		byte		//操作类型
	State		byte		//操作好友状态
	FriendId		uint32		//好友的id操作
	MapId		byte		//好友邀请地图
}

func (this *Msg_1502) GetCmd() uint16 {
	return MSG_1502
}
 
func (this *Msg_1502) WriteProto(p *packet.Packet) {
	p.WriteUInt32(this.RoleId)
	p.WriteByte(this.Type)
	p.WriteByte(this.State)
	p.WriteUInt32(this.FriendId)
	p.WriteByte(this.MapId)
}
 
func (this *Msg_1502) ReadProto(p *packet.Packet) {
	this.RoleId = p.ReadUInt32()
	this.Type = p.ReadByte()
	this.State = p.ReadByte()
	this.FriendId = p.ReadUInt32()
	this.MapId = p.ReadByte()
}
//留言板
type Msg_1503 struct {
	RoleId		uint32		//玩家id
	FriendId		uint32		//好友的id操作
	State		byte		//状态：获取聊天记录，发送聊天
	Message		string		//发送聊天的消息 限制30个字符
	ChatInfos		[]GChatInfo		//聊天记录 限制20条
}

func (this *Msg_1503) GetCmd() uint16 {
	return MSG_1503
}
 
func (this *Msg_1503) WriteProto(p *packet.Packet) {
	p.WriteUInt32(this.RoleId)
	p.WriteUInt32(this.FriendId)
	p.WriteByte(this.State)
	p.WriteString(this.Message)
	if this.ChatInfos == nil {
		this.ChatInfos = make([]GChatInfo, 0)
	}
	p.WriteUInt16(uint16(len(this.ChatInfos)))
	for i := 0; i < len(this.ChatInfos); i++ {
		this.ChatInfos[i].WriteProto(p)
	}
}
 
func (this *Msg_1503) ReadProto(p *packet.Packet) {
	this.RoleId = p.ReadUInt32()
	this.FriendId = p.ReadUInt32()
	this.State = p.ReadByte()
	this.Message = p.ReadString()
	this.ChatInfos = make([]GChatInfo, 0)
	ChatInfos_len := p.ReadUInt16()
	for i := 0; i < int(ChatInfos_len); i++ {
		ChatInfos_p := &GChatInfo{}
		ChatInfos_p.ReadProto(p)
		this.ChatInfos = append(this.ChatInfos, *ChatInfos_p)
	}
}
//社团队友：限制30人
type Msg_1504 struct {
	RoleId		uint32		//玩家id
	FriendId		uint32		//帮助好友id
	GassId		uint32		//社团id
	State		byte		//社团状态
	MemeberInfos		[]GMemeberInfo		//社团队友
	LifeInfos		[]GLifeInfo		//社团生命列表
}

func (this *Msg_1504) GetCmd() uint16 {
	return MSG_1504
}
 
func (this *Msg_1504) WriteProto(p *packet.Packet) {
	p.WriteUInt32(this.RoleId)
	p.WriteUInt32(this.FriendId)
	p.WriteUInt32(this.GassId)
	p.WriteByte(this.State)
	if this.MemeberInfos == nil {
		this.MemeberInfos = make([]GMemeberInfo, 0)
	}
	p.WriteUInt16(uint16(len(this.MemeberInfos)))
	for i := 0; i < len(this.MemeberInfos); i++ {
		this.MemeberInfos[i].WriteProto(p)
	}
	if this.LifeInfos == nil {
		this.LifeInfos = make([]GLifeInfo, 0)
	}
	p.WriteUInt16(uint16(len(this.LifeInfos)))
	for i := 0; i < len(this.LifeInfos); i++ {
		this.LifeInfos[i].WriteProto(p)
	}
}
 
func (this *Msg_1504) ReadProto(p *packet.Packet) {
	this.RoleId = p.ReadUInt32()
	this.FriendId = p.ReadUInt32()
	this.GassId = p.ReadUInt32()
	this.State = p.ReadByte()
	this.MemeberInfos = make([]GMemeberInfo, 0)
	MemeberInfos_len := p.ReadUInt16()
	for i := 0; i < int(MemeberInfos_len); i++ {
		MemeberInfos_p := &GMemeberInfo{}
		MemeberInfos_p.ReadProto(p)
		this.MemeberInfos = append(this.MemeberInfos, *MemeberInfos_p)
	}
	this.LifeInfos = make([]GLifeInfo, 0)
	LifeInfos_len := p.ReadUInt16()
	for i := 0; i < int(LifeInfos_len); i++ {
		LifeInfos_p := &GLifeInfo{}
		LifeInfos_p.ReadProto(p)
		this.LifeInfos = append(this.LifeInfos, *LifeInfos_p)
	}
}
