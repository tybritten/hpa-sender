# Example using Kind

### Prerequisites
* Have [Kind](https://kind.sigs.k8s.io) installed
* Create the cluster and install all the knative eventing components and the metrics server (needed for the HPA) by running `./eventing.sh`.

Setup the API Event Sender
`kubectl apply -f k8s-events.yaml`

Now Deploy the HPA-sender
`kubectl apply -f hpa-sender-deployment.yaml`

Now we're ready to deploy a basic test app for HPA.

Edit secret at the begining of the `hpa-app.yaml` file to set your URL and headers/body for the webhooks.

Then save and apply it
`kubectl apply -f hpa-app.yaml`

To check on the app
'kubectl get pods'
and the HPA
`kubectl get hpa`

Now we can cause the load to go up and down by following the instructions [here](https://unofficial-kubernetes.readthedocs.io/en/latest/tasks/run-application/horizontal-pod-autoscale-walkthrough/)


```bash
$ kubectl run -i --tty load-generator --image=busybox /bin/sh

Hit enter for command prompt

$ while true; do wget -q -O- http://php-apache.default.svc.cluster.local; done
```

Which should trigger events to your url. 
