run:
	cd ~/Desktop/projects/DukeDB/cmd/duke/ && go run .

compile:
	go build ./cmd/duke/main.go

compile-with-race:
	go build -race ./cmd/duke/main.go 

run-two-nodes:
	./main -self-addr "localhost:8000" -self-node-id "a" -seed-node true -delay 1 -api-at ":9000"  & \
		./main -self-addr "localhost:8001" -self-node-id "b" -peer-addr "localhost:8000" -peer-node-id "a" -delay 8 -api-at ":9001" & \
	wait

kill-2pro:
	sudo lsof -ti:8000,8001 | xargs kill -9

kill-3pro:
	sudo lsof -ti:8000,8001,8002 | xargs kill -9

restart:
	make kill-pro ; make compile && make run-two-nodes

run-three-nodes:
	./main -self-addr "localhost:8000" -self-node-id "a" -seed-node=true -api-at ":9000" & \
		./main -self-addr "localhost:8001" -self-node-id "b" -peer-addr "localhost:8000" -peer-node-id "a" -delay 2 -api-at ":9001" & \
		./main -self-addr "localhost:8002" -self-node-id "c" -peer-addr "localhost:8000" -peer-node-id "a" -delay 5 -api-at ":9002"
	wait

run-four-nodes:
	./main -self-addr "localhost:8000" -self-node-id "a" -seed-node=true -api-at ":9000" & \
	./main -self-addr "localhost:8001" -self-node-id "b" -peer-addr "localhost:8000" -peer-node-id "a" -delay 1 -api-at ":9001" & \
	./main -self-addr "localhost:8002" -self-node-id "c" -peer-addr "localhost:8001" -peer-node-id "b" -delay 2 -api-at ":9002" & \
  ./main -self-addr "localhost:8003" -self-node-id "d" -peer-addr "localhost:8000" -peer-node-id "a" -delay 5 -api-at ":9003"
	wait

run-five-nodes:
	./main -self-addr "localhost:8000" -self-node-id "a" -seed-node=true -api-at ":9000" & \
	./main -self-addr "localhost:8001" -self-node-id "b" -peer-addr "localhost:8000" -peer-node-id "a" -delay 1 -api-at ":9001" & \
	./main -self-addr "localhost:8002" -self-node-id "c" -peer-addr "localhost:8001" -peer-node-id "b" -delay 2 -api-at ":9002" & \
  ./main -self-addr "localhost:8003" -self-node-id "d" -peer-addr "localhost:8000" -peer-node-id "a" -delay 5 -api-at ":9003" & \
  ./main -self-addr "localhost:8004" -self-node-id "e" -peer-addr "localhost:8002" -peer-node-id "c" -delay 6 -api-at ":9004"
	wait

kill-5pro:
	sudo lsof -ti:8000,8001,8002,8003,8004 | xargs kill -9

