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

>该项目经过生产环境实际检验，目前已在线上平稳运行两年时间
>
>截至目前，小程序用户大约37w+，日活

| 类型                   | 截图                                                         |
| ---------------------- | ------------------------------------------------------------ |
| 用户量37w，每天增量2k+ | <img src="https://cdn.fengxianhub.top/resources-master/image-20260315124945669.png" alt="image-20260315124945669" style="zoom: 33%;" /> |
| 峰值pv 6k+             | <img src="https://cdn.fengxianhub.top/resources-master/image-20260315124937972.png" alt="image-20260315124937972" style="zoom:33%;" /> |



