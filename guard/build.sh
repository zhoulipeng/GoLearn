# 首先安装依赖的json 库在 GOROOT目录下
# cd /usr/lib/golang/src
# git clone https://github.com/bitly/go-simplejson
#mkdir -p src/github.com/bitly
#midir -p src/github.com/takama
#cd src/github.com/takama
#git clone https://github.com/takama/daemon.git
#cd ../../../
export GOPATH="${PWD}"
go build -o guard.exe
