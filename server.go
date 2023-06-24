package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	OnLineMap map[string]*User
	mapLock   sync.RWMutex

	Message chan string
}

func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message

		this.mapLock.Lock()
		for _, cli := range this.OnLineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "{" + user.Addr + "}" + user.Name + ":" + msg
	this.Message <- sendMsg
}
func (this *Server) Handler(conn net.Conn) {
	//fmt.Println("建立成功")

	user := NewUser(conn)
	this.mapLock.Lock()
	this.OnLineMap[user.Name] = user
	this.mapLock.Unlock()
	this.BroadCast(user, "已上线")
	select {}
}
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnLineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

func (this *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		this.Handler(conn)
	}
}
