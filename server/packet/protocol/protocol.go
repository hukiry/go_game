package protocol

import (
	"game/server/packet"
	"sync"
)
 
const MSG_1001 uint16 = 1001 // 请求登录
const MSG_1002 uint16 = 1002 // 上传游戏数据
const MSG_1003 uint16 = 1003 // 注销或绑定
const MSG_1004 uint16 = 1004 // 心跳时间
const MSG_1005 uint16 = 1005 // 
const MSG_1006 uint16 = 1006 // 
 
const MSG_1101 uint16 = 1101 // ---邮件和公告数据操作
 
const MSG_1201 uint16 = 1201 // 所有活动：定时请求，领取操作
 
const MSG_1301 uint16 = 1301 // 排行榜
const MSG_1302 uint16 = 1302 // 反馈数据：存入到数据库，3个工作日内回复
 
const MSG_1401 uint16 = 1401 // 商店数据
 
const MSG_1501 uint16 = 1501 // 好友系统30个
const MSG_1502 uint16 = 1502 // 
const MSG_1503 uint16 = 1503 // 留言板
const MSG_1504 uint16 = 1504 // 社团队友：限制30人
 
const MSG_1601 uint16 = 1601 // 元宇宙数据:元宇宙表，发给谁，数据就存入给对方
 
const MSG_1901 uint16 = 1901 // 服务器错误处理
var msgMap = map[uint16]func() packet.IProto{
	MSG_1001: func() packet.IProto { return &Msg_1001{} },
	MSG_1002: func() packet.IProto { return &Msg_1002{} },
	MSG_1003: func() packet.IProto { return &Msg_1003{} },
	MSG_1004: func() packet.IProto { return &Msg_1004{} },
	MSG_1005: func() packet.IProto { return &Msg_1005{} },
	MSG_1006: func() packet.IProto { return &Msg_1006{} },
	MSG_1101: func() packet.IProto { return &Msg_1101{} },
	MSG_1201: func() packet.IProto { return &Msg_1201{} },
	MSG_1301: func() packet.IProto { return &Msg_1301{} },
	MSG_1302: func() packet.IProto { return &Msg_1302{} },
	MSG_1401: func() packet.IProto { return &Msg_1401{} },
	MSG_1501: func() packet.IProto { return &Msg_1501{} },
	MSG_1502: func() packet.IProto { return &Msg_1502{} },
	MSG_1503: func() packet.IProto { return &Msg_1503{} },
	MSG_1504: func() packet.IProto { return &Msg_1504{} },
	MSG_1601: func() packet.IProto { return &Msg_1601{} },
	MSG_1901: func() packet.IProto { return &Msg_1901{} },

}

var lock sync.Mutex
func GetMsgPB(cmd uint16) packet.IProto {
	lock.Lock()
	defer lock.Unlock()
	function := msgMap[cmd]
	return function()
}

