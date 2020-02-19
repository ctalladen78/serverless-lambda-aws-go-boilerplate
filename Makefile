.PHONY: build clean deploy

build:
	# dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/authorize authorize/main.go authorize/session_ctrl.go 
	env GOOS=linux go build -ldflags="-s -w" -o bin/todo_list todo_list/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/todo_new todo_new/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/signin signin/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/signup signup/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose
