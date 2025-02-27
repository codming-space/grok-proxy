# Grok Proxy: 高性能 OpenAI 兼容的 Grok API 代理

![GitHub release](https://img.shields.io/github/v/release/codming-space/grok-proxy) ![Docker Image Size](https://ghcr-badge.egpl.dev/codming-space/grok-proxy/size) ![Go Version](https://img.shields.io/github/go-mod/go-version/codming-space/grok-proxy) ![License](https://img.shields.io/github/license/codming-space/grok-proxy)

Grok Proxy 是一个高性能的 Go 实现的代理服务器，提供 OpenAI 兼容的 API 接口，同时将请求转发至 Grok AI 服务。本项目是 [CNFlyCat/GrokProxy](https://github.com/CNFlyCat/GrokProxy.git) Python 实现的 Go 语言重写版本，具有更高的性能和更低的资源占用。

## 特点

- OpenAI 兼容 API 接口
- Cookie 轮换管理
- User-Agent 随机轮换
- 支持流式和非流式响应模式
- API 密钥认证
- 高性能 Go 实现


## 使用 Docker 运行

### 1. 拉取 Docker 镜像

```bash
docker pull ghcr.io/codming-space/grok-proxy:latest
```

### 2. 准备配置文件

创建一个 `configs` 目录并添加 `config.yaml` 文件:
```bash
mkdir -p configs
```

创建 `configs/config.yaml` 文件，内容如下:

```yaml
cookies:
  - "sso=your-grok-cookie-1"
  - "sso=your-grok-cookie-2"

password: "your_api_password"

user_agent:
  - "Mozilla/5.0 (Macintosh; Intel Mac OS X 13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.6367.202 Safari/537.36"
  - "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"
  - "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:115.0) Gecko/20100101 Firefox/115.0"
```

### 3. 启动容器
```bash
docker run -d \
  --name grok-proxy \
  -p 8000:8000 \
  -v $(pwd)/configs:/app/configs \
  --restart unless-stopped \
  ghcr.io/codming-space/grok-proxy:latest
```

### 4. 使用 docker-compose

或者，你可以使用 docker-compose:

创建 `docker-compose.yml` 文件:

```yaml
version: '3.8'

services:
  grok-proxy:
    image: ghcr.io/codming-space/grok-proxy:latest
    ports:
      - "8000:8000"
    volumes:
      - ./configs:/app/configs
    restart: unless-stopped
```

启动服务:
```bash
docker-compose up -d
```



## 自行构建镜像

如果你想自己构建 Docker 镜像，可以按照以下步骤进行：

### 1. 克隆仓库
```bash
git clone https://github.com/codming-space/grok-proxy.git
cd grok-proxy
```

### 2. 构建镜像
```bash
docker build -t grok-proxy .
```

### 3. 运行本地构建的镜像
```bash
docker run -d \
  --name grok-proxy \
  -p 8000:8000 \
  -v $(pwd)/configs:/app/configs \
  grok-proxy
```



## 本地编译运行

### 安装 Go

确保已安装 Go 1.21 或更高版本。

### 编译
```bash
# 克隆仓库
git clone https://github.com/codming-space/grok-proxy.git
cd grok-proxy

# 下载依赖
go mod download

# 编译
go build -o grok-proxy ./cmd/server
```

### 运行
```bash
./grok-proxy
```



## API 使用

### 获取模型列表
```bash
curl http://localhost:8000/v1/models \
  -H "Authorization: Bearer your_api_password"
```

### 聊天补全 (非流式)
```bash
curl http://localhost:8000/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your_api_password" \
  -d '{
    "model": "grok-3",
    "messages": [
      {
        "role": "user",
        "content": "Hello, how are you?"
      }
    ]
  }'
```

### 聊天补全 (流式)
```bash
curl http://localhost:8000/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your_api_password" \
  -d '{
    "model": "grok-3",
    "stream": true,
    "messages": [
      {
        "role": "user",
        "content": "Hello, how are you?"
      }
    ]
  }'
```



## 关于

本项目是 [CNFlyCat/GrokProxy](https://github.com/CNFlyCat/GrokProxy.git) 的 Go 语言重写版本，旨在提供更高性能、更低资源占用的实现。原作者为 [CNFlyCat](https://github.com/CNFlyCat)。

## 贡献

欢迎提交 Pull Request 和 Issue。

## 许可证

[MIT](./LICENSE)
