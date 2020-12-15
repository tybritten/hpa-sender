# HorizontalPodScaler Sending

### Prerequisites
* Knative Eventing Installed and setup with a Broker
* Knative ApiServerSource setup
* ServiceAccount for the pod that has API access to secrets and HorizontalPodScalers. 

An example Role for the ServiceAccount:
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: default
  name: ce-webhook-sender
rules:
- apiGroups: [""] # "" indicates the core API group
  resources: ["secrets", "horizonalpodautoscalers"]
  verbs: ["get", "watch", "list"]

---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: ce-webhook-sender-binding
  namespace: default
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: ce-webhook-sender
subjects:
- kind: ServiceAccount
  name: default
  namespace: default
```


To use, run in cluster, expose a service on port 8080. 

Create a Knative trigger pointing at this service:
```yaml
apiVersion: eventing.knative.dev/v1
kind: Trigger
metadata:
  name: my-service-trigger
spec:
  broker: default
  subscriber:
    ref:
      apiVersion: serving.knative.dev/v1
      kind: Service
      name: my-service

```
Once this is setup this service will watch for HorizontalPodAutoscaler Events and if it sees one, it will check that HPA for an annotation-
```yaml
hpa-event: namespace/secretname
```
And if this secret exists, it will send an event to your configured webhook from the secret. The secret format lis like this:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: secretname
type: Opaque
stringData:
  url: https://mywebookdestination
  headers: >
      {
        "api-key": "key here",
        "some_other": "traceheader"
      }
  body: >
    {
      "type": "scaling",
      "time": "_time_",
      "message": "_message_"
    }
```

The webhook call will be a `POST` to the `url` you supplied, and will include any headers you have in `headers:` and the body that's in `body:` (JSON only). There are two special replacement variables `_time_` and `_message_`- if you use them in your body, the service will replace them with the time and/or the message from the HPA scaling event.

