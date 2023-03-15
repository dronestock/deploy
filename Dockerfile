FROM dockerproxy.com/alpine/k8s:1.26.2 AS kubectl





FROM ccr.ccs.tencentyun.com/storezhang/alpine:3.17.2


LABEL author="storezhang<华寅>" \
    email="storezhang@gmail.com" \
    qq="160290688" \
    wechat="storezhang" \
    description="Drone持续集成系统Kubernetes插件，增加以下功能：1、支持Deployment应用；2、支持Stateful应用；3、支持自动部署到K8s集群"


# 复制Kubernetes控制程序
COPY --from=kubectl /usr/bin/kubectl /usr/bin/kubectl
# 复制主程序
COPY deploy /bin


RUN set -ex \
    \
    \
    \
    && apk update \
    && apk --no-cache add curl ca-certificates \
    \
    \
    \
    # 增加执行权限
    && chmod +x /bin/deploy \
    \
    \
    \
    && rm -rf /var/cache/apk/*


# 执行命令
ENTRYPOINT /bin/deploy


# 强制使用Kubernetes默认用户
ENV USERNAME default


# 强制用户名为Kubernetes默认用户名
ENV USERNAME default
