run:
	cd ~/Desktop/projects/DukeDB/cmd/duke/ && go run .

compile:
	go build ./cmd/duke/main.go

run-two-nodes:
	./main -selfAddr "localhost:8000" -selfNodeID "a" -peerAddr "localhost:8001" -peerNodeID "b" -delay 3 & \
	./main -selfAddr "localhost:8001" -selfNodeID "b" -peerAddr "localhost:8000" -peerNodeID "a" -delay 8 & \
	wait

kill-pro:
	sudo lsof -ti:8000,8001 | xargs kill -9

restart:
	make kill-pro ; make compile && make run-two-nodes
