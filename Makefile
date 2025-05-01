.PHONY: wire run
wire:
	cd cmd && wire
run:
	cd cmd && go run main.go wire_gen.go