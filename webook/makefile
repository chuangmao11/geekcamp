.PHONY: docker
docker:
	@rm webook || true
	@GOOS=linux GOARCH=arm go build -tags=k8s -o webook .
	@docker rmi -f chuangmao11/webook:v0.0.1
	@docker build -t chuangmao11/webook:v0.0.1 .