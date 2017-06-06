# 首先安装依赖的json 库在 GOROOT目录下
# cd /usr/lib/golang/src
# git clone https://github.com/bitly/go-simplejson
#mkdir -p src/github.com/bitly
#midir -p src/github.com/takama
#cd src/github.com/takama
#git clone https://github.com/takama/daemon.git
#cd ../../../
export GOPATH="${PWD}"
export GOTRACEBACK=crash 
go build -gcflags "-N -l" guard.go
# please chech http service with command "ss -lanp|grep 12345"
# don't use command "lsof -n|grep 12345"
