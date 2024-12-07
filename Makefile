up:
	docker-compose up -d
	docker-compose -f file-register/build/docker-compose.yml up -d
	docker-compose -f file-storage/build/docker-compose.yml up -d
	docker-compose -f subscription/build/docker-compose.yml up -d

stop:
	docker-compose stop
	docker-compose -f file-register/build/docker-compose.yml stop
	docker-compose -f file-storage/build/docker-compose.yml stop
	docker-compose -f subscription/build/docker-compose.yml stop
