kill `ps -ef |grep "micro_"|grep -v grep|awk '{print $2}'`
ps -ef |grep `ps -ef |grep "go run main.go"|grep -v grep|awk '{print $2}'`|awk '{print $2}' |xargs kill
