
# !/bin/bash

#构建所有的服务脚本

source ./servicedef.sh

#编译成linux下可执行文件
linux_script(){
    for loop in ${server_array[@]}
    do 
    CGO_ENABLED=0  GOOS=linux  GOARCH=amd64 go build -ldflags "-s -w" -o ./../bin/${loop}_server  ./../service/${loop}/${loop}.go
    done
}

#编译成mac下可执行文件
mac_script(){
    for loop in ${server_array[@]}
    do 
    go build -ldflags "-s -w" -o ./../bin/${loop}_server  ./../service/${loop}/${loop}.go
    done
}

help_script(){
    echo ${0}' mac    编译成mac系统下可运行二进制文件'
    echo ${0}' linux  编译成linux系统下可运行二进制文件'
}


#入口函数
handle(){
    case $1 in
    mac)
        mac_script
        ;;
    linux)
        linux_script
        ;;
    help)
        help_script
        ;;
    *)
        echo '请运行 '${0} 'mac|linux|help';
        ;;
    esac
}
 
if [ $# -eq 1 ]
then
    handle $1
else
    help_script
fi