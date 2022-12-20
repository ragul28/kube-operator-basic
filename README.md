# kube-operator-basic

## Getting Started for dev

* Create project using init command
```sh
mkdir -p kube-operator-basic
kubebuilder init --domain sample-operator.dev --repo github.com/ragul28/kube-operator-basic
```

* Create custom resource
```sh
kubebuilder create api --group basic --version v1 --kind BasicVolume
```

* Update api types & genarate manifests
```sh
make generate
make manifests
```

* Update the controller reconcile logic & install crd to cluster
```sh
make install
```

* Apply the Sample CRD
```sh
kubectl apply -f config/samples/
```

* Run controller in local
```sh
make run
```

### Running on the cluster
1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

2. Build and push your image to the location specified by `IMG`:
	
```sh
make docker-build docker-push IMG=<some-registry>/kube-operator-basic:tag
```
	
3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/kube-operator-basic:tag
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
