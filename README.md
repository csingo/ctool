# Go 框架使用说明

## 目录结构

```
project             项目目录
├── .gitignore
├── README.md
├── go.mod
├── go.sum
├── Makefile
├── app             应用目录
│   ├── command
│   ├── controller
│   ├── middleware
│   ├── service     应用服务目录，此目录文件由代码工具自动生成后再做修改
│   └── validator
├── autoload        自动加载目录，此目录文件由代码工具自动生成
├── base            应用目录，此目录文件由代码工具自动生成
│   └── app         应用sdk
├── config          配置目录
│   ├── typs        配置定义目录
│   └── vars        配置实例化目录
├── core            框架核心目录
├── docker          部署流程控制文件目录
├── resource        资源文件目录
├── response        接口返回定义目录
└── server
    └── main.go     程序入口
```

## 环境部署

1. 安装 golang，配置 GOPATH，参考 golang 官方文档
2. 设置环境变量，PATH=PATH:$GOAPTH/bin
3. 安装框架代码生成工具  
   go install gitee.com/csingo/ctool@v1.0.0
4. 使用代码生成工具创建项目

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






## the end










