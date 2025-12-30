package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

var (
	authToken         = getEnv("token", "")
	port              = getEnv("PORT", "3000")
	activeConnections int32
	maxConnections    int32 = 100
	totalConnections  int64
	failedConnections int64
	
	upgrader = websocket.Upgrader{
		CheckOrigin:     func(r *http.Request) bool { return r.Header.Get("Origin") == "" },
		HandshakeTimeout: 10 * time.Second,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
	}
	
	bufPool = sync.Pool{
		New: func() interface{} { return make([]byte, 32*1024) },
	}
)

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

func parseAddress(addr string) (host string, portStr string, err error) {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return "", "", errors.New("地址为空")
	}

	if strings.HasPrefix(addr, "[") {
		// IPv6 格式: [ipv6]:port
		end := strings.Index(addr, "]")
		if end == -1 {
			return "", "", errors.New("无效的 IPv6 格式")
		}
		host = addr[1:end]
		if end+2 >= len(addr) {
			return "", "", errors.New("IPv6 缺少端口号")
		}
		portStr = addr[end+2:]
	} else {
		// IPv4 或域名格式
		parts := strings.Split(addr, ":")
		if len(parts) != 2 {
			return "", "", errors.New("无效的地址格式，应为 host:port")
		}
		host = parts[0]
		portStr = parts[1]
	}

	if host == "" || portStr == "" {
		return "", "", errors.New("主机或端口为空")
	}

	return host, portStr, nil
}

func handleProxy(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	atomic.AddInt64(&totalConnections, 1)

	if atomic.AddInt32(&activeConnections, 1) > maxConnections {
		atomic.AddInt32(&activeConnections, -1)
		atomic.AddInt64(&failedConnections, 1)
		log.Printf("[拒绝] %s: 达到最大并发限制 (%d)", clientIP, maxConnections)
		http.Error(w, "Too Many Connections", http.StatusTooManyRequests)
		return
	}
	defer atomic.AddInt32(&activeConnections, -1)

	protocol := r.Header.Get("Sec-WebSocket-Protocol")
	
	if authToken != "" && protocol != authToken {
		atomic.AddInt64(&failedConnections, 1)
		log.Printf("[拒绝] %s: Token 校验失败 (期望: %s, 实际: %s)", clientIP, authToken, protocol)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	header := http.Header{}
	if protocol != "" {
		header.Set("Sec-WebSocket-Protocol", protocol)
	}

	ws, err := upgrader.Upgrade(w, r, header)
	if err != nil {
		atomic.AddInt64(&failedConnections, 1)
		log.Printf("[握手失败] %s: %v", clientIP, err)
		return
	}
	defer ws.Close()

	const defaultTimeout = 5 * time.Minute
	ws.SetReadDeadline(time.Now().Add(defaultTimeout))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(defaultTimeout))
		return nil
	})

	_, msg, err := ws.ReadMessage()
	if err != nil {
		log.Printf("[读取失败] %s: %v", clientIP, err)
		return
	}

	if !strings.HasPrefix(string(msg), "CONNECT:") {
		log.Printf("[协议错误] %s: 未知的消息格式", clientIP)
		ws.WriteMessage(websocket.TextMessage, []byte("ERROR:Invalid message format"))
		return
	}

	payload := strings.TrimPrefix(string(msg), "CONNECT:")
	parts := strings.SplitN(payload, "|", 3)
	if len(parts) < 1 || parts[0] == "" {
		log.Printf("[协议错误] %s: 无效的 CONNECT 消息", clientIP)
		ws.WriteMessage(websocket.TextMessage, []byte("ERROR:Invalid CONNECT message"))
		return
	}

	targetAddr := parts[0]
	firstFrame := ""
	if len(parts) > 1 {
		firstFrame = parts[1]
	}

	remote, err := net.DialTimeout("tcp", targetAddr, 15*time.Second)
	if err != nil {
		atomic.AddInt64(&failedConnections, 1)
		log.Printf("[拨号失败] %s -> %s: %v", clientIP, targetAddr, err)
		ws.WriteMessage(websocket.TextMessage, []byte("ERROR:"+err.Error()))
		return
	}
	defer remote.Close()

	if firstFrame != "" {
		if _, err := remote.Write([]byte(firstFrame)); err != nil {
			atomic.AddInt64(&failedConnections, 1)
			log.Printf("[首帧失败] %s: %v", clientIP, err)
			ws.WriteMessage(websocket.TextMessage, []byte("ERROR:Failed to write first frame"))
			return
		}
	}

	if err := ws.WriteMessage(websocket.TextMessage, []byte("CONNECTED")); err != nil {
		atomic.AddInt64(&failedConnections, 1)
		log.Printf("[响应失败] %s: 发送 CONNECTED 失败: %v", clientIP, err)
		return
	}

	log.Printf("[连接建立] %s -> %s (活跃: %d)", clientIP, targetAddr, atomic.LoadInt32(&activeConnections))

	done := make(chan struct{}, 2)

	// WebSocket -> TCP
	go func() {
		defer func() { done <- struct{}{} }()
		for {
			mt, data, err := ws.ReadMessage()
			if err != nil {
				break
			}
			if mt == websocket.BinaryMessage {
				if _, err := remote.Write(data); err != nil {
					break
				}
			} else if mt == websocket.TextMessage {
				if string(data) == "CLOSE" {
					break
				}
			}
		}
	}()

	// TCP -> WebSocket
	go func() {
		defer func() { done <- struct{}{} }()
		buf := bufPool.Get().([]byte)
		defer bufPool.Put(buf)
		for {
			remote.SetReadDeadline(time.Now().Add(defaultTimeout))
			n, err := remote.Read(buf)
			if n > 0 {
				if err := ws.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
					break
				}
			}
			if err != nil {
				break
			}
		}
	}()

	<-done
	log.Printf("[连接关闭] %s -> %s", clientIP, targetAddr)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			active := atomic.LoadInt32(&activeConnections)
			total := atomic.LoadInt64(&totalConnections)
			failed := atomic.LoadInt64(&failedConnections)
			log.Printf("[统计] 活跃连接: %d, 总连接: %d, 失败: %d", active, total, failed)
		}
	}()

	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/", handleProxy)

	server := &http.Server{
		Addr:         "0.0.0.0:" + port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("ClawCloud Go ech服务启动 | 端口: %s | 最大并发: %d", port, maxConnections)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在优雅关闭服务器...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("关闭超时，强制退出: %v", err)
	}
	log.Println("服务器已关闭")
}
