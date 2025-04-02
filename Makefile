TAG=v1.0.0
receiver-build-image:
	docker build -f receiver/Dockerfile -t carlson-zyc/otel-receiver-server:$(TAG) .