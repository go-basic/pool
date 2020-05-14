package main

import (
	"fmt"
	"github.com/go-basic/pool"
	"net"
	"os"
	"os/signal"
	"time"
)

const addr string = "127.0.0.1:80"

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	go server()
	//等待tcp server启动
	time.Sleep(2 * time.Second)
	client()
	fmt.Println("使用: ctrl+c 退出服务")
	for sig := range c {
		fmt.Printf("received ctrl+c(%v)\n", sig)
		os.Exit(0)
	}
	fmt.Println("服务退出")
}

type DemoCloser struct {
	Conn     net.Conn
	activeAt time.Time
}

func (p *DemoCloser) Close() error {
	return p.Conn.Close()
}

func (p *DemoCloser) GetActiveTime() time.Time {
	return p.activeAt
}

func client() {

	p, err := pool.NewGenericPool(2, 5, 30*time.Second, func() (poolable pool.Poolable, e error) {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}
		return &DemoCloser{Conn: conn, activeAt: time.Now()}, nil
	})
	if err != nil {
		fmt.Println("err=", err)
	}

	//从连接池中取得一个连接
	v, err := p.Get()

	//do something
	//conn=v.(net.Conn)

	//将连接放回连接池中
	p.Put(v)

	//释放连接池中的所有连接
	//p.Release()

	//查看当前连接中的数量
	current := p.Len()
	fmt.Println("len=", current)
}

func server() {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Error listening: ", err)
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on ", addr)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err)
		}
		fmt.Printf("Received message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())
	}
}
