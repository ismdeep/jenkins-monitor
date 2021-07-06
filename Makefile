IMAGE_VERSION=0.2.0

install:
	go mod tidy
	pip3.9 install -r requirements.txt

doraemon-web:
	go run github.com/ismdeep/jenkins-monitor -c config.json

docker-build:
	docker build -t ismdeep/jenkins-monitor:$(IMAGE_VERSION) .
	docker build -t ismdeep/jenkins-monitor:latest .

docker-push:
	docker push ismdeep/jenkins-monitor:$(IMAGE_VERSION)
	docker push ismdeep/jenkins-monitor:latest

docker-run-doraemon-web:
	docker run --rm -v $(CURDIR)/config.json:/config.json -v $(CURDIR)/k8s.json:/k8s.json ismdeep/jenkins-monitor:0.1.1
