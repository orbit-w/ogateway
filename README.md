# OGateway

一个高性能的网关服务，支持长连接流式传输，使用KCP协议实现可靠的UDP通信。

## 项目概述

OGateway是一个基于Go语言开发的网络网关服务，主要用于处理客户端与服务器之间的通信。它支持长连接流式传输，可以高效处理大量并发连接。

## 通信协议

### 网络连接

- 使用KCP协议（基于UDP的可靠传输协议）
- 支持长连接流式传输
- 服务器默认监听端口：8900

### 消息编码规则

#### 网络层编码

##### 上行消息（客户端 -> 服务器）

网络层的消息格式如下：

```
size (int32) | gzipped (bool) | type (byte) | [length (uint32) | body (bytes) ]...
```

- `size`: 消息总长度（4字节整数）
- `gzipped`: 消息体是否经过gzip压缩（1字节布尔值）
- `type`: 消息类型（1字节整数）
- `length`: 消息内容的长度（4字节整数）
- `body`: 实际的消息数据（变长字节数组）

##### 下行消息（服务器 -> 客户端）

网络层的消息格式如下：

```
size (int32) | gzipped (bool) | type (byte) | body (bytes)
```

- `size`: 消息总长度（4字节整数）
- `gzipped`: 消息体是否经过gzip压缩（1字节布尔值）
- `type`: 消息类型（1字节整数）
- `body`: 消息体（变长字节数组）

#### 业务层编码

业务层的消息体（body）格式如下：

```
[协议号（4byte）| seq（4byte，optional）| 消息长度（4byte）| 消息内容（bytes）]...
```

- `协议号`: 标识消息类型的ID（4字节整数）
- `seq`: 消息序列号，可选字段（4字节整数）
- `消息长度`: 消息内容的长度（4字节整数）
- `消息内容`: 实际的消息数据（变长字节数组）

## 配置说明

配置文件位于 `configs/config.toml`：

## 运行项目

```bash
# 克隆项目
git clone https://github.com/orbit-w/ogateway.git
cd ogateway

# 运行服务器
go run main.go
```
