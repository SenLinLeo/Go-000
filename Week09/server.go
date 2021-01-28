package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func handleConn(c net.Conn, limitCh chan struct{}) {
	limitCh <- struct{}{}
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("2006-01-02 15:04:05\n")) // 循环 读写
		if err != nil {
			c.Close()
			<-limitCh
			fmt.Println("client主动断开连接...")
			return
		}
		time.Sleep(1 * time.Second)
		fmt.Println("当前开启协程数：", len(limitCh))
	}
	<-limitCh
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal("监听出错了： ", err)
		return
	}

	limitChan := make(chan struct{}, 2)

	for {
		conn, err := listener.Accept() // 循环 listen
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConn(conn, limitChan)
	}
}
