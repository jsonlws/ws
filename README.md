//编译文件
//进入script

./buildserver mac  [编译成mac系统下可执行文件]
./buildserver linux [编译成linux系统下可执行文件]


运行服务
./service.sh  start

获取服务状态
./service.sh  status

停止服务
./service.sh  stop

重启服务
./service.sh  restart


//生成rpc go代码命令
//进入到grpc\sendmsg目录执行
protoc -I . sendMsg.proto --go_out=plugins=grpc:.