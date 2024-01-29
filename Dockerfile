FROM dockerproxy.com/alpine/k8s:1.29.1 AS kubectl


FROM ccr.ccs.tencentyun.com/storezhang/alpine:3.19.1 AS builder

# 复制Kubernetes控制程序
COPY --from=kubectl /usr/bin/kubectl /usr/local/bin/kubectl
# 复制主程序
COPY deploy /usr/local/bin


FROM ccr.ccs.tencentyun.com/storezhang/alpine:3.19.1

LABEL author="storezhang<华寅>" \
    email="storezhang@gmail.com" \
    qq="160290688" \
    wechat="storezhang" \
    description="Drone持续集成系统Kubernetes插件，增加以下功能：1、支持Deployment应用；2、支持Stateful应用；3、支持自动部署到K8s集群"


# 一次性复制所有程序，如果有多个COPY命令需要通过多Builder模式减少COPY登岛
COPY --from=builder /usr/local/bin/ /usr/local/bin/

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
    && chmod +x /usr/local/bin/deploy \
    \
    \
    \
    && rm -rf /var/cache/apk/*


# 执行命令
ENTRYPOINT /usr/local/bin/deploy

# 强制使用Kubernetes默认用户
ENV USERNAME default
