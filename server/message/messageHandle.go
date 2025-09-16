package message

import (
	"game/server/packet"
	"sync"
)

type MessageClient struct {
	IpAddress string
	RoleId    uint32
}

var msgHandler map[uint16]func(hostId *MessageClient, proto packet.IProto) *packet.Packet
var lock sync.Mutex

var loginHandler func(hostMsg *MessageClient, proto packet.IProto) *packet.Packet

// 注册回调函数
func RegeditHandler(cmd uint16, handler func(*MessageClient, packet.IProto) *packet.Packet) {
	lock.Lock()
	defer lock.Unlock()

	if msgHandler == nil {
		msgHandler = make(map[uint16]func(hostId *MessageClient, proto packet.IProto) *packet.Packet)
	}
	msgHandler[cmd] = handler
}

func SetLoginHandler(handler func(hostMsg *MessageClient, proto packet.IProto) *packet.Packet) {
	lock.Lock()
	defer lock.Unlock()
	loginHandler = handler
}

func GetHandler(cmd uint16) func(*MessageClient, packet.IProto) *packet.Packet {
	lock.Lock()
	defer lock.Unlock()
	if msgHandler[cmd] == nil {
		return nil
	}
	return msgHandler[cmd]
}

func GetLoginHandler() func(hostMsg *MessageClient, proto packet.IProto) *packet.Packet {
	lock.Lock()
	defer lock.Unlock()
	return loginHandler
}
