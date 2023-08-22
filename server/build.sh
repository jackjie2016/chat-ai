#!/usr/bin/env bash

# 使用方法：
# ./build.sh wms_admin_server 8


image=wechat_service
version=v1

echo "==========start=============" >>build.log
#echo "删除本地镜像" >>build.log
#docker rmi $(docker image | grep "${image}" | grep "${version}" | awk '{print $3}') >>build.log
echo "./build.sh  ${image}:${version}">>build.log
echo "开始制作镜像：$image:$version " >>build.log
 docker build -f ./Dockerfile -t ${image}:${version} . >>build.log
echo "给镜像：$image:$version 打标签 " >>build.log
docker tag ${image}:${version} registry.cn-hangzhou.aliyuncs.com/zifeng6257/${image}:${version} >>build.log
echo "推送镜像：$image:$version 到 harbor.7in6.com" >>build.log
#winpty docker push registry.cn-hangzhou.aliyuncs.com/zero-mall/${image}:${version} >>build.log
docker push registry.cn-hangzhou.aliyuncs.com/zifeng6257/${image}:${version} >>build.log

echo "完成" >>build.log
echo "==========end=============" >>build.log

