# Build Instructions

## To build
go build bitbucket.org/project-alkemi/switch-board-service/cmd/switch-board-service

## To Run
go run bitbucket.org/project-alkemi/switch-board-service/cmd/switch-board-service

## To Test
go test bitbucket.org/project-alkemi/switch-board-service/test/switch-board-service

## To Clean
go clean bitbucket.org/project-alkemi/switch-board-service/cmd/switch-board-service

## To Generate Documentation
swag init -g cmd/switch-board-service/main.go --output docs/swagger/
