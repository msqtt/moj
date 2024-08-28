## 项目设计

### 代码结构

项目使用六边形架构组织，对于一个服务应用来说，业务逻辑不依赖任何具体的应用代码，而是依赖统一的抽象接口。应用层实现这些抽象接口从而做到**控制反转**，减少业务逻辑与应用代码的耦合性，提高系统的可扩展性。

如图所示：

<div align="center">
    <img src="./assets/hexagonal-architecture.drawio.png"/>
</div>

在业务逻辑中引入 DDD + CQRS ，代码结构示例图：

<div align="center">
    <img src="./assets/hexagonal-architecture-2.drawio.png"/>
</div>


为了进一步降低微服务间的代码耦合，实现各应用要能做到独立部署与更新，应用间的交互只能通过**消息队列**与**GRPC**，禁止任何应用层代码间的依赖。由此，提高了代码的防侵入性。

在 Golang 中实现，就需要每个 Application 创建独立的 module ，且每个 module 都有自己的依赖，相互之间没有任何依赖。

> domain models 也需要创建单独的 module，其他的 application model 只依赖于它。

目录结构：
```b
.
├── apps # 各服务应用的代码 ---+
│   ├── game              # |
│   ├── judgement         # |
│   ├── question          # |
│   ├── record            # | 依赖 / 实现
│   ├── user              # |
│   └── web-bff           # |
├── builds                # |
├── docs                  # |
├── domain # 领域模型的代码 <--+
├── go.work
├── go.work.sum
...
```
> 本项目使用了 Golang 的 Workspace 功能来组织这些 module ，这是因为把 domain 和 application 代码都放同一个仓库中便于开发，如果每份代码都分别发布在不同的仓库也是可行的。

### 领域模型

事件风暴，如下图所示：

<div align="center">
    <img src="./assets/moj-eventstorm.png"/>
</div>

领域模型的详细设计，如下图所示：

<div align="center">
    <img src="./assets/moj-ddd.png"/>
</div>

### 应用层

#### 消息队列

#### 持久层

#### 判题缓存

#### 接口协议

#### 判题机实现