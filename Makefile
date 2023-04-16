.PHONY: render
render:
	dot -Tpng -Ksfdp -O scheme.gv

.PHONY: run
run:
	go run ./cmd/app/main.go