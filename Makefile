run:
	go run main.go

run-background:
	nohup go run main.go &

stop-background:
	sudo kill $$(sudo lsof -t -i:8080)