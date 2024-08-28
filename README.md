## MOJ

### 项目简介

moj 是一个简单的 OnlineJudge 分布式服务应用，设计基于领域驱动 ，使用 Golang 开发。

本项目仅作为本人学习微服务过程中开发的简单样例，业务功能尚未完善，后续或许会逐步完善，仅供参考，请勿在正式生产环境下使用。

如有任何问题或建议，欢迎通过提交 issue 交流讨论。

#### 业务服务

包含基础业务：
* 用户服务
* 题目服务
* 判题服务
* 比赛服务

> 详细内容请查看[需求文档](./docs/requirement.md)。

#### 服务治理

本项目以**服务网格**作为服务治理方案，使用 Kubernetes 作为容器编排方案。

在通信协议上:
* 各服务节点间使用 GRPC 协议进行通信。
* Web-BFF 节点以 GraphQL 为规范，使用 HTTP 接口与前端通信。



### 安装说明

环境需求:

* Linux Kernel 5.8+ (support cgroup v2)
* Docker 19.03+
* Kubernetes 1.16+
* Helm 3.2.0+

安装步骤：

```bash
> helm repo add msqtt https://msqtt.github.io/helm-charts
> helm repo update
> helm install moj msqtt/moj
# install on specific namespace
# > helm install moj msqtt/moj --namespace moj
```
配置说明：
```bash
# 查看部署可选配置 values.yaml
> helm show values msqtt/moj
```

#### ServiceMesh

~~本项目的设计理念是 KISS ，所以~~推荐使用架构简单且轻量化的 Linkerd 作为 ServiceMesh ，安装方式请参考[官方文档](https://linkerd.io/)。

#### 监控设施

* 对于**指标监控**，
本项目推荐搭配使用 Prometheus ，安装方式请参考[官方文档](https://prometheus.io)。
* 对于**日志监控**，
本项目推荐搭配使用 Loki 架构，安装方式请参考[官方文档](https://grafana.com/docs/loki/)。
* 对于**看板应用和监控告警**，本项目推荐使用 Grafana ，安装方式请参考[官方文档](https://grafana.com/grafana/)。
* 对于**分布式追踪**，本项目推荐搭配使用 Jaeger ，安装方式请参考[官方文档](https://www.jaegertracing.io/)。

### 使用说明

完全启动服务后， 集群内 GraphQL 页面访问地址：

`http://moj-web-bff.<namespace>:18080/`

GraphQL 接口访问地址：

`http://moj-web-bff.<namespace>:18080/query/`

### 设计相关

详细内容请查看[设计文档](./docs/design.md)。

