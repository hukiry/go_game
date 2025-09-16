package server

import (
	"encoding/binary"
	"game/server/message"
	"game/server/packet"
	"game/server/packet/protocol"
	"game/util"
	"io"
	"log"
	"net"
)

// Client 客户端连接结构体
type Client struct {
	Conn      net.Conn
	IpAddress string
	Server    *TcpServer
	Forward   chan []byte
	RoleId    uint32
}

// 处理函数
func (this *Client) HandleClient() {
	// closing
	defer func() {
		this.Conn.Close()
		this.Server.MuClient.Lock()
		delete(this.Server.Clients, this.IpAddress)
		this.Server.MuClient.Unlock()
	}()

	header := make([]byte, 2)
	ch := make(chan []byte, 100000)
	go this.startAgent(ch)
	for {
		//读取数据长度
		n, err := io.ReadFull(this.Conn, header)
		if err != nil {
			util.LogError("error receiving header, bytes:", n, "reason:", err)
			break
		}

		// 包长度
		size := binary.BigEndian.Uint16(header) - 2
		data := make([]byte, size)
		// 读取数据包
		n, err = io.ReadFull(this.Conn, data)
		if err != nil {
			util.LogError("error receiving msg, bytes:", n, "reason:", err)
			break
		}

		//消息包通知
		ch <- data
	}
	close(ch)
}

func (this *Client) startAgent(incoming chan []byte) {
	util.Log("game server connected", this.IpAddress)
	forward := make(chan []byte, 100000)
	this.Forward = forward
	// closing
	defer func() {
		close(forward)
		util.Log("closing disconnected\n", this.IpAddress)
	}()

	messageClient := &message.MessageClient{IpAddress: this.IpAddress}
	for {
		select {
		case bodyData, ok := <-incoming: // request from game server
			if !ok {
				return
			}
			// read seqid
			reader := packet.Decode(bodyData)
			this._printMsg(reader.CMD(), "接受cmd=", reader.CMD(), ", size=", reader.Length())
			// handle request
			handle := this.handlerAgent(messageClient, reader.CMD(), reader)
			// send result
			if handle.CMD() > 0 {
				this._printMsg(reader.CMD(), "send cmd=", handle.CMD(), ", size=", handle.Length())
				n, err := this.Conn.Write(handle.ToBytes())
				if err != nil {
					log.Println("Error send reply to GS, bytes:", n, "reason:", err)
				}
				this.RoleId = messageClient.RoleId
			}

		case obj := <-forward: // forwarding packets(ie. seqid == 0)
			this._send(1099, obj)
		}
	}

}

func (this *Client) _printMsg(cmd uint16, params ...any) {
	if cmd != protocol.MSG_1004 {
		util.Log(params...)
	}
}

func (this *Client) handlerAgent(messageClient *message.MessageClient, cmd uint16, p *packet.Packet) *packet.Packet {
	// todo 根据协议号读取函数，传参数据，并处理完后返回字节码
	//读取数据
	proto := protocol.GetMsgPB(cmd)
	proto.ReadProto(p)

	//登录模块
	if cmd == protocol.MSG_1001 {
		loginHandler := message.GetLoginHandler()
		if loginHandler != nil {
			return loginHandler(messageClient, proto)
		} else {
			goto HandlerCall
		}
	} else {
		//获取回调函数
		handleCall := message.GetHandler(cmd)
		if handleCall != nil {
			return handleCall(messageClient, proto)
		} else {
			goto HandlerCall
		}
	}

HandlerCall:
	util.LogError(" 此协议未注册：", cmd)

	return packet.NewSpacePack(cmd)
}

func (this *Client) _send(cmd uint16, data []byte) {
	writer := packet.NewPack(cmd, data)
	n, err := this.Conn.Write(writer.ToBytes())
	if err != nil {
		log.Println("Error send reply to GS, bytes:", n, "reason:", err)
	}
}
