run:
	cd ~/Desktop/projects/DukeDB/cmd/duke/ && go run .

compile:
	go build ./cmd/duke/main.go

run-two-nodes:
	./main -selfAddr "localhost:8000" -peerAddr "localhost:8001" -peerNodeID "b" & \
	./main -selfAddr "localhost:8001" -peerAddr "localhost:8000" -peerNodeID "a" & \
	wait

kill-pro:
	sudo lsof -ti:8000,8001 | xargs kill -9
