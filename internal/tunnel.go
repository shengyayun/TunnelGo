package internal

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
	"tunnel/tools"
)

type Tunnel struct {
	Name   string
	Auth   string `yaml:"auth"`
	Local  string `yaml:"local"`
	Server string `yaml:"server"`
	Remote string `yaml:"remote"`
}

var serverConnPool = make(map[string]*ssh.Client)

func (tunnel *Tunnel) Run() {
	log.Printf("%s is Running\n", tunnel.Name)
	// 本地端口监听
	listener, err := net.Listen("tcp", tunnel.Local)
	if err != nil {
		log.Fatalf("%s: Failed to listen: %v\n", tunnel.Name, err)
	}
	defer listener.Close()

	server := strings.Split(tunnel.Server, "@")
	config := &ssh.ClientConfig{
		User: server[0],
		Auth: []ssh.AuthMethod{
			tools.PublicKeyFile(tunnel.Auth),
		}, HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	for {
		localConn, err := listener.Accept()
		if err != nil {
			fmt.Printf("%s: Listener accept error: %s\n", tunnel.Name, err)
			os.Exit(1)
		}

		go tunnel.forward(localConn, server[1], config)
	}
}

func (tunnel *Tunnel) forward(localConn net.Conn, addr string, config *ssh.ClientConfig) {
	var err error
	var serverConn *ssh.Client
	var remoteConn net.Conn
	var ok bool

	try := 5
	for try > 0 {
		try--

		if serverConn, ok = serverConnPool[tunnel.Name]; !ok {
			// 如果连接池中没有，新建一个
			serverConn, err = ssh.Dial("tcp", addr, config)
			if err != nil { // 连接Server失败，再来一次
				fmt.Printf("%s: Server dial error: %s\n", tunnel.Name, err)
				time.Sleep(time.Millisecond * 200)
				continue
			} else {
				// 入池
				serverConnPool[tunnel.Name] = serverConn
			}
		}

		remoteConn, err = serverConn.Dial("tcp", tunnel.Remote)
		if err != nil { // 连接失败，清除Server连接，重来一次
			log.Printf("%s: Failed to dial remote: %v\n", tunnel.Name, err)
			serverConn.Close()
			delete(serverConnPool, tunnel.Name)
			time.Sleep(time.Millisecond * 200)
			continue
		}

		var wg sync.WaitGroup
		wg.Add(2)
		copyConn := func(writer, reader net.Conn) {
			_, err := io.Copy(writer, reader)
			if err != nil {
				log.Fatalf("%s: Failed to io copy: %v\n", tunnel.Name, err)
			}
			wg.Done()
		}
		go copyConn(localConn, remoteConn)
		go copyConn(remoteConn, localConn)

		wg.Wait()

		break
	}
}
