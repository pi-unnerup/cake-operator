# Rancher Muffin Time
Welcome to Muffin Time!

In this workshop we will modify, deploy and update the Kubernetes components going into making a cake:

- Operator
- Custom Resource Definition
- Custom Resource

## Prereqs

Prerecuisites for running this lab are: 

1. A default cluster that you can deploy resources to. `kubectl cluster-info` should show you your current context
2. Docker installed. `docker version` should return the version
3. Make installed. `make --version` should return the version
4. (optional, only if you want to build your own operator) golang v1.24. `go version` should return the version
5. (optional, only if you want to build your own operator) A Docker registry that you can push images to. You can create a free registry on [dockerhub](https://www.docker.com/products/docker-hub/)

## 1. Setup
## Rancher

If you are using your own cluster or a local cluster on your own machine such as Rancher or Docker Desktop, you will need to first install the operator:

```sh
make deploy IMG=unnerup/cake-operator:march
```

This will
- install kustomize in a local bin inside this repo
- deploy the operator into the cluster

You should see a number of resources created this way:
```sh
namespace/cake-operator-system created
serviceaccount/cake-operator-controller-manager created
role.rbac.authorization.k8s.io/cake-operator-leader-election-role created
clusterrole.rbac.authorization.k8s.io/cake-operator-manager-role created
clusterrole.rbac.authorization.k8s.io/cake-operator-metrics-reader created
clusterrole.rbac.authorization.k8s.io/cake-operator-proxy-role created
rolebinding.rbac.authorization.k8s.io/cake-operator-leader-election-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/cake-operator-manager-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/cake-operator-proxy-rolebinding created
configmap/cake-operator-manager-config created
service/cake-operator-controller-manager-metrics-service created
deployment.apps/cake-operator-controller-manager created
```

Later, if you want to uninstall the operator you can run:

```sh
make undeploy
```

## 2. Exploring the cluster

We can explore the resources available to us with:

```sh
kubectl api-resources
```

You should see a long list of resource types listed, such as pods, services and deployments. Cake, however, appears to be missing. Let's install the CustomResourceDefinition for Cake so we can start working with Cake objects. 

```sh
./bin/kustomize build config/crd | kubectl apply -f -
```

Now we can see our newly created resource type by running:

```sh
kubectl api-resources
```

A little later we will explore the content of `config/crd` by adding some features to our cake objects.

At last we are ready to deploy our first cake!

```sh
kubectl apply -f config/samples/tutorials_v1_cake.yaml
```

Your cake resource has created in the same way as any other resource and you can view it by:

```sh
kubectl get cake
```

The cake object also created some other resources:

```sh
kubectl get pods
kubectl get services
```

Now go and have a look at your creation! Since we are running rancher, we will need to expose a port on our localhost. We do this by [toggling the port-forwarder](https://docs.rancherdesktop.io/ui/port-forwarding/) in Rancher.

Now go to your browser on `localhost:<port>` where the port is the one you chose in rancher.

Your cake should look something like this:
![cake rendering](./img/plain-cake.png)

Congratulations on baking your first cake!

## 3. Expanding the recipe

We now have our first cake, but it is a little bit boring - it doesn't have any toppings for example. Lets spice up our cake. 

Go to `config/samples/tutorials_v1_cake.yaml` and update the `DECORATION` value from "no" one of the other available options, for example "ghost". Apply the changes:

NOTE: THIS WILL NOT DO ANYTHING

```sh
kubectl apply -f config/samples/tutorials_v1_cake.yaml
```

Nothing changed! The DECORATION value did not trigger the cake pod to update. One of the purposes of the operator is to define the controller loop for what triggers "reconciliation" or in other words when the controller should act to set the world right. The `port` and `replicaCount` specs have reconciliation control loops (you can find it in controllers/cake_controller.go), but the `DECORATION`, `COLOUR`, `MESSAGE` and `BACKGROUND` do not.

Some example reconciliation code:

```sh
if existingService.Spec.Ports[0].NodePort != port {
        log.Info("🔁 Port number changes, update the service! 🔁")
        existingService.Spec.Ports[0].NodePort = port
        err = r.Update(ctx, existingService)
        if err != nil {
            log.Error(err, "❌ Failed to update Service", "Service.Namespace", existingService.Namespace, "Service.Name", existingService.Name)
            return ctrl.Result{}, err
        }
    }
```

We now have three options: we can write a reconciliation block for DECORATION and reinstall the operator, we can change one of the values with an reconcillication loop (for example the port) or we can delete the cake resource and recreate it - in the interest of time, we will do the latter, but feel free to give it a go yourself.

```sh
kubectl delete cake cake-sample
kubectl apply -f config/samples/tutorials_v1_cake.yaml
```

When you refresh your browser you should see that your cake now has decorations.

Play around with different values for all of your cake specs. Note that BACKGROUND can be both hashes and names of colours.

You can submit your cake now (see section 5), or move onto the next section where we will update the operator to include an additional spec.

## 4. Updating the Operator

We can work with our cake object to change the decoration, colouration, background colour and message, but we would like to be able to update the title as well. This is not a current function supported by our existing operator. Let's add it!

For this section, it is advantageous to know a little bit of golang code, but not a requirement as we will walk through everything in stages. 

You will be working with 

- ./controller/cake_controller.go
- ./api/v1/cake_types.go

Open `./controller/cake_controller.go` and find the section which defines the required container spec. It will look like this:

```sh
Env: []corev1.EnvVar{
    {
        Name:  "COLOUR",
        Value: cakeCR.Spec.COLOUR,
    },
    {
        Name:  "MESSAGE",
        Value: cakeCR.Spec.MESSAGE,
    },
    {
        Name:  "DECORATION",
        Value: cakeCR.Spec.DECORATION,
    },
    {
        Name:  "BACKGROUND",
        Value: cakeCR.Spec.BACKGROUND,
    },
},
```

Add in an additional spec to the list as such:

```sh
    {
        Name:  "TITLE",
        Value: cakeCR.Spec.TITLE,
    },
```

Save and close the file. 
Open `./api/v1/cake_types.go` and find the section which defines the structure of the cake spec. It will look like this:

```sh
type CakeSpec struct {
    // Number of replicas for the Nginx Pods
    ReplicaCount int32 `json:"replicaCount"`
    // Exposed port for the Nginx server
    Port int32 `json:"port"`
    //COLOUR can be one of "white" or "colour"
    COLOUR string `json:"COLOUR"`
    //Decoration can be one of "ghost" or "heart"
    DECORATION string `json:"DECORATION"`
    //Background colour, e.g. Aquamarine
    BACKGROUND string `json:"BACKGROUND"`
    MESSAGE    string `json:"MESSAGE"`
}
```

Add in an additional spec to the bottom of the list as such:

```sh
    TITLE    string `json:"TITLE"`
```

You will now need to build and push the image - for this you will need your own registry. If you do not have a registry easily available, we have already built a ready image for this: `unnerup/cake-operator:TITLE`, otherwise build and push your operator image:

NOTE: This can take a few minutes.

```sh
docker build -t <your-registry>/cake-operator:TITLE .
docker push <your-registry>/cake-operator:TITLE
```

Now deploy whichever operator image you choose to use, using either pre-built unnerup image or your own registry:

```sh
make deploy IMG=unnerup/cake-operator:TITLE
```

or

```sh
make deploy IMG=<your-registry>/cake-operator:TITLE
```

Let's now edit our CustomResourceDefinition to reflect our changes. Open `./config/crd/bases/tutorials.piunnerup.com_cakes.yaml`. Under `spec.versions.schema.openAPIV3Schema.properties.spec.properties` you can see entries for BACKGROUND, COLOUR, DECORATION, MESSAGE, port and replicaCount. Add an entry for TITLE as such:

```sh
TITLE:
  type: string
```

Underneath this section you will see a required section. Add `- TITLE` to here as well. 

Save the file and apply it:

```sh
kubectl apply -f ./config/crd/bases/tutorials.piunnerup.com_cakes.yaml
```

Now we are ready to add a TITLE to our cake object. Open your cake resource again `config/samples/tutorials_v1_cake.yaml` and add a `TITLE` to the spec field alongside the other specs. Title your cake whatever you would like. Then delete and apply your cake object again.

```sh
kubectl delete cake cake-sample
kubectl apply -f ./config/crd/bases/tutorials.piunnerup.com_cakes.yaml
```

View your newly titled cake in your browser. 
