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


run-ten-nodes:
	./main -self-addr "localhost:8000" -self-node-id "a" -seed-node=true -api-at ":9000" & \
	./main -self-addr "localhost:8001" -self-node-id "b" -peer-addr "localhost:8000" -peer-node-id "a" -delay 1 -api-at ":9001" & \
	./main -self-addr "localhost:8002" -self-node-id "c" -peer-addr "localhost:8001" -peer-node-id "b" -delay 2 -api-at ":9002" & \
  ./main -self-addr "localhost:8003" -self-node-id "d" -peer-addr "localhost:8000" -peer-node-id "a" -delay 3 -api-at ":9003" & \
  ./main -self-addr "localhost:8004" -self-node-id "e" -peer-addr "localhost:8002" -peer-node-id "c" -delay 4 -api-at ":9004" & \
	./main -self-addr "localhost:8005" -self-node-id "f" -peer-addr "localhost:8004" -peer-node-id "e" -delay 5 -api-at ":9005" & \
	./main -self-addr "localhost:8006" -self-node-id "g" -peer-addr "localhost:8003" -peer-node-id "d" -delay 6 -api-at ":9006" & \
  ./main -self-addr "localhost:8007" -self-node-id "h" -peer-addr "localhost:8005" -peer-node-id "e" -delay 7 -api-at ":9007" & \
  ./main -self-addr "localhost:8008" -self-node-id "i" -peer-addr "localhost:8007" -peer-node-id "h" -delay 8 -api-at ":9008"
	./main -self-addr "localhost:8009" -self-node-id "j" -peer-addr "localhost:8007" -peer-node-id "h" -delay 10 -api-at ":9009"
	wait


kill-10pro:
	sudo lsof -ti:8000,8001,8002,8003,8004,8005,8006,8007,8008,8009 | xargs kill -9

