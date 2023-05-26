FROM alpine:3.15

ARG defaultStorageAccessDomain=www.xxx.com
ARG defaultStorageAccessScheme=https
ARG defaultDatabaseDsn='mysql://user:password@tcp(mysql8)/chat?collation=utf8mb4_unicode_ci&loc=Asia/Shanghai&parseTime=true'

ENV DEFAULT__DATABASE__DSN=${defaultDatabaseDsn} \
    DEFAULT__STORAGE__ACCESS__DOMAIN=${defaultStorageAccessDomain} \
    DEFAULT__STORAGE__ACCESS__SCHEME=${defaultStorageAccessScheme}

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk update && \
    apk add tzdata && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

COPY build/main /app/main

EXPOSE 80
ENTRYPOINT ["/app/main"]
