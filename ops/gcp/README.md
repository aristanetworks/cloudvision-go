# Let's Encrypt Arista

##
In order to use Let's Encrypt certificates internally, we need to have:

1. a publicly exposed container running a webserver (nginx in our case) to answer Let's Encrypt challenge requests and get the certificates for a domain.
2. Point the dns for this domain to this container public ip address.

1/ Is done by having a container with a static ip address running in a GCP kubernetes cluster.
The image definition of this container is here: TODO
The image is deploy to a Google Container Registry instance so the GCP k8s cluster can use it to start the container.

2/ is done by using `arista.io` top domain and having a dual dns:
The public dns will point to the container deployed in 1/
The arista internal dns will point to the kubernetes cluster having the Let's Encrypt certificates deployed.

## Building the letsencrypt docker image:

from this directory, run:
```sh
docker build -t gcr.io/sw-jenkins-build-storage/letsencrypt letsencrypt
```

## Pushing the letsencrypt docker image to gcp registry

```sh
gcloud docker -- push gcr.io/sw-jenkins-build-storage/letsencrypt
```

## Deploying the container

Note: Be sure to have your `kubectl` command point to the right kubernetes cluster `gke_sw-jenkins-build-storage_us-west1-a_gcops-cluster`.

```sh
kubectl apply -f letsencrypt.yml
```


### Notes

A public ip address needs to be created:

```sh
gcloud compute addresses create letsencrypt-public --region=us-west1
```

For the rest, the Service with LoadBalancerIP/Ingress etc, I'm not sure to understand how/why it works :-O Need to investigate. Maybe simplify all this.