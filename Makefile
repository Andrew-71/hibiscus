build:
	go build

run:
	go build && ./hibiscus-txt

dev:
	go build && ./hibiscus-txt  --config config/dev-config.txt