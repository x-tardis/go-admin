
# 型号
model = go-admin
# 固件版本
# git describe --tags `git rev-list --tags --max-count=1`
version = v0.0.1 #`git describe --tags`
# api版本
APIversion = v0.0.1
# 设置固件名称
firmwareName = go-admin

execveFile := ${firmwareName}

# 路径相关
PROJDIR=..

# 编译平台
platform = CGO_ENABLED=0
# 编译选项,如tags,多个采用','分开 sqlite3
opts = -trimpath
# 编译flags
path = github.com/thinkgos/sharp/builder
flags = -ldflags "-X '${path}.BuildTime=`date "+%F %T %z"`' \
	-X '${path}.GitCommit=`git rev-parse --short=8 HEAD`' \
	-X '${path}.GitFullCommit=`git rev-parse HEAD`' \
	-X '${path}.Version=${version}' \
	-X '${path}.Model=${model}' \
	-X '${path}.APIVersion=${APIversion}' -s -w"

system:
	@echo "----> system executable building..."
	@${platform} go build ${opts} ${flags} -o ${execveFile} .
#	@upx --best --lzma ${execveFile}
	@#bzip2 -c ${execveFile} > ${execveFile}.bz2
	@echo "----> system executable build successful"

swag:
	@echo "----> swagger docs building..."
	@swag init --parseDependency ${PROJDIR}/apis
	@#swag init
	@echo "----> swagger docs build successful..."

clean:
	@echo "----> cleaning..."
	@go clean
	@rm -rf ${BINDIR}/*
	@echo "----> clean successful"

help:
	@echo " ------------- How to build ------------- "
	@echo " make         -- build target for system"
	@echo " make swag 	 -- build swagger doc"
	@echo " make clean   -- clean build files"
	@echo " ------------- How to build ------------- "

.PHONY: system swag docker clean help

