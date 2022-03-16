docker-build:
	docker build -t pluralith-ci .

docker-run:
	docker run --name pluralith-ci pluralith-ci
