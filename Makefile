GOCMD=go
GOBUILD=$(GOCMD) build
GOGET=$(GOCMD) get
BIN_NAME=validator
PKI_FOLDER=./pki
DOCKER_TAG=rcanderson23/validator:v0.0.1-06

.PHONY: build
build:
	$(GOBUILD) -o $(BIN_NAME) -v

.PHONY: build-run
build-run:
	$(GOBUILD) -o $(BIN_NAME) -v; ./${BIN_NAME}

.PHONY: standalone
standalone:
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) -a -ldflags '-w -extldflags "-static"' -o $(BIN_NAME) -v

.PHONY: gen-pki
gen-pki:
	mkdir -p ${PKI_FOLDER}
	openssl req -nodes -new -x509 -keyout ${PKI_FOLDER}/ca.key -out ${PKI_FOLDER}/ca.crt -subj "/CN=Validator Webhook"
	openssl genrsa -out ${PKI_FOLDER}/server.key 2048
	openssl req -new -key ${PKI_FOLDER}/server.key -subj "/CN=validator.kube-system.svc" | openssl x509 -req -CA ${PKI_FOLDER}/ca.crt -CAkey ${PKI_FOLDER}/ca.key -CAcreateserial -out ${PKI_FOLDER}/server.crt

.PHONY: gen-container
gen-container:
	docker build -t=${DOCKER_TAG} .

.PHONY: push-container
push-container:
	docker push ${DOCKER_TAG}

.PHONY: deploy
deploy:
	kubectl create secret generic validator-tls --from-file=${PKI_FOLDER}/server.crt --from-file=${PKI_FOLDER}/server.key
