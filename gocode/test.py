import socket

def send_request(host, port, command, args=""):
    # 创建TCP客户端
    client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    try:
        # 连接到服务器
        client.connect((host, port))
        
        # 构建请求
        request = f"{command} {args}\n"
        print(f"Sending to server: {request.strip()}")
        
        # 发送请求
        client.sendall(request.encode('utf-8'))
        
        # 接收响应
        response = client.recv(1024)
        print("Received from server:", response.decode('utf-8').strip())
    except Exception as e:
        print("An error occurred:", e)
    finally:
        # 关闭连接
        client.close()

# 测试 PING 命令
send_request('localhost', 8080, 'PING')

# 测试 UPLOAD 命令
send_request('localhost', 8080, 'UPLOAD', 'testfile.txt')

# 测试 DOWNLOAD 命令
send_request('localhost', 8080, 'DOWNLOAD')

# 测试 LOGIN 命令
send_request('localhost', 8080, 'LOGIN', 'admin password')