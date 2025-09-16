package protocol
 
import "game/server/packet"
 
//元宇宙信息：发起或点赞
type MetaInfo struct {
	FriendId		int32		//发起者的id
	LikeNum		int32		//地图的点赞数量
	NumberId		byte		//发起者地图id
	State		byte		//挑战状态，0=未接受挑战，1=接受过挑战
	Comment		string		//地图的所有评论
}

func (this *MetaInfo) GetCmd() uint16 {
	return 0
}
 
func (this *MetaInfo) WriteProto(p *packet.Packet) {
	p.WriteInt32(this.FriendId)
	p.WriteInt32(this.LikeNum)
	p.WriteByte(this.NumberId)
	p.WriteByte(this.State)
	p.WriteString(this.Comment)
}
 
func (this *MetaInfo) ReadProto(p *packet.Packet) {
	this.FriendId = p.ReadInt32()
	this.LikeNum = p.ReadInt32()
	this.NumberId = p.ReadByte()
	this.State = p.ReadByte()
	this.Comment = p.ReadString()
}
//元宇宙数据:元宇宙表，发给谁，数据就存入给对方
type Msg_1601 struct {
	State		byte		//0=对战消息，1=发起挑战，2=点赞，3=评论
	RoleId		int32		//我的角色id
	FriendId		int32		//好友角色id
	NumberId		byte		//好友点赞的地图id
	Comment		string		//好友发起评论
	MetaInfos		[]MetaInfo		//地图返回的挑战消息
}

func (this *Msg_1601) GetCmd() uint16 {
	return MSG_1601
}
 
func (this *Msg_1601) WriteProto(p *packet.Packet) {
	p.WriteByte(this.State)
	p.WriteInt32(this.RoleId)
	p.WriteInt32(this.FriendId)
	p.WriteByte(this.NumberId)
	p.WriteString(this.Comment)
	if this.MetaInfos == nil {
		this.MetaInfos = make([]MetaInfo, 0)
	}
	p.WriteUInt16(uint16(len(this.MetaInfos)))
	for i := 0; i < len(this.MetaInfos); i++ {
		this.MetaInfos[i].WriteProto(p)
	}
}
 
func (this *Msg_1601) ReadProto(p *packet.Packet) {
	this.State = p.ReadByte()
	this.RoleId = p.ReadInt32()
	this.FriendId = p.ReadInt32()
	this.NumberId = p.ReadByte()
	this.Comment = p.ReadString()
	this.MetaInfos = make([]MetaInfo, 0)
	MetaInfos_len := p.ReadUInt16()
	for i := 0; i < int(MetaInfos_len); i++ {
		MetaInfos_p := &MetaInfo{}
		MetaInfos_p.ReadProto(p)
		this.MetaInfos = append(this.MetaInfos, *MetaInfos_p)
	}
}
