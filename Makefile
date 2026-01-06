start-agent-a:
	go run main.go agent-a
start-agent-b:
	go run main.go agent-b
start-agent-root:
	go run main.go root-agent web api webui
test:
	./scripts/submit.sh

# tracing
start-lgtm:
	docker compose -f compose-otel.yaml up -d
start-beyla:
	sudo beyla --config beyla-config.yaml
