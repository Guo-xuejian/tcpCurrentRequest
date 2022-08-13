package main

import (
	"log"
	"net"
)

func process(con net.Conn) {
	//循环接收客户端发送的数据
	defer func(con net.Conn) {
		err := con.Close()
		if err != nil {
			log.Println("连接关闭失败")
		}
	}(con) //关闭con

	for {
		//创建一个新的切片
		buf := make([]byte, 1024)

		//con.Read(buf)
		//1.等待客户端通过con发送信息
		//2.如果客户端没有write[发送]，协程就会阻塞于此
		log.Printf("服务器等待客户端 %s 发送信息\n", con.RemoteAddr().String())
		n, err := con.Read(buf)
		if err != nil {
			log.Println("客户端已退出,err:", err)
			return
		} else {
			//3.服务器显示客户端信息
			//log.Printf("收到了客户端（IP：%v）%d 个字节数据",con.RemoteAddr().String(),n)
			log.Printf("收到了客户端 %s 数据:%s ", con.RemoteAddr().String(), string(buf[:n]))
		}
		_, err = con.Write([]byte("hello"))
	}

}
func main() {
	log.Println("服务器开始监听...")
	//1.tcp表示使用网络协议是tcp
	//2.0.0.0.0:8888表示在本地监听8888端口
	lister, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Println("监听失败...err: ", err)
		return
	}

	defer func(lister net.Listener) {
		err := lister.Close()
		if err != nil {
			log.Println("监听关闭失败")
		}
	}(lister) //延时关闭listen

	//循环等待客户端连接
	for {
		//等待客户端连接
		log.Println("等待客户端连接")
		conn, err := lister.Accept()
		if err != nil {
			log.Println("连接Accept() 失败，err: ", err)
		} else {
			log.Printf("Accept() suc conn=%v,客户端IP=%v\n", conn, conn.RemoteAddr().String())
		}
		go process(conn)
	}
	//log.Printf("lister=%v\n",lister)
}
