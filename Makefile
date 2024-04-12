
build:
	go build -o ./bin/api ./cmd/api
	go build -o ./bin/scraper ./cmd/scraper
	go build -o ./bin/parser ./cmd/parser
	go build -o ./bin/taxonomizer ./cmd/taxonomizer
	go build -o ./bin/maintainer ./cmd/maintainer
	go build -o ./bin/natspub ./cmd/natspub