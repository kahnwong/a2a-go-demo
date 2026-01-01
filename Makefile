start-agent-a:
	go run main.go agent-a
start-agent-b:
	go run main.go agent-b
start-agent-root:
	go run main.go root-agent web api webui
submit:
	./scripts/submit.sh
