# mxclub-server

#### 目录结构

- apps
  - mxclub-mini：小程序
  - mxclub-admin：管理后台

- domain - 领域模型
  - aggregate - 聚合
  - entity - 实体
  - event - 领域事件
  - vo - 值对象
  - po - 持久化对象
  - *.go - 领域服务
- adapter - 端口适配器
  - controller - 控制器
  - repository - 仓库
- server - 服务端程序入口
  - conf - 配置文件
  - main.go - 主函数
- infra - 基础设施组件

