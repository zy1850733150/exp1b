package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 登陆值绑定
type LoginValues struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 用户结构体
type User struct {
	Username string `json:"username"`
}

// 处理连接
func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		// 读取请求行
		requestLine, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading request:", err)
			break
		}

		fmt.Print("Received: ", requestLine)

		// 解析请求行
		parts := strings.Fields(requestLine)
		if len(parts) == 0 {
			writeResponse(writer, "ERROR", "400", "Invalid request")
			continue
		}

		command := parts[0]

		switch command {
		case "UPLOAD":
			// 处理文件上传
			handleFileUpload(writer, parts[1:])
		case "DOWNLOAD":
			// 处理文件下载请求
			handleFileDownload(writer)
		case "LOGIN":
			// 处理登录请求
			handleLogin(writer, parts[1:])
		case "PING":
			// 处理心跳检测
			writeResponse(writer, "PONG", "200", "Heartbeat received")
		default:
			writeResponse(writer, "ERROR", "404", "Command not found")
		}

		// 刷新缓冲区，确保响应被发送
		writer.Flush()
	}
}

// 写响应
func writeResponse(writer *bufio.Writer, response, statusCode, message string) {
	fmt.Fprintf(writer, "%s %s %s\n", response, statusCode, message)
	writer.Flush()
}

// 处理文件上传
func handleFileUpload(writer *bufio.Writer, args []string) {
	if len(args) < 1 {
		writeResponse(writer, "ERROR", "400", "Missing file name")
		return
	}

	filename := args[0]
	filepath := filepath.Join("uploads", filename)

	// 这里应该是文件上传的逻辑
	// 假设文件内容直接写入文件
	file, err := os.Create(filepath)
	if err != nil {
		writeResponse(writer, "ERROR", "500", "Failed to create file")
		return
	}
	defer file.Close()

	// 读取文件内容
	file.WriteString("file content")

	writeResponse(writer, "OK", "200", "File uploaded")
}

// 处理文件下载请求
func handleFileDownload(writer *bufio.Writer) {
	filename := "uploads/example.txt" // 假设要下载的文件名

	// 这里应该是文件下载的逻辑
	// 读取文件内容
	file, err := os.Open(filename)
	if err != nil {
		writeResponse(writer, "ERROR", "500", "Failed to open file")
		return
	}
	defer file.Close()

	// 发送文件内容
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		writeResponse(writer, "ERROR", "500", "Failed to read file")
		return
	}

	writeResponse(writer, "OK", "200", string(fileContent))
}

// 处理登录请求
func handleLogin(writer *bufio.Writer, args []string) {
	if len(args) < 2 {
		writeResponse(writer, "ERROR", "400", "Missing username or password")
		return
	}

	username := args[0]
	password := args[1]

	// 假设的登陆名和密码
	if username == "admin" && password == "password" {
		writeResponse(writer, "OK", "200", "Login successful")
	} else {
		writeResponse(writer, "ERROR", "401", "Unauthorized")
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Listening on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}
		// 设置连接的读写超时时间
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		go handleConnection(conn)
	}
}
