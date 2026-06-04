run:
	cd ~/Desktop/projects/DukeDB/cmd/duke/ && go run .

compile:
	go build ./cmd/duke/main.go

run-two-nodes:
	./main -self-addr "localhost:8000" -self-node-id "a" -peer-addr "localhost:8001" -peer-node-id "b" -delay 3  & \
	./main -self-addr "localhost:8001" -self-node-id "b" -peer-addr "localhost:8000" -peer-node-id "a" -delay 8 & \
	wait

kill-2pro:
	sudo lsof -ti:8000,8001 | xargs kill -9

kill-3pro:
	sudo lsof -ti:8000,8001,8002 | xargs kill -9

restart:
	make kill-pro ; make compile && make run-two-nodes

run-three-nodes:
	./main -self-addr "localhost:8000" -self-node-id "a" -seed-node true & \
	./main -self-addr "localhost:8001" -self-node-id "b" -peer-addr "localhost:8000" -peer-node-id "a" -delay 8 & \
	./main -self-addr "localhost:8002" -self-node-id "c" -peer-addr "localhost:8000" -peer-node-id "a" -delay 5
	wait
