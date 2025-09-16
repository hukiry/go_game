package protocol
 
import "game/server/packet"
 
//商店数据
type Msg_1401 struct {
	State		byte		//0=接收成功，1=请求商店数据，2=支付
	RoleId		int32		//角色id
	ShopId		int32		//商店id
	Price		int16		//商店id
}

func (this *Msg_1401) GetCmd() uint16 {
	return MSG_1401
}
 
func (this *Msg_1401) WriteProto(p *packet.Packet) {
	p.WriteByte(this.State)
	p.WriteInt32(this.RoleId)
	p.WriteInt32(this.ShopId)
	p.WriteInt16(this.Price)
}
 
func (this *Msg_1401) ReadProto(p *packet.Packet) {
	this.State = p.ReadByte()
	this.RoleId = p.ReadInt32()
	this.ShopId = p.ReadInt32()
	this.Price = p.ReadInt16()
}
