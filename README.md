# hub-api

简单的多个大语言模型 API 汇总工具

## 特性

- 支持多种负载均衡算法：轮询、权重、最小负载
- 节点失败自动进入静默一段时间
- 支持模型映射

## 使用

下载 [https://github.com/sohaha/hub-api/releases](https://github.com/sohaha/hub-api/releases)

和其他程序一样直接启动即可，如果配置文件不存在就会自动创建。

### 部署开机自动启动

通过终端+命令可以直接安装成系统服务，就可以自动开启了

```bash
./程序 install
```


### Docker 启动
```bash
docker run -d --name hub-api -p 8181:8181 ghcr.io/sohaha/sohaha/hub-api:master
```

## 配置

配置分节点配置(provider.json)和程序配置(config.toml)

### 程序配置
config.toml

```toml
[base]
# 开启调试模式
debug = true
# 启动端口
port = 8181

[chat]
# 负载均衡算法: 0:按权重随机,1:最小连接优先,2:循环,3:按权重优先
algorithm = 0
# 访问接口使用的 key, 多个使用逗号分隔
key = 'sb123,sk-sb456'
# 内置定时重试失败节点
test_interval = 60000
# 后备模型，如果该模型失败会切换到后备模型
[chat.fallback]
'claude-3.7-sonnet' = 'o4-mini'
'o4-mini' = 'deepseek-v3'

```

默认是泛绑定，**如果希望只本地访问**

```toml
[base]
port = "127.0.0.1:8181"
```

### 节点配置
provider.json

可以直接修改配置文件或通过接口更改

1. 修改配置文件然后重启程序

2. 直接通过接口更新，无需重启

```http
POST /chat/provider
Content-Type: application/json

{
    "provider": {
        "name": "provider1",
        "base": "https://api.provider1.com/v1",
        "key": "sk-xxx",
        "models": [
            "gpt-4o-mini"
        ],
        "weight": 10
    }
}
```

#### 配置说明

简单示例:
```json
{
    "节点名称": {
        "base": "https://节点域名/v1",
        "key": "sk-xxx",
        "models": ["gpt-4o-mini"],
        "weight": 10,
        "max": 100,
        "cooldown": 6000,
        "fallback": false
    }
}
```

配置字段说明:

**节点名称**
必须唯一

**base**
接口的 base URL

**key**
接口的 token，多个使用逗号,分隔

**models**
支持的模型列表

如果需要映射模型使用 Object : { "模型名称": "真实的模型名称" }
```json
{
    "models": {"4o-mini": "gpt-4o-mini"},
}
```

**weight**
节点权重，数字越大权重越高

**max**
节点最大连接数

如果同时存在的链接超过这个值就会使用其他节点

**cooldown**
节点静默时间，单位毫秒

如果请求失败在指定时间段内不会再使用这个节点的模型

**fallback**
该节点是否后备节点

只在全部节点都失败时候才会触发

## 内置接口

### 发起会话

POST /v1/chat/completions

### 访问模型列表

GET /v1/models

### 查看模型请求次数

GET /chat/total

### 配置节点

POST /chat/provider

### 查看节点信息

GET /chat/nodes