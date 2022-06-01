# ws
websocket程序


需要在项目根目录创建bin  logs这两个文件夹

项目启动在linux或mac系统下直接进入script目录
执行 chmod +x *sh 赋予执行shell脚本权限
编译成对应系统下可执行二进制文件命令
./buildserver.sh  mac   
./buildserver.sh  linux


window系统下只能进行手动编译无法使用脚步，进入service目录下对应执行
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -ldflags "-s -w" -o ./../../bin/ws  [服务文件名]


使用startall.sh脚步进行服务的启动，重启，停止以及服务状态的查看
在服务器上执行
./startall.sh start    [启动服务]
./startall.sh stop     [停止服务] 
./startall.sh restart  [重启服务] 
./startall.sh status   [获取服务的状态]
./startall.sh log   [log_file]   [实时查看日志信息]

windows系统下执行 [对应服务名].exe -f ./../config/config.json

在根目录执行 gprc生产命令 protoc -I grpc/ grpc/sendmsg/sendMsg.proto --go_out=plugins=grpc:grpc/sendmsg


//模拟心跳回复数据结构
{"action_type":"ping","old_index":[服务器心跳每次返回的值],"sender":[发送用户id]}

注意事项在在进行服务器数据之间发送时要将数据进行加密一次再往服务器发送