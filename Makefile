.PHONY: build
build:
	go build -o ./user/cmd -v ./user/cmd/main.go 
	go build -o ./auth/cmd -v ./auth/cmd/main.go
	go build -o ./notifications/cmd -v ./notifications/cmd/main.go
	go build -o ./support/cmd -v ./support/cmd/main.go	
	go build -v .

.PHONY: user
user:
	go build -o ./user/cmd -v ./user/cmd/main.go 

.PHONY: core
core:
	go build -v .

.PHONY: auth
auth:
	go build -o ./auth/cmd -v ./auth/cmd/main.go

.PHONY: notifs
notifs:
	go build -o ./notifications/cmd -v ./notifications/cmd/main.go

.PHONY: support
support:
	go build -o ./support/cmd -v ./support/cmd/main.go	

.PHONY: test
test:
	go test -v race -timeout 30s ./...


.DEFAULT_GOAL := build