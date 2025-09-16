package logic

import (
	"game/db/data"
	"game/server/message"
	"game/server/packet"
	"game/server/packet/protocol"
	"game/util/timer"
	"strconv"
)

func InitMessage() {

}

func init() {
	message.SetLoginHandler(msg1001)
	message.RegeditHandler(protocol.MSG_1002, msg1002)
	message.RegeditHandler(protocol.MSG_1003, msg1003)
	message.RegeditHandler(protocol.MSG_1004, msg1004)
}

// 登录数据处理： 主机id，用户id生成
func msg1001(client *message.MessageClient, proto packet.IProto) *packet.Packet {
	//todo 保存数据
	result := proto.(*protocol.Msg_1001)
	user := data.User{ItemMoney: map[byte]uint32{}, ItemPack: map[byte]string{}}
	user.DeviceId = result.DeviceId
	user.ItemPack[data.Login_IpAddress] = client.IpAddress
	user.ItemPack[data.Login_LangCode] = result.Lang
	user.ItemPack[data.Login_Type] = strconv.Itoa(int(result.Login_type))
	user.ItemPack[data.Login_Platform] = strconv.Itoa(int(result.Platform))
	user.ReadDB()

	client.RoleId = user.RoleID

	//todo 给其赋值
	msg := &protocol.Msg_1002{}
	msg.Token = user.Token
	msg.TimeStamp = uint32(timer.GetLocalTime().Unix())

	temp := make([]protocol.ItemsResource, 0)
	for key, value := range user.ItemMoney {
		temp = append(temp, protocol.ItemsResource{Type: key, Number: int32(value)})
	}
	msg.Items = temp

	temp2 := make([]protocol.ItemsPack, 0)
	for key, value := range user.ItemPack {
		temp2 = append(temp2, protocol.ItemsPack{Type: key, Value: value})
	}
	msg.ItemsPacks = temp2
	msg.State = 1
	msg.RoleId = user.RoleID
	p := packet.NewSpacePack(msg.GetCmd())
	msg.WriteProto(p)
	return p
}

// 数据上传
func msg1002(client *message.MessageClient, proto packet.IProto) *packet.Packet {
	//todo 处理获取数据
	result := proto.(*protocol.Msg_1002)
	user := data.User{ItemMoney: map[byte]uint32{}, ItemPack: map[byte]string{}}
	user.RoleID = client.RoleId
	user.DeviceId = result.DeviceId
	for i := 0; i < len(result.Items); i++ {
		var key = result.Items[i].Type
		var number = uint32(result.Items[i].Number)
		user.ItemMoney[key] = number
	}

	for i := 0; i < len(result.ItemsPacks); i++ {
		var v = result.ItemsPacks[i]
		user.ItemPack[v.Type] = v.Value
	}
	user.SaveTime = uint32(timer.GetLocalTime().Unix())
	user.SaveDB()

	//todo 给其赋值
	msg := &protocol.Msg_1002{}
	msg.State = 0
	p := packet.NewSpacePack(proto.GetCmd())
	msg.WriteProto(p)
	return p
}

func msg1003(client *message.MessageClient, proto packet.IProto) *packet.Packet {
	result := proto.(*protocol.Msg_1003)
	user := data.User{}
	user.RoleID = client.RoleId
	if result.State == 1 {
		//登出
	} else if result.State == 2 {
		user.Token = result.Token
		user.SaveDB()
	} else if result.State == 2 {

	}
	msg := &protocol.Msg_1003{}
	p := packet.NewSpacePack(proto.GetCmd())
	msg.WriteProto(p)
	return p
}

// ping
func msg1004(client *message.MessageClient, proto packet.IProto) *packet.Packet {
	//todo 开始发送数据
	msg := &protocol.Msg_1004{}
	msg.TimeStamp = uint32(timer.GetLocalTime().Unix())
	//todo 封包
	p := packet.NewSpacePack(proto.GetCmd())
	msg.WriteProto(p)
	return p
}
