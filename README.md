mkdir cake-operator && cd cake-opertaor
operator-sdk init --project-name cake-operator --domain piunnerup.com --repo github.com/pi-unnerup/cake-operator
operator-sdk create api --group tutorials --version v1 --kind Cake --resource --controller
make manifests

replace in api/v1/cake_types.go in CakeSpec struct
```
    // Number of replicas for the Nginx Pods
    ReplicaCount int32 `json:"replicaCount"`
    // Exposed port for the Nginx server
    Port int32 `json:"port"`
```

make manifests

update ./controllers/cake_controller.go: https://docs.ovh.com/gb/en/kubernetes/deploying-go-operator/#:~:text=Update%20the%20./controllers/ovhnginx_controller.go%20file%3A

cake_controller.go change image to unnerup/muffin-time:1

make install run //runs on local, don't kill process. Can also run make install for not local

in ./config/samples/tutorials_v1_cake.yaml add
```
spec:
  port: 80
  replicaCount: 1
```

//kubectl create ns test-go-operator
kubectl apply -f ./config/samples/tutorials_v1_cake.yaml //-n test-go-operator //takes some time to come up

Go to EXTERNAL-IP:NODEPORT (on rancher, localhost:8099)



# nginx-go-operator
// TODO(user): Add simple overview of use/purpose

## Description
// TODO(user): An in-depth paragraph about your project and overview of use

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster
1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

2. Build and push your image to the location specified by `IMG`:
	
```sh
make docker-build docker-push IMG=<some-registry>/nginx-go-operator:tag
```
	
3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/nginx-go-operator:tag
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller to the cluster:

```sh
make undeploy
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/) 
which provides a reconcile function responsible for synchronizing resources untile the desired state is reached on the cluster 

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

