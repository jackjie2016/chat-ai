#!/usr/bin/env bash

# 使用方法：
# ./build.sh wms_admin_server 8
# docker login --username=32521*****@qq.com registry.cn-hangzhou.aliyuncs.com

image=wechat_service_web
version=v1

echo "删除本地镜像"
#docker rmi $(docker image | grep "${image}" | grep "${version}" | awk '{print $3}')
echo "./build.sh  ${image}:${version}">>build.log
echo "开始制作镜像：$image:$version "
docker build -f ./Dockerfile -t ${image}:${version} .
#docker build -f ./Dockerfile -t wms_front:v24 .
echo "给镜像：$image:$version 打标签 "
docker tag ${image}:${version} registry.cn-hangzhou.aliyuncs.com/zifeng6257/${image}:${version}
#docker tag wms_front:v26 registry.cn-hangzhou.aliyuncs.com/zifeng6257/wms_front:v26
echo "推送镜像：$image:$version 到 harbor.7in6.com"
docker push registry.cn-hangzhou.aliyuncs.com/zifeng6257/${image}:${version}
#docker push registry.cn-hangzhou.aliyuncs.com/zifeng6257/wms_front:v26
echo "完成"
