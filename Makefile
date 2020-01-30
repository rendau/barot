.DEFAULT_GOAL := build

BINARY_NAME = barot_api
BUILD_PATH = cmd/build

build:
	mkdir -p $(BUILD_PATH)
	cp .conf.yml $(BUILD_PATH)/
	CGO_ENABLED=0 go build -o $(BUILD_PATH)/$(BINARY_NAME) main.go

run:
	docker-compose up -d

stop:
	docker-compose down

test:
	docker-compose -f docker-compose.test.yml up --build --exit-code-from barot_tests
	@EXIT_CODE=$$?
	docker-compose -f docker-compose.test.yml down
	@exit $$EXIT_CODE

clean:
	rm -rf $(BUILD_PATH)
