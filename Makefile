# import deploy config
# You can change the default deploy config with `make cnf="deploy_special.env" release`
dpl ?= deploy.env
include $(dpl)
export $(shell sed 's/=.*//' $(dpl))

dpl_root ?= ../deploy.env
include $(dpl_root)
export $(shell sed 's/=.*//' $(dpl_root))

include ../Makefile


.PHONY: build
build:
	go mod tidy
	go mod vendor
	#go test -v ./... -gcflags='-l'
	#cd cmd/tenantDeployer && CGO_ENABLED=0 go build -o ../artifacts/${BINARY_NAME}

	cd cmd && CGO_ENABLED=0 env GOOS=linux GOARCH=amd64 go build -o ../artifacts/${BINARY_NAME}

#for official branches
build-fetcher:
	go mod tidy
	go mod vendor
	cd cmd && CGO_ENABLED=0 go build -o ../artifacts/${BINARY_NAME}

image: build
	docker build -f build/Dockerfile -t ${DOCKER_REPO_PATH}/${DOCKER_REPO_NAME}:latest .
	docker tag ${DOCKER_REPO_PATH}/${DOCKER_REPO_NAME}:latest ${DOCKER_REPO_PATH}/${DOCKER_REPO_NAME}:${DOCKER_TAG}
#	docker tag ${DOCKER_REPO_PATH}/${DOCKER_REPO_NAME}:latest ${DOCKER_REPO_PATH}/${DOCKER_REPO_NAME}:${BINARY_NAME}-${BUILD_DATE}

public-image: auth
	docker push ${DOCKER_REPO_PATH}/${DOCKER_REPO_NAME}:latest




test:
	 CGO_ENABLED=0 env GOOS=linux GOARCH=amd64 go test internal/pkg/dd_deployer/*



