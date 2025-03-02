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
  concurrent: 10
  timeout: 10
  interval: 10
```


- `items`: 检查项，可选值为 `speed` `youtube` `openai` `netflix` `disney`
- `concurrent`: 并发数量
- `timeout`: 超时时间 单位毫秒
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

  - `method`: 保存方法，可选值为 `webdav` `http` `gist` `r2`
  - `port`: 保存端口
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

- `api-url`: mihomo api url
- `api-secret`: mihomo api secret

## rename

```yaml
rename:
  flag: true
  method: "mix"
```

- `flag`: 是否启用重命名
- `method`: 重命名方式 可选值为 `mix` `api` `regex`

> 重命名方式为 `mix` 时，会先进行`regex`重命名，然后进行`api`重命名

