.PHONY: build run stop clean

build:
	sudo docker-compose build

run:
	sudo docker-compose up -d

stop:
	sudo docker-compose down

clean:
	sudo docker-compose down -v
