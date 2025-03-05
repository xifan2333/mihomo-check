# 配置文件详解

### check

```yaml
check:
  items:
    - speed
    - youtube
    - openai
    - netflix
    - disney
  concurrent: 100
  timeout: 2000
  interval: 10
```


- `items`: 检查项，可选值为 `speed` `youtube` `openai` `netflix` `disney`
- `concurrent`: 并发数量,此程序占用资源较少，并发可以设置较高
- `timeout`: 超时时间 单位毫秒 节点的最大延迟
- `interval`: 检测间隔时间 单位分钟 最低必须大于10分钟

### save

```yaml
save:
  method: http
  port: 8081
  webdav-url: "https://webdav-url/dav/"
  webdav-username: "webdav-username"
  webdav-password: "webdav-password"
  github-token: "github-token"
  github-gist-id: "github-gist-id"
  github-api-mirror: "https://worker-url/github"
  worker-url: https://worker-url
  worker-token: token 
```

- `method`: 保存方法，可选值为 `webdav` `http` `gist` `r2` `local`
- `port`: `http` 保存方式下的端口
- webdav:
    - `webdav-url`: webdav url
    - `webdav-username`: webdav 用户名
    - `webdav-password`: webdav 密码
- gist:
  - `github-token`: gist token
  - `github-gist-id`: gist id
- r2:
  - `worker-url`: worker url
  - `worker-token`: worker token
## mihomo

```yaml
mihomo:
  api-url: "http://192.168.31.11:9090"
  api-secret: "mihomo-api-secret"
```
此选项是为了检测完成后自动更新provider

- `api-url`: mihomo api url
- `api-secret`: mihomo api secret

## rename

```yaml
rename:
  flag: true
  method: "mix"
```

- `flag`: 重命名后是否增加旗帜信息
- `method`: 重命名方式 可选值为 `mix` `api` `regex`

> api 方式重命名更加准确，但耗时较长  
> regex 方式重命名更加快速，但如果`rename.yaml`文件规则不完善，可能会有部分节点无法重命名  
> mix 方式不做选择，全都要！会先进行`regex`重命名，没有匹配的再进行`api`重命名

## Proxy

```yaml
proxy:
  type: "http" # Options: http, socks
  address: "http://192.168.31.11:7890" # Proxy address
```
此处代理用于拉取订阅和保存使用，例如保存到gist时，则需要设置此选项
