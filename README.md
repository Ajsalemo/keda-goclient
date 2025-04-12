# keda-goclient

An application functioning as a REST API that does the following:
- Uses `client-go` to interface with Kubernetes. This has been developed in `wsl2` with a local Docker enabled Kubernetes single node cluster
- Acts as a client to create KEDA `ScaledJobs` and `ScaledObjects`
- `scaler-kubernetes-webhook` `.yamls` can be deployed to set up both Mutating and Admission Webhooks that will be invoked during KEDA resource creation. The logic currently is that a `metadata.label` of either `scaledjob` or `scaledobject` will be adding during resource creation by the Mutating webhook


----------
Below is the request body structure that some of these endpoints accept. `POST` bodies should follow the same general structure that a Kubernetes or KEDA `.yaml` expects

`/api/secret/create`

```json
{
    "name": "yoursecret",
    "parameter": "personalAccessToken",
    "value": "somevalue" 
}
```

`/api/scaledjob/create`

```json
{
    "name": "test-keda-scaledjob",
    "containers": [
        {
            "name": "scaledjobname",
            "image": "someimage:sometag",
            "imagePullPolicy": "IfNotPresent",
            "env": [
                {
                    "name": "someenv",
                    "value": "somevalue"
                }
            ]
        }
    ],
    "triggers": [
        {
            "type": "somekedatrigger",
            "metadata": {
                // Define your metadata here. This is just an example
                "ownerFromEnv": "REPO_OWNER",
                "runnerScope": "repo",
                "repoFromEnv": "REPO_NAME",
                "targetWorkflowQueueLength": "1"
            },
            "authenticationRef": {
                "name": "yoursecret"
            }
        }
    ]
}
```

`/api/scaledjob/delete/:scaledJobName`

```json
DELETE https://localhost:31323/api/scaledjob/delete/yourscaledjob
```

[TBD]