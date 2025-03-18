docker:
	@rm webook || true
	@GOOS=linux GOARCH=arm go build -o webook-live .
	@docker rmi -f zms/webook-live:v0.0.1
	@docker build -t zms/webook-live:v0.0.1 .