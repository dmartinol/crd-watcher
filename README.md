# crd-watcher
> This reference implementation is based on the [sample-controller](https://github.com/kubernetes/sample-controller) example

## Scope
Define a pattern to build event-driven applications in a Kubernetes environment using `CustomResource`s to model events.

## Overview
This application defines an asynchronous system of two components that interact by exchanging `CustomResource` instances: 
* `director` is the orchestrator that starts the execution of a sample flow
* `worker` is the executor of flows

Every execution is modelled by an instance of the `RequestState` resource, a resource defined in [crds/request-state.yaml](./crds/request-state.yaml)
that defines the status of a request, together with other attributes that specify which `job` has to be performed:
```yaml
apiVersion: com.redhat.ecosystem.sample/v1
kind: RequestState
metadata:
  name: example
  namespace: default
spec:
  request-uid: aaa-bbb-ccc
  job: JOB2
  state: COMPLETED
status:
  history:
    - job: JOB1
      state: STARTED
      timestamp: t0
    - job: JOB1
      state: COMPLETED
      timestamp: t1
    - job: JOB2
      state: STARTED
      timestamp: t2
    - job: JOB2
      state: COMPLETED
      timestamp: t3
```
Different jobs can be modelled
In the `status.history` field we can track the full list of state changes for a given request.

## Sample flow
The flow implemented by this reference implementation is the following:
* `director` creates a `RequestState` with `job='JOB1', state=REQUESTED`
* `worker` updates it to `job=JOB1, state=STARTED`
* `worker` updates it to `job=JOB1, state=COMPLETED`
* `director` updates it to `job=JOB2, state=REQUESTED`
* `worker` updates it to `job=JOB2, state=STARTED`
* `worker` updates it to `job=JOB2, state=COMPLETED`
* `director` receives `job='JOB2', state=COMPLETED` and stops the flow


## Generating the K8s client code
This command generates the Go client fcode to manage the [RequestState](./crds/request-state.yaml) `CustomResourceDefinition`:
```yaml
./hack/update-codegen.sh
```
You don't need to run it unless you changes anything in the `CustomResourceDefinition`. 

If you applied any change, remember to update the 
[RequestState](./pkg/apis/requeststate/v1/types.go) type definition accordingly.

## Installing the CRDs
Install the CRD with an [example instance](./crds/sample.yaml), verify it was created and finally delete it:
```bash
oc apply -f crds
oc get RequestState -oyaml -A
oc delete RequestState -A --all
```

## Running the example
Runtime options:
* `-mode=director/worker` (def is director)
* `-namespace=NAMESPACE` (def is `default`)
* `-kubeconfig=PATH_TO_KUBECONFIG` (like: `-kubeconfig="$HOME"/.kube/config`)

Run the `worker` application:
```bash
go run main.go -kubeconfig="$HOME"/.kube/config -mode worker -namespace test-crd
```

Run the `director` application:
```bash
go run main.go -kubeconfig="$HOME"/.kube/config -mode director -namespace test-crd
``

Check the status of generated resources:
```bash
oc get RequestState -oyaml -n test-crd
```

Clean up:
```bash
oc delete RequestState -n test-crd --all
```

## Useful references
* [How to Create a Kubernetes Custom Controller Using client-go](https://itnext.io/how-to-create-a-kubernetes-custom-controller-using-client-go-f36a7a7536cc)
* [How to generate client codes for Kubernetes Custom Resource Definitions (CRD)](https://itnext.io/how-to-generate-client-codes-for-kubernetes-custom-resource-definitions-crd-b4b9907769ba)

