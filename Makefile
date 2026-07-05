.PHONY: build web test smoke smoke-api smoke-docker

build: web
	rm -rf cmd/disk-tool/static/*
	cp -r web/dist/* cmd/disk-tool/static/
	go build -o bin/disk-tool ./cmd/disk-tool

web:
	cd web && npm ci && npm run build

test:
	go test -race ./...

smoke:
	docker compose run --rm smoke

smoke-api:
	bash scripts/smoke-api.sh ./bin/disk-tool

smoke-docker:
	bash scripts/smoke-docker.sh
