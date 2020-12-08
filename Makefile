# 应用名称
name = go-admin
# 型号
model = ${name}
# 固件版本
# git describe --tags `git rev-list --tags --max-count=1`
version = v0.0.1 #`git describe --tags`
# api版本
APIversion = v0.0.1
# 设置固件名称
firmwareName = ${name}

execveFile := ${firmwareName}

# 路径相关
PROJDIR=.
BINDIR=bin
# 编译平台
platform = CGO_ENABLED=0
# 编译选项,如tags,多个采用','分开 sqlite3
opts = -trimpath -tags=jsoniter
# 编译flags
path = github.com/thinkgos/sharp/builder
flags = -ldflags "-X '${path}.BuildTime=`date "+%F %T %z"`' \
	-X '${path}.GitCommit=`git rev-parse --short=8 HEAD`' \
	-X '${path}.GitFullCommit=`git rev-parse HEAD`' \
	-X '${path}.Name=${name}' \
	-X '${path}.Model=${model}' \
	-X '${path}.Version=${version}' \
	-X '${path}.APIVersion=${APIversion}' -s -w"

system:
	@echo "----> system executable building..."
	@mkdir -p ${BINDIR}
	@${platform} go build ${opts} ${flags} -o ${BINDIR}/${execveFile} ${PROJDIR}/app
	@#upx --best --lzma ${execveFile}
	@#bzip2 -c ${execveFile} > ${execveFile}.bz2
	@echo "----> system executable build successful"

run: system
	@${BINDIR}/${execveFile} server

swag:
	@echo "----> swagger docs building..."
	@swag init -d ${PROJDIR}/app --parseDependency ${PROJDIR}/apis
	@echo "----> swagger docs build successful"

clean:
	@echo "----> cleaning..."
	@go clean
	@rm -r ${BINDIR}
	@echo "----> clean successful"

help:
	@echo " ------------- How to build ------------- "
	@echo " make         -- build target for system"
	@echo " run          -- build and run target for system"
	@echo " make swag    -- build swagger doc"
	@echo " make clean   -- clean build files"
	@echo " ------------- How to build ------------- "

.PHONY: system swag clean help

