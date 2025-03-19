.PHONY: docker
docker:
	@rm webook || true
	@GOOS=linux GOARCH=arm go build -tags=k8s -o webook .
	@docker rmi -f zms/webook-live:v0.0.1
	@docker build -t zms/webook-live:v0.0.1 .