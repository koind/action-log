run: build up

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down --remove-orphans

reup:
	docker-compose down --remove-orphans ;\
	docker-compose build ;\
	docker-compose up -d ;\

rmi:
	docker rmi $(docker images -a -q)

rm:
	docker rm $(docker ps -a -f status=exited -q)

test:
	set -e ;\
	docker-compose -f docker-compose.test.yml up --build -d ;\
	test_status_code=0 ;\
	docker-compose -f docker-compose.test.yml run integration_tests go test -v ./... || test_status_code=$$? ;\
	docker-compose -f docker-compose.test.yml down ;\
	echo "test_status_code="$$test_status_code ;\
	exit $$test_status_code ;\