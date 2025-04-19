# keda-goclient

An application functioning as a REST API that does the following:
- Uses `client-go` to interface with Kubernetes. This has been developed in `wsl2` with a local Docker enabled Kubernetes single node cluster
- Acts as a client to create KEDA `ScaledJobs` and `ScaledObjects`. In this repo, it was specifically for creating GitHub Action runners.
- `scaler-kubernetes-webhook` `.yamls` can be deployed to set up both Mutating and Admission Webhooks that will be invoked during KEDA resource creation. The logic currently is that a `metadata.label` of either `scaledjob` or `scaledobject` will be adding during resource creation by the Mutating webhook


----------
Below is the request body structure that some of these endpoints accept. `POST` bodies should follow the same general structure that a Kubernetes or KEDA `.yaml` expects

`POST /api/secret/create`

```json
{
    "name": "yoursecret",
    "parameter": "personalAccessToken",
    "value": "yourPatValue" 
}
```

`POST /api/scaledjob/create`

```json
{
    "name": "test-keda-scaledjob",
    "containers": [
        {
            "name": "github-runner",
            "image": "self-hosted-github-action-runner:local",
            "imagePullPolicy": "IfNotPresent",
            "env": [
                {
                    "name": "GITHUB_PAT",
                    "value": "github_pat_xxx"
                },
                {
                    "name": "REPO_OWNER",
                    "value": "YourUser"
                },
                {
                    "name": "REPO_NAME",
                    "value": "yourrepo"
                },
                {
                    "name": "REPO_URL",
                    "value": "https://github.com/YourUser/yourrepo"
                },
                {
                    "name": "REGISTRATION_TOKEN_API_URL",
                    "value": "https://api.github.com/repos/YourUser/yourrepo/actions/runners/registration-token"
                }
            ]
        }
    ],
    "triggers": [
        {
            "type": "github-runner",
            "metadata": {
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

`POST /api/scaledobject/create`

```json
{
    "name": "yourdeployment",
    "triggers": [
        {
            "type": "github-runner",
            "metadata": {
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

`POST /api/deployment/create`

```json
{
    "name": "githubrunner-deployment",
    "replicas": 1,
    "containers": [
        {
            "name": "github-runner",
            "image": "self-hosted-github-action-runner:local",
            "imagePullPolicy": "IfNotPresent",
            "env": [
                {
                    "name": "GITHUB_PAT",
                    "value": "github_pat_xxx"
                },
                {
                    "name": "REPO_OWNER",
                    "value": "YourUser"
                },
                {
                    "name": "REPO_NAME",
                    "value": "self-hosted-github-action-runner"
                },
                {
                    "name": "REPO_URL",
                    "value": "https://github.com/YourUser/yourrepo"
                },
                {
                    "name": "REGISTRATION_TOKEN_API_URL",
                    "value": "https://api.github.com/repos/YourUser/yourrepo/actions/runners/registration-token"
                }
            ]
        }
    ]
}
```

`DELETE /api/scaledjob/delete/:scaledJobName`

```json
DELETE https://localhost:31323/api/scaledjob/delete/yourscaledjob
```

`DELETE /api/scaledjob/delete/:scaledJobName`

```json
DELETE https://localhost:31323/api/scaledobject/delete/yourscaledobject
```

