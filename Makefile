Run:
	go run .

run:
	go run .

Preload:
	go run . -l

preload:
	go run . -l

MigrateUp:
	migrate -database "sqlite://./Database/Database.sqlite" -path ./Database/Migrations up

MigrateDown:
	migrate -database "sqlite://./Database/Database.sqlite" -path ./Database/Migrations down

Test:
	bash ./runTests.sh

ServBuild:
	go build -o startBack

RunBuild:
	./startBack

RunBuildPreload:
	./startBack -l

DockerBuild:
	docker build -t social-back .

DockerRun:
	docker run -p 8080:8080 -it social-back