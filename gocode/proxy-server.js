const http = require('http');
const net = require('net');
const util = require('util');

// 用于将Stream流转换为Promise对象
const pipeline = util.promisify(require('stream').pipeline);

const proxyServer = http.createServer((req, res) => {
  // 设置CORS头部
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS');
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type');

  // 处理OPTIONS请求（预检请求）
  if (req.method === 'OPTIONS') {
    res.writeHead(204);
    res.end();
    return;
  }

  // 连接到TCP服务器
  const tcpClient = net.createConnection({ port: 8080, host: 'localhost' }, () => {
    console.log(`Connected to TCP server ${req.method} ${req.url}`);
    
    // 将HTTP请求转换为TCP请求
    // 这里假设TCP协议使用简单的文本协议，根据实际协议调整
    let tcpRequest = `${req.method} ${req.url} HTTP/1.1\r\n`;
    Object.keys(req.headers).forEach((header) => {
      if (header !== 'host' && header !== 'transfer-encoding') {
        tcpRequest += `${header}: ${req.headers[header]}\r\n`;
      }
    });
    tcpRequest += `\r\n`;

    // 发送请求头
    tcpClient.write(tcpRequest);

    // 发送请求体
    pipeline(req, tcpClient);
  });

  // 将TCP响应转换为HTTP响应
  pipeline(tcpClient, res);
});

proxyServer.listen(3000, () => {
  console.log('Proxy server is running on port 3000');
});