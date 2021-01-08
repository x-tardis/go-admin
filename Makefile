# 应用名称
name = go-admin
# 型号
model = ${name}
# 固件版本
version = `git describe --always --tags`
# api版本
APIversion = v0.0.1
# 设置固件名称
firmwareName = ${name}

execveFile := ${firmwareName}

# 路径相关
ProjectDir=.
BinDir=${CURDIR}/bin

# 编译平台
platform = CGO_ENABLED=0
# 编译选项,如tags,多个采用','分开 sqlite3
opts = -trimpath -tags=jsoniter
# 编译flags
path = github.com/thinkgos/x/builder
flags = -ldflags "-X '${path}.BuildTime=`date "+%F %T %z"`' \
	-X '${path}.GitCommit=`git rev-parse --short=8 HEAD`' \
	-X '${path}.GitFullCommit=`git rev-parse HEAD`' \
	-X '${path}.Name=${name}' \
	-X '${path}.Model=${model}' \
	-X '${path}.Version=${version}' \
	-X '${path}.APIVersion=${APIversion}' -w" # -s 引起gops无法识别go版本号,upx压缩也同样

system:
	@echo "----> system executable building..."
	@mkdir -p ${BinDir}
	@${platform} go build ${opts} ${flags} -o ${BinDir}/${execveFile} ${ProjectDir}/app
	@#upx --best --lzma ${execveFile}
	@#bzip2 -c ${execveFile} > ${execveFile}.bz2
	@echo "----> system executable build successful"

run: system
	@${BinDir}/${execveFile} server

swag:
	@echo "----> swagger docs building..."
	@swag init -d ${ProjectDir}/app --parseDependency ${ProjectDir}/apis
	@echo "----> swagger docs build successful"

clean:
	@echo "----> cleaning..."
	@go clean
	@rm -r ${BinDir}
	@echo "----> clean successful"

help:
	@echo " ------------- How to build ------------- "
	@echo " make         -- build target for system"
	@echo " run          -- build and run target for system"
	@echo " make swag    -- build swagger doc"
	@echo " make clean   -- clean build files"
	@echo " ------------- How to build ------------- "

.PHONY: system run swag clean help

