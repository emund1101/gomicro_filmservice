cd bin
nohup ./micro_film > /dev/null 2>&1 &
nohup ./micro_order >/dev/null 2>&1 &
nohup ./micro_user >/dev/null 2>&1 &

cd ../gateway
nohup go run main.go  > /dev/null 2>&1 &
