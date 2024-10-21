.PHONY: build up down test

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down