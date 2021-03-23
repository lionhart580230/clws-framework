#!/bin/bash

echo ">> 正在编译..."
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server main.go
echo ">> 编译完成..."


WORK_DIR=/tmp/docker_build_code_clws-framework
IMAGE_NAME=code_clws-framework
GROUP_NAME=lionhart
REGISTRY_URL=registry.cn-hongkong.aliyuncs.com

if test -z "$1" ; then
    echo ">> 请输入版本号"
    exit 1;
fi

if [ -d $WORK_DIR ]; then
  rm -rf ${WORK_DIR}
fi

mkdir ${WORK_DIR}
cp server ${WORK_DIR}/

echo "$1" > ${WORK_DIR}/Version

cat>${WORK_DIR}/Dockerfile<<EOF
FROM alpine:3.10.3

WORKDIR /root/

RUN apk add ca-certificates

COPY Version /root/Version
COPY server /root/server

CMD ["/root/server"]
EOF


echo ">> 构建镜像..."
docker build -t ${IMAGE_NAME} -f ${WORK_DIR}/Dockerfile ${WORK_DIR}/

echo ">> 上传镜像到docker仓库..."
docker tag $(docker images | grep ${IMAGE_NAME} | head -1 | awk '{print $3}') ${REGISTRY_URL}/${GROUP_NAME}/${IMAGE_NAME}:latest

docker push ${REGISTRY_URL}/${GROUP_NAME}/${IMAGE_NAME}:latest

docker tag $(docker images | grep ${IMAGE_NAME} | head -1 | awk '{print $3}') ${REGISTRY_URL}/${GROUP_NAME}/${IMAGE_NAME}:$1
docker push ${REGISTRY_URL}/${GROUP_NAME}/${IMAGE_NAME}:$1

#echo ">> 清理本地镜像..."
#docker rmi -f $(docker images | grep ${IMAGE_NAME} | awk '{print $3}')

#echo ">> 重启服务"
#docker pull ${REGISTRY_URL}/${GROUP_NAME}/${IMAGE_NAME}:latest
#docker-compose up -d