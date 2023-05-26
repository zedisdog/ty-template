.PHONY: checkout clean compile build-image push build-front

REGISTRY=registry.cn-chengdu.aliyuncs.com/ziyujituan/chat
TARGET=latest
VERSION=develop
IMAGE=$(REGISTRY):$(TARGET)

checkout:
	git checkout $(VERSION)

clean:
	@if [ -e build/main ];then \
		rm build/main; \
	fi

build-front:
	cd ../frontend && \
	npm run build && \
	cp -r ./build/* ../backend/internal/modules/frontend/

compile: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./build/main cmd/main/main.go

build-image: compile
	@if [ latest != $(TARGET) ];then \
		echo "build staging image"; \
		docker build -t $(IMAGE) .; \
	else \
		echo "build release image"; \
		docker build \
			--build-arg defaultDatabaseDsn='mysql://user:pass@tcp(mysql8)/chat?collation=utf8mb4_unicode_ci&loc=Asia/Shanghai&parseTime=true' \
			--build-arg defaultStorageAccessDomain=www.xxx.com \
		-t $(IMAGE) .; \
	fi

push: build-front build-image
	@docker login registry.cn-chengdu.aliyuncs.com --username=user --password=pass && \
	docker push $(IMAGE)
