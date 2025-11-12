start-agent-a:
	go run agent_a.go
start-agent-b:
	go run agent_b.go
start-agent-root:
	go run main.go web api webui
submit:
	./scripts/submit.sh
