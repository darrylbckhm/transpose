build:
	go build -o lib/go/orchestrators/ingestion/crawler/main -modfile lib/go/orchestrators/ingestion/crawler/go.mod lib/go/orchestrators/ingestion/crawler/main.go
	go build -o lib/go/orchestrators/ingestion/fetch/main -modfile lib/go/orchestrators/ingestion/fetch/go.mod lib/go/orchestrators/ingestion/fetch/main.go

run:
	go run lib/go/orchestrators/ingestion/fetch/main.go
	go run lib/go/orchestrators/ingestion/crawler/main.go

compile:
	# 32-Bit Systems
	# FreeBDS
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 lib/go/orchestrators/ingestion/fetch/main.go
	# MacOS
	GOOS=darwin GOARCH=386 go build -o bin/main-darwin-386 lib/go/orchestrators/ingestion/fetch/main.go
	# Linux
	GOOS=linux GOARCH=386 go build -o bin/main-linux-386 lib/go/orchestrators/ingestion/fetch/main.go
	# Windows
	GOOS=windows GOARCH=386 go build -o bin/main-windows-386 lib/go/orchestrators/ingestion/fetch/main.go
	# 64-Bit
	# FreeBDS
	GOOS=freebsd GOARCH=amd64 go build -o bin/main-freebsd-amd64 lib/go/orchestrators/ingestion/fetch/main.go
	# MacOS
	GOOS=darwin GOARCH=amd64 go build -o bin/main-darwin-amd64 lib/go/orchestrators/ingestion/fetch/main.go
	# Linux
	GOOS=linux GOARCH=amd64 go build -o bin/main-linux-amd64 lib/go/orchestrators/ingestion/fetch/main.go
	# Windows
	GOOS=windows GOARCH=amd64 go build -o bin/main-windows-amd64 lib/go/orchestrators/ingestion/fetch/main.go
