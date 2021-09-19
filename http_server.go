package main

import (
	context2 "context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"golang.org/x/sync/errgroup"
)

/**
 * @Author: dushuaihua
 * @Description:
 * @File:  http_server
 * @Version: 1.0.0
 * @Date: 2021/9/19 11:55
 */

type Httpserver struct {
	Ip string
	Port string
	Server *http.Server
}

type HttpHandler struct {}

func main()  {
	context, cancel := context2.WithCancel(context2.Background())
	defer cancel()
	group := new(errgroup.Group)
	// 1.启动httpserver
	// 2.有启动和退出方法
	// 3.能够处理linux signal
	// 4.errgroup运行一个线程监听linux操作

	httpS := NewHttpServer("127.0.0.1","8081")
	// 新建一个路由表
	routerlist := http.NewServeMux()

	// 添加路由函数
	routerlist.HandleFunc("/test", TestHandler)

	// 赋值给Server.Handler
	httpS.Server.Handler = routerlist


	group.Go(func() error {
		err := httpS.Server.ListenAndServe()
		return err
	})

	// 设置监听linux信号的
	group.Go(
		func() error {
			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
			<-sigs
			err := httpS.Server.Shutdown(context)
			return err
		})
	group.Wait()
}

func NewHttpServer(ip string, port string) *Httpserver {
	s := &http.Server{
		Addr : ip+":"+port,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return &Httpserver{
		Ip: ip,
		Port: port,
		Server: s,
	}
}

func TestHandler(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("hhhh"))
}