## 项目简介

此项目是一个基于Go+Vue开发的Kubernetes多集群管理平台，平台还集成CMDB资产管理、容器平台、CICD功能等。
<p align="center">
  <a href="https://golang.google.cn/">
    <img src="https://img.shields.io/badge/Golang-1.23-green.svg" alt="golang">
  </a>
  <a href="https://gin-gonic.com/">
    <img src="https://img.shields.io/badge/Gin-1.9.1-red.svg" alt="gin">
  </a>
  <a href="https://gorm.io/">
    <img src="https://img.shields.io/badge/Gorm-1.25.10-orange.svg" alt="gorm">
  </a>
  <a href="https://redis.io/">
    <img src="https://img.shields.io/badge/redis-5.0-brightgreen.svg" alt="redis">
  </a>
  <a href="https://vuejs.org/">
    <img src="https://img.shields.io/badge/Vue-3.0.0-orange.svg" alt="vue">
  </a>
  <a href="https://element-plus.org/zh-CN/">
    <img src="https://img.shields.io/badge/Element%20Plus-%202.8.1-blue" alt="Element Plus">
  </a>
</p>

> DYCLOUD运维平台， 本项目使用Go、 Gin、Gorm开发， 前端使用的是Vue3+Element Plus框架。


```shell
├── api
│   └── v1
├── config
├── core
├── docs
├── global
├── initialize
│   └── internal
├── middleware
├── model
│   ├── request
│   └── response
├── packfile
├── resource
│   ├── excel
│   ├── page
│   └── template
├── router
├── service
├── source
└── utils
    ├── timer
    └── upload
```

| 文件夹       | 说明                    | 描述                        |
| ------------ | ----------------------- | --------------------------- |
| `api`        | api层                   | api层 |
| `--v1`       | v1版本接口              | v1版本接口                  |
| `config`     | 配置包                  | config.yaml对应的配置结构体 |
| `core`       | 核心文件                | 核心组件(zap, viper, server)的初始化 |
| `docs`       | swagger文档目录         | swagger文档目录 |
| `global`     | 全局对象                | 全局对象 |
| `initialize` | 初始化 | router,redis,gorm,validator, timer的初始化 |
| `--internal` | 初始化内部函数 | gorm 的 longger 自定义,在此文件夹的函数只能由 `initialize` 层进行调用 |
| `middleware` | 中间件层 | 用于存放 `gin` 中间件代码 |
| `model`      | 模型层                  | 模型对应数据表              |
| `--request`  | 入参结构体              | 接收前端发送到后端的数据。  |
| `--response` | 出参结构体              | 返回给前端的数据结构体      |
| `packfile`   | 静态文件打包            | 静态文件打包 |
| `resource`   | 静态资源文件夹          | 负责存放静态文件                |
| `--excel` | excel导入导出默认路径 | excel导入导出默认路径 |
| `--page` | 表单生成器 | 表单生成器 打包后的dist |
| `--template` | 模板 | 模板文件夹,存放的是代码生成器的模板 |
| `router`     | 路由层                  | 路由层 |
| `service`    | service层               | 存放业务逻辑问题 |
| `source` | source层 | 存放初始化数据的函数 |
| `utils`      | 工具包                  | 工具函数封装            |
| `--timer` | timer | 定时器接口封装 |
| `--upload`      | oss                  | oss接口封装        |

![img.png](image/img.png)
![img_5.png](image/img_5.png)
![img_1.png](image/img_1.png)
![img_2.png](image/img_2.png)
![img_3.png](image/img_3.png)
![img_4.png](image/img_4.png)