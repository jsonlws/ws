# !/bin/bash


#定义所有服务的文件注意这几个服务不能交换顺序
server_array=(
    ws #ws主程序
    grpcserv #grpc服务端
    #grpcclient  #grpc客户端程序
)