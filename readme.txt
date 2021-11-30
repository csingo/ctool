目录结构

project
    - app                           [app服务目录]
        - command
        - controller
        - middleware
        - orm
        - validator
        - service
    - autoload                      [自动加载目录]
        command.go
        config.go
        controller.go
        middleware.go
        grpc.go
    - base                          [pb文件目录, 一个服务一个文件夹]
    - bin                           [编译后的二进制执行文件目录]
    - config                        [配置目录]
        - var                       [配置对象实例目录]
            CommandConf.go
            ConfigCenterConf.go
            DatabaseConf.go
            RedisConf.go
            RouteConf.go
            ServerConf.go
        - typ                       [配置对象定义目录]
            CommandConf.go
            ConfigCenterConf.go
            DatabaseConf.go
            RedisConf.go
            RouteConf.go
            ServerConf.go
    - core                          [框架核心目录]
        - ccommand
        - cconfig
        - chelper
        - cmiddleware
        - cserver
        - cdb
        - credis
    - docker                        [部署配置目录]
        Dockerfile
        deployment.yml
        docker-compose.yml
    - resource                      [资源文件目录]
    - response                      [定义响应类型]
        typ.go
        function.go
    main.go                         [程序入口]
    Makefile                        [程序构建文件]

代码生成工具

tool project::create --name=project [--path=.]
    创建项目
    复制目录:   autoload | base | bin | config | core | docker | resource | response | services | server | Makefile

    --name  项目名称
    --path  指定生成项目的目录，不存在则新建。默认当前目录

tool config::create --name=TestConf --listen=false

tool app::create --name=app [--path=.]
    创建应用
    在项目中新增一个应用，其中包含目录: command | controllers | middlewares | orms | validators

    --name          应用名称
    --appPath       应用所有目录，不存在则新建。默认当前目录

tool app::service --app=app [--protoPath=.]
tool app::sdk --sdk=app2 [--protoPath=.]
tool app::command --app=app --name=TestCommand
tool app::controller --app=app --name=TestController
tool app::orm --app=app --name=TestOrm

tool proto::import --protoPath=.