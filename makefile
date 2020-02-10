mod:
	# This make rule requires Go 1.11+
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor

test:
	docker-compose -f integration_tests.yml up --build --abort-on-container-exit
