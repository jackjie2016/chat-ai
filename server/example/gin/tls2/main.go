package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	// 复制响应头
	for key, value := range resp.Header {
		for _, v := range value {
			w.Header().Add(key, v)
		}
	}

	// 设置响应状态码
	w.WriteHeader(resp.StatusCode)

	// 复制响应内容
	io.Copy(w, resp.Body)
}

func handleHTTPS(w http.ResponseWriter, req *http.Request) {
	// 连接目标主机
	conn, err := net.DialTimeout("tcp", "https://www.google.com", 10*time.Second)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // 仅适用于测试环境，请注意实际生产环境的安全性要求
			},
			Dial: func(network, addr string) (net.Conn, error) {
				deadline := time.Now().Add(10 * time.Second)
				c, err := net.DialTimeout(network, addr, 10*time.Second)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}

	// 发送CONNECT请求进行握手
	fmt.Println("req.Host", req.Host)
	connectReq, err := http.NewRequest(http.MethodConnect, "https://www.google.com", nil)
	if err != nil {
		log.Println(err)
		return
	}
	connectReq.Header.Set("Proxy-Connection", "Keep-Alive")

	connectResp, err := client.Do(connectReq)
	if err != nil {
		log.Println(err)
		return
	}
	defer connectResp.Body.Close()

	if connectResp.StatusCode != http.StatusOK {
		log.Println("CONNECT request failed:", connectResp.Status)
		return
	}
	cert, err := tls.LoadX509KeyPair("/etc/squid/ssl/proxy.7in6.com.pem ", "/etc/squid/ssl/proxy.7in6.com.key")
	if err != nil {
		// 处理错误
	}
	// 将连接转换为TLS连接
	tlsConn := tls.Client(conn, &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true, // 仅适用于测试环境，请注意实际生产环境的安全性要求
	})
	defer tlsConn.Close()

	// 连接成功后，继续转发原始的HTTP请求
	clientTransport := &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			return tlsConn, nil
		},
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 仅适用于测试环境，请注意实际生产环境的安全性要求
		},
	}

	client = &http.Client{
		Transport: clientTransport,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	// 复制响应头
	for key, value := range resp.Header {
		for _, v := range value {
			w.Header().Add(key, v)
		}
	}

	// 设置响应状态码
	w.WriteHeader(resp.StatusCode)

	// 复制响应内容
	io.Copy(w, resp.Body)
}

func main() {
	proxyAddr := ":7891" // 修改为实际的代理地址

	// 创建HTTP服务器
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodConnect {
			handleHTTPS(w, req)
		} else {
			handleHTTP(w, req)
		}
	})

	log.Println("Proxy server is running on", proxyAddr)
	log.Fatal(http.ListenAndServe(proxyAddr, nil))
}
