package server

import (
	"fmt"
	"game/util"
	"net"
	"sync"
	"time"
)

type ClientObject struct {
	sendID     uint32 // 发送方用户ID 不是后台时用到
	receiverID uint32 // 接收放用户ID
	sendType   int    // 0=广播所有 1=广播指定用户
	JsonStr    string // 投递的 字符串
}

func (this *ClientObject) ToBytes() []byte {
	return []byte(this.JsonStr)
}

// TcpServer TCP服务器结构体
type TcpServer struct {
	Listener  net.Listener
	Clients   map[string]*Client
	MuClient  sync.Mutex
	IsRunning bool
}

func (this *TcpServer) acceptLoop() {
	for this.IsRunning {

		conn, err := this.Listener.Accept() // 建立连接
		if err != nil {
			if !this.IsRunning {
				return
			}
			util.LogError("accept failed, err:", err)
			continue
		}
		this.setGameServerKeepAlive(conn, 30*time.Second)
		cli := &Client{
			Conn:      conn,
			Server:    this,
			IpAddress: conn.RemoteAddr().String(),
		}

		this.MuClient.Lock()
		this.Clients[cli.IpAddress] = cli
		this.MuClient.Unlock()

		go cli.HandleClient()
	}
}

func (this *TcpServer) setGameServerKeepAlive(conn net.Conn, keepAlive time.Duration) {
	// 将连接转换为 TCPConn 以设置 Keep-Alive
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		util.Log("nicht TCP connet")
		conn.Close()
		return
	}
	if err := tcpConn.SetKeepAlive(true); err != nil {
		util.Log("set Keep-Alive 失败: %v", err)
	}
	if err := tcpConn.SetKeepAlivePeriod(keepAlive); err != nil {
		util.Log("set Keep-Alive 间隔失败: %v", err)
	}
}

func (this *TcpServer) allServers() []string {
	this.MuClient.Lock()
	defer this.MuClient.Unlock()

	_all := make([]string, len(this.Clients))
	idx := 0
	for _, v := range this.Clients {
		_all[idx] = v.IpAddress
		idx++
	}

	return _all
}

func (this *TcpServer) getServer(ipAddress string) *Client {
	this.MuClient.Lock()
	defer this.MuClient.Unlock()
	return this.Clients[ipAddress]
}

// SysCastAll 由另外一个客户端进行广播
func (this *TcpServer) SysCastAll(ipAddress string, obj ClientObject) {
	all := this.allServers()
	for _, v := range all {
		if v != ipAddress {
			ch := this.getServer(v)
			if ch != nil {
				ch.Forward <- obj.ToBytes()
			}
		}
	}
}

func (this *TcpServer) Start(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("listen failed, err: %v", err)
	}
	this.Listener = listener
	util.Log("服务器启动！")

	this.IsRunning = true
	go this.acceptLoop()
	return nil
}

// Stop 停止服务器
func (this *TcpServer) Stop() {
	if this.Listener != nil {
		this.IsRunning = false
		this.Listener.Close()
	}

	this.MuClient.Lock()
	for _, cli := range this.Clients {
		cli.Conn.Close()
	}
	this.Clients = make(map[string]*Client)
	this.MuClient.Unlock()
}

func NewServer() *TcpServer {
	return &TcpServer{
		Clients: make(map[string]*Client),
	}
}

func server_Test() {
	//logic.InitMessage()
	//tcpServer := NewServer()
	//err := tcpServer.Start("0.0.0.0:9301")
	//if err != nil {
	//	util.LogError("启动服务器失败: %v\n", err)
	//}
}
