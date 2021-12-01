# Go 框架使用说明

## 目录结构
```
project
├── README.md
├── Makefile
├── go.mod
├── go.sum
├── .gitignore
├── app
│   ├── command
│   │   └── TestCommand.go
│   ├── controller
│   │   └── TestController.go
│   ├── middleware
│   │   └── TestMiddleware.go
│   ├── service
│   │   └── TestService.go
│   └── validator
│       └── TestValidator.go
├── autoload
│   ├── command.go
│   ├── config.go
│   ├── controller.go
│   ├── loader.go
│   ├── middleware.go
│   ├── rpc.go
│   └── service.go
├── base
│   └── app
└── tsconfig.json
```