#!/bin/bash

source ./servicedef.sh

#判断程序是否已经在运行
status_script(){
    for loop in ${server_array[@]}
    do 
        servicePid=`ps aux|grep ${loop}_server|grep -v grep|awk '{print $2}'`
        if [ -n "${servicePid}" ]
        then
            echo -e ${loop}'主程序在运行中\033[32m running...  \033[0m'
        else
            echo -e ${loop}'程序未启动\033[31m stop  \033[0m'
        fi
    done
}


#启动脚本，先判断脚本是否已经在运行
start_script(){
    echo '启动服务中...'
    for loop in ${server_array[@]}
    do 
        nohup ./../bin/${loop}_server -f ./../config/config.json  > ./../logs/${loop}.log 2>&1 &
        echo -e ${loop}'程序启动成功\033[32m success  \033[0m'
    done
    echo '所有服务启动完毕'
}

#停止脚本
stop_script(){
    for loop in ${server_array[@]}
    do 
        servicePid=`ps aux|grep ${loop}_server|grep -v grep|awk '{print $2}'`
        if [ -n "${servicePid}" ]
        then
        kill -TERM ${servicePid} >/dev/null 2>&1
        echo -e ${loop}'程序已停止'
        else
        echo -e ${loop}'程序未启动,无需停止'
        fi
    done
}
 
#重启脚本
reload_script(){
    stop_script
    sleep 2
    start_script
}

#帮助
help_script(){
    echo ${0}' status  指令获取运行程序状态'
    echo ${0}' start   指令启动程序'
    echo ${0}' stop    指令停止运行程序'
    echo ${0}' restart 指令重启程序'
    echo ${0}' log [log_file]  指令实时查看日志信息'
}

#日志
log_script(){
    if echo "${server_array[@]}" | grep -w ${1} &>/dev/null; 
    then
        tail -f ./../logs/${1}.log
    else
        echo -e "文件名错误或缺少需要查看日志文件名,可选项如下"
        for loop in ${server_array[@]}
        do
            echo -e  ${loop}
        done     
    fi
}


#入口函数
handle(){
    case $1 in
    start)
        start_script
        ;;
    stop)
        stop_script
        ;;
    status)
        status_script
        ;;
    restart)
        reload_script
        ;;
    log)
        log_script $2
        ;;    
    *)
        echo '请运行 '${0} 'status|start|stop|restart|log';
        ;;
    esac
}
 
if [ $# -ge 1 ]
then
    handle $1 $2
else
    help_script
fi