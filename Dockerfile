FROM storezhang/alpine:3.16.2


LABEL author="storezhang<华寅>" \
email="storezhang@gmail.com" \
qq="160290688" \
wechat="storezhang" \
description="Drone持续集成系统Kubernetes插件，增加以下功能：1、支持Deployment应用；2、支持Stateful应用；3、支持自动部署到K8s集群"


# 复制文件
COPY plugin /bin


RUN set -ex \
    \
    \
    \
    && apk update \
    && apk --no-cache add docker \
    \
    \
    \
    # 增加执行权限
    && chmod +x /bin/plugin \
    \
    \
    \
    && rm -rf /var/cache/apk/*


# 执行命令
ENTRYPOINT /bin/plugin
