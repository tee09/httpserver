/*
编写一个 HTTP 服务器，大家视个人不同情况决定完成到哪个环节，但尽量把 1 都做完：
1. 接收客户端 request，并将 request 中带的 header 写入 response header
2. 读取当前系统的环境变量中的 VERSION 配置，并写入 response header
3. Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
4. 当访问 localhost/healthz 时，应返回 200
*/
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"
)

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello World<h1>"))
	//设置version
	os.Setenv("VERSION", "0.0.4")
	version := os.Getenv("VERSION")
	w.Header().Set("OS VERSION is: %s \n", version)
	//将request header 设置到response header中
	for k, v := range r.Header {
		for _, vv := range v {
			fmt.Printf("Header Key: %s, Header value: %s \n", k, v)
			w.Header().Set(k, vv) //写入resp header
		}
		//fmt.Println(k, v)
	}

	// 记录日志并输出，取clientip
	//clintip := r.RemoteAddr
	//fmt.Println(clintip)
	//如果经过nat转换则clintIP拿到的是nat 负载均衡 proxy之后的地址
	//不是真实用户地址
	// X-REAL-IP
	// X-FORWORD-FOR 这两个header里面有真实地址
	clientip := getCurrentIP(r)
	//log.Printf("Response code: %d", 200)
	httpCode := http.StatusOK
	log.Printf("clinetip: %s  code: %v", clientip, httpCode)
}

//检查检查func
func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "working")
}

//构造取IP函数
func getCurrentIP(r *http.Request) string {
	ip := r.Header.Get("X-REAL-IP")
	if ip == "" {
		//remoteaddr IP:PORT 当请求头不存在则直接取IP,取remoteaddr：前面一段
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}
	return ip
}

//ClientIP尽最大努力实现获取客户端IP
//解析X-REAL-IP和X-FORWARDED-FOR
func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""

}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.HandleFunc("healthz", healthz)

	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
		log.Fatalf("start server failed, error: %s\n", err.Error())
	}
}
