Run:
	go run .

run:
	go run .

MigrateUp:
	migrate -database "sqlite://./Database/Database.sqlite" -path ./Database/Migrations up

MigrateDown:
	migrate -database "sqlite://./Database/Database.sqlite" -path ./Database/Migrations down

Test:
	bash ./runTests.sh