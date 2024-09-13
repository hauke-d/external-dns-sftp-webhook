EXECUTABLE=build/webhook

.PHONY: run
run: build
	./$(EXECUTABLE)

.PHONY: build
build:
	go build -o ./$(EXECUTABLE)
