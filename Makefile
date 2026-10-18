OS ?= $(shell go env GOOS)
ARCH ?= $(shell go env GOARCH)

IMAGE_NAME := "cert-manager-webhook-loopia"
IMAGE_TAG := "latest"

# etcd, kube-apiserver and kubectl for the conformance test, matching the
# k8s.io library versions in go.mod.
ENVTEST_VERSION ?= v0.25.2
ENVTEST_K8S_VERSION ?= 1.36.2
SETUP_ENVTEST := go run sigs.k8s.io/controller-runtime/tools/setup-envtest@$(ENVTEST_VERSION)

test:
	ASSETS="$(CURDIR)/$$($(SETUP_ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir testdata/bin -p path)" && \
	TEST_ASSET_ETCD="$$ASSETS/etcd" \
	TEST_ASSET_KUBE_APISERVER="$$ASSETS/kube-apiserver" \
	TEST_ASSET_KUBECTL="$$ASSETS/kubectl" \
	go test -v .

clean:
	chmod -R u+w testdata/bin 2>/dev/null; rm -rf testdata/bin

build:
	docker build -t "$(IMAGE_NAME):$(IMAGE_TAG)" .

.PHONY: test clean build rendered-manifest.yaml
rendered-manifest.yaml:
	helm template \
	    --name cert-manager-webhook-loopia \
        --set image.repository=$(IMAGE_NAME) \
        --set image.tag=$(IMAGE_TAG) \
        deploy/cert-manager-webhook-loopia > "$(OUT)/rendered-manifest.yaml"