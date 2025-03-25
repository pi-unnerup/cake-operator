## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Prerequisites

Prerecuisites for running this lab are: 

1. A default cluster that you can deploy resources to. `kubectl cluster-info` should show you your current context
2. Docker installed. `docker version` should return the version
3. Make installed. `make --version` should return the version
4. golang v1.24. `go version` should return the version
5. A Docker registry that you can push images to. You can create a free registry on [dockerhub](https://www.docker.com/products/docker-hub/)

### Setup

The operator functionality runs as a container inside the cluster.

Build and push your image to the location specified by `IMG`:
	
```sh
make docker-build IMG=<some-registry>/cake-operator:tag
```

```sh
make docker-push IMG=<some-registry>/cake-operator:tag
```

Create your manifests
```sh
make manifests
```

`ls -la config/crd/bases/` will show you the timestamp of the newly created manifest

### Deploy the operator and Custom Resource Definitiion

Deploy the operator to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/cake-operator:tag
make install
```
Notice that lots of Kubernetes Resources have been created as part of the operator. 

You should be able to confirm the CRD is installed by 

```sh
kubectl api-resources -n cake-operator-system
```

where your 'Cakes' resource type now shows up.

### Bake the cake

Now it's time to create your cake object using your operator and CRD created in the previous steps:

```sh
kubectl apply -f config/samples/tutorials_v1_cake.yaml
```

Your cake is now up and running! `kubectl describe cake cake-sample` confirms some details about your cake such as that it is running on nodeport 30300. Please note if you're running in a local cluster the 'node' may not be your local machine and so in order to view your cake in the browser please forward it to your localhost. [Here](https://docs.rancherdesktop.io/ui/port-forwarding/) is a neat guide on how to do it on Rancher Desktop in 10 seconds. 

Go to your browser and view your cake; localhost:30300 (or whichever port you assigned on your node)

### Customise the cake

Change your cake object `config/samples/tutorials_v1_cake.yaml` by adding some custom configuration, such as a colour and a decoration, then delete and apply the object again. 

```sh
kubectl delete -f config/samples/tutorials_v1_cake.yaml
kubectl apply -f config/samples/tutorials_v1_cake.yaml
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


### Modifying the API definitions

If you are editing the API definitions, don't forget to generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

This operator was built using [this guide](https://medium.com/developingnodes/mastering-kubernetes-operators-your-definitive-guide-to-starting-strong-70ff43579eb9)

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

