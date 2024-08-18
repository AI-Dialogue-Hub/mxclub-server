# mxclub-server

## 目录结构

```shell
.
├── README.md
├── apps
│   ├── mxclub-admin // 管理后台
│   └── mxclub-mini // 小程序
├── domain // 领域层
│   ├── common
│   ├── message
│   ├── order
│   ├── payment
│   ├── product
│   └── user
├── go.mod
├── go.sum
├── infra // 基础设施层（没写逻辑，逻辑都写在领域层了）
│   ├── order
│   └── user
├── pkg
│   ├── api 
│   ├── common // 存放各种db、缓存等中间件的封装
│   ├── constant
│   └── utils
└── script
    ├── Dockerfile.admin
    ├── Dockerfile.mini
    ├── Makefile
    └── redis_lua
```









