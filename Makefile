docker-build:
	docker build -t danceladus/pluralith-ci .

docker-run:
	docker run --name pluralith-ci danceladus/pluralith-ci

docker-push:
	docker push danceladus/pluralith-ci

