SHELL := /bin/bash

.PHONY: *

mocks:
	find . -type f -name 'mock_*.go' -delete
	mockery

image:
	docker build -t pilot:latest -t localhost:8562/pilot:latest .

setup_e2e: image
	k3d registry create skiff -p 0.0.0.0:8562 || true
	docker push localhost:8562/pilot:latest
	echo "Set env vars K3D_REGISTRY=k3d-skiff:8562 and PILOT_TEST_IMAGE=\$K3D_REGISTRY/pilot:latest"
