# Go 框架使用说明

## 目录结构

```
project                 项目目录
├── .gitignore
├── README.md
├── go.mod
├── go.sum
├── Makefile
├── app                 应用目录
│   └── myapp           应用
│       ├── command
│       ├── controller
│       ├── middleware
│       ├── service     应用服务目录，此目录文件由代码工具自动生成后再做修改
│       └── validator
├── autoload            自动加载目录，此目录文件由代码工具自动生成
├── base                应用目录，此目录文件由代码工具自动生成
│   └── app             应用sdk
├── config              配置目录
│   ├── typs            配置定义目录
│   └── vars            配置实例化目录
├── core                框架核心目录
├── docker              部署流程控制文件目录
├── resource            资源文件目录
├── response            接口返回定义目录
└── server
    └── main.go         程序入口
```

## 环境部署

1. 安装 golang，配置 GOPATH，参考 golang 官方文档
2. 设置环境变量，PATH=PATH:$GOAPTH/bin
3. 安装框架代码生成工具
```
go install github.com/csingo/ctool@v1.0.5
```
4. 使用代码生成工具创建项目

## 服务编写规范
1. 创建的 go 文件都是用首字母大写的驼峰式命名
2. 每个业务中的多个 project 共同维护一个 proto 仓库，用于自动生成 服务、枚举值、常用结构体、sdk
3. 内部服务调用尽量使用 sdk 的方式，特殊情况特殊分析
4. 创建 project 的流程
```
ctool project::create --name=myproject                # 对应一个 Git 仓库
ctool app::create --name=myapp                        # 对应仓库中 app 目录内的文件夹
ctool app::service --app=myapp --protoPath=fullpath   # 对应 app 中的 service 文件夹中的每个文件，注意: 需要填写全路径

......
```

## 代码生成工具介绍
_注意：代码工具仅作用于当前命令行所在文件路径_

- 初始化工具
```
ctool tool::init
```

- 创建项目
```
ctool project::create --name=myproject
```

- 创建应用
```
ctool app::create --name=myapp
```

- 创建服务
```
ctool app::service --app=myapp --protoPath=fullpath
```

- 创建 controller / command / middleware
```
ctool app::controller --app=myapp --name=mycontroller
ctool app::command --app=myapp --name=mycontroller
ctool app::middleware --app=myapp --name=mycontroller
```

- 创建 config
```
ctool config::create --name=MyConf
```

- 创建 sdk
```
ctool sdk::create --app=myapp --protoPath=fullpath
```






## the end










