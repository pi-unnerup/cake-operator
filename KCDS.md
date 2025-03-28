# KCDS Muffin Time
Welcome to Muffin Time!

In this workshop we will modify, deploy and update the Kubernetes components going into making a cake:

- Operator
- Custom Resource Definition
- Custom Resource

## 1. Setup
## If you are using the Workshop Provided cluster

To get a cluster, go to
https://tinyurl.com/yxrk6ynu
and download a config file.

In order to start using the file, set the KUBECONFIG variable on your terminal to the path to the file you just downloaded.

```sh
export KUBECONFIG=<path/to/file>
```

Before starting, make sure that your CLI is pointing to the right cluster:

```sh
kubectl config current-context
```

The name of your current context should be "kubernetes-admin@kubernetes". If the context reads anything else, you should ensure that your KUBECONFIG is set to the path of the file provided by your trainer 

Lastly we will install the operator. The image is the same as the code in this repo - we will go through some of it later:

```sh
make deploy IMG=unnerup/cake-operator:v1.0.1
```

## If you are using your own cluster

If you are using your own cluster or a local cluster on your own machine such as Rancher or Docker Desktop, you will need to first install the operator:

```sh
make deploy IMG=unnerup/cake-operator:v1.0.1
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

Now go and have a look at your creation! If you are running on the workshop-provided cluster go to your browser and see your cake on `<node-IP>:30300` or `<node-dns:30300>
You can find your external IP by looking at your kubeconfig and finding the server. 
If you are using the workshop provided clusters, you can also find the node address as such:

```sh
kubectl config view -o jsonpath='{.clusters[].cluster.server}' | sed 's~https://~~g' | sed 's/:6443/\n/g'
```

Note that the cake will run on http, e.g. `ec2-18-130-139-178.eu-west-2.compute.amazonaws.com:30300`

If you are running on a local cluster such as Rancher, you will likely see it on `localhost:80`.

Your cake should look something like this:
![cake rendering](./img/plain-cake.png)

Now lets change the port of our cake to 30400. In `config/samples/tutorials_v1_cake.yaml` change the port to 30400 and apply your cake again:

```sh
kubectl apply -f config/samples/tutorials_v1_cake.yaml
```

Your cake will now be available on port 30400.
Congratulations on baking your first cake!

## 3. Expanding the recipe

We now have our first cake, but it is a little bit boring - it doesn't have any toppings for example. Lets spice up our cake. 

Go to `config/samples/tutorials_v1_cake.yaml` and update the `DECORATION` value from "no" one of the other available options, for example "ghost". Apply the changes:

NOTE: THIS WILL NOT DO ANYTHING

```sh
kubectl apply -f config/samples/tutorials_v1_cake.yaml
```

Nothing changed! When we changed the port, the object updated, but the DECORATION value did not trigger the cake pod to update. One of the purposes of the operator is to define the controller loop for what triggers "reconciliation" or in other words when the controller should act to set the world right. The `port` and `replicaCount` specs have reconciliation control loops (you can find it in controllers/cake_controller.go), but the `DECORATION`, `COLOUR`, `MESSAGE` and `BACKGROUND` do not.

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

We now have two options: we can write a reconciliation block for DECORATION and reinstall the operator or we can delete the cake resource and recreate it - in the interest of time, we will do the latter, but feel free to give it a go yourself.

```sh
kubectl delete cake cake-sample
kubectl apply -f config/samples/tutorials_v1_cake.yaml
```

When you refresh your browser you should see that your cake now has decorations.

Play around with different values for all of your cake specs. Note that BACKGROUND can be both hashes and names of colours.

You can submit your cake now (see section 5), or move onto the next section where we will update the operator to include an additional spec.

## 4. (Optional) Updating the Operator

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

## 5. Submit your cake

Time to submit your cake! To enter into the competition for 3 x £25 AWS giftcards, take a screenshot of your cake and attach it here https://tinyurl.com/2v8t4dyv
You can submit until 5pm today, at which point we will pick 3 winning designs and email the winners. Don't forget to have one of the trainers scan your badge on the way out to enter!
