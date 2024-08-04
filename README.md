# mxclub-server

## 目录结构

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

## 设计逻辑

### 订单完成状态状态流传逻辑设计

使用BPM对订单完成之后的流程进行管理，我们可以对订单结束后的操作分为以下几个步骤

1. 更新订单状态
2. 计算并分配报酬
3. 发送通知消息
4. 检查用户等级
5. 提醒用户评价







