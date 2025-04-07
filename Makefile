TAG=v1.0.0
receiver-build-image:
	docker build -f receiver/Dockerfile -t carlson-zyc/otel-receiver-server:$(TAG) .

run-svc1:
	go run ./svc1

run-svc2:
	go run ./svc2