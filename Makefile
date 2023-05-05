run:
	go run ./main.go -o mongo
print:
	go run ./main.go -o print
build:
	go build -o ./bin/main.exe ./main.go
exeRun:
	./bin/main.exe -o mongo
exePrint:
	./bin/main.exe -o print