package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

func main() {
	// 设置代理服务器的监听地址和端口
	proxyAddr := "localhost:7891"

	// 创建代理服务器
	proxy := &http.Server{
		Addr: proxyAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(r.Host)
			fmt.Println("转发中")

			var resp *http.Response

			if strings.Contains(r.Host, ":443") {
				fmt.Println("字符串中包含443")
				// 创建与目标服务器的连接
				targetConn, err := net.Dial("tcp", r.Host)
				if err != nil {
					log.Printf("无法连接到目标服务器：%v", err)
					http.Error(w, "服务不可用", http.StatusServiceUnavailable)
					return
				}
				defer targetConn.Close()
				// TLS握手
				config := &tls.Config{
					InsecureSkipVerify: true,
				}
				tlsConn := tls.Client(targetConn, config)
				if err := tlsConn.Handshake(); err != nil {
					log.Printf("TLS握手失败：%v", err)
					http.Error(w, "服务不可用", http.StatusServiceUnavailable)
					return
				}

				// 创建与目标服务器的请求
				req := r.Clone(context.Background())

				// 发送请求给目标服务器
				err = req.Write(tlsConn)
				if err != nil {
					log.Printf("发送请求失败：%v", err)
					http.Error(w, "服务不可用", http.StatusServiceUnavailable)
					return
				}

				// 将目标服务器的响应返回给客户端
				resp, err = http.ReadResponse(bufio.NewReader(tlsConn), req)
				if err != nil {
					log.Printf("读取响应失败：%v", err)
					http.Error(w, "服务不可用", http.StatusServiceUnavailable)
					return
				}
				defer resp.Body.Close()

				// 将响应头写入客户端
				for key, values := range resp.Header {
					for _, value := range values {
						w.Header().Add(key, value)
					}
				}

				// 将响应体写入客户端
				io.Copy(w, resp.Body)
			} else {
				fmt.Println("URL没有443端口")
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

				for key, values := range resp.Header {
					for _, value := range values {
						w.Header().Add(key, value)
					}
				}

				// 将目标服务器响应的内容写入到响应客户端
				io.Copy(w, resp.Body)
			}

		}),
	}

	// 启动代理服务器
	log.Printf("正向代理服务器已启动，监听地址：%s\n", proxyAddr)
	log.Fatal(proxy.ListenAndServe())
}
