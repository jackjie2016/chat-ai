package main

import (
	"io"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// 创建一个新的请求客户端
	client := &http.Client{}

	// 复制原始请求
	req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header = r.Header.Clone()

	// 发送请求到目标服务器
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	log.Printf("转发请求中")
	// 将响应的头部信息复制到响应客户端
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// 将目标服务器响应的内容写入到响应客户端
	io.Copy(w, resp.Body)
}

func main() {
	// 创建一个反向代理服务器
	proxy := &http.Server{
		Addr: ":7891",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler(w, r)
		}),
	}

	// 启动代理服务器
	log.Printf("Proxy server is running on http://localhost:7891")
	log.Fatal(proxy.ListenAndServe())
}
