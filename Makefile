build:
	go build cmd/main.go

clear:
	@rm -rf ./tmp/blocks
	@mkdir ./tmp/blocks

run:
	go run cmd/main.go

print:
	go run cmd/main.go print

add: 
	go run cmd/main.go add -block data 

