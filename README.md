# limits-admission-webhook

A k8s admission webhook to validate deployments

## Prepare

* `kubectl apply -f manifests/deployment`
* `kubectl apply -f manifests/admission-webhook.yaml`

## Test

This command must fail with `... denied the request: no resource limits set`
* `kubectl apply -f manifests/invalid-deployment.yaml`