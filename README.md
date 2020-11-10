# iGBS Cloud

基于Gin框架的Reaful API应用程序设计。

## 系统技术栈

gin
mqtt
mysql
redis
minio


## 文件介绍

| 文件名  | 介绍                                     |
|---------|------------------------------------------|
| dao     | 数据抽象层， 返回的是获取的数据或者ecode |
| service | 服务层， 返回格式化好的数据或者ecode     |
| route   | 路由层， 暴露app的接口                   |

## 架构

redis保存所有运行时数据

所有基站与盒子保存的信息（运行时数据），快速响应对基站和盒子的绑定解绑关系。

## 资源
* [jenkins](http://192.168.31.51:8080/) 账号：coint 密码：123456789
* [portainer](http://192.168.31.51:9000/auth) 账号：coint 密码：123456789
* [grafana](http://192.168.31.51:3000/) 账号：admin 密码：123456789
* [emqtt](http://192.168.31.51:18083/) 账号：admin 密码：publish
* [mysql](http://192.168.31.51:3307/) 账号：root 密码：1234
* [prometheus](http://192.168.31.51:9090/) 
* [docs](docs/readme.md)


## question 

1. 多个客户端订阅了同一个主题，发布者发布主题时，每个客户端都会同时收到这个主题的消息。在客户端集群部署的场景下会出现消息重复处理的问题。

EMQ支持共享订阅，多个客户端订阅了同一个主题，发布者发布主题时，只有其中一个客户端接收到消息。
订阅前缀 `$share/<group>/[topic]`
多组客户端订阅了`$share/<group>/[topic]`、`$share/<group2>/[topic]`，发布者发布到topic，则消息会发布到每个group中，但是每个group中只有一个客户端会接收到消息。


## thanks

* [grafana看板](https://grafana.com/grafana/dashboards?orderBy=name&direction=asc)
* [RSSI定位算法](https://www.jianshu.com/p/2d071581b468)
* [ci部署](https://www.cnblogs.com/rainshi/p/12167285.html)

