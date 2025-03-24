# keda-goclient

An application functioning as a REST API that does the following:
- Uses `client-go` to interface with Kubernetes. This has been developed in `wsl2` with a local Docker enabled Kubernetes single node cluster
- Acts as a client to create KEDA `ScaledJobs` and `ScaledObjects`
- `scaler-kubernetes-webhook` `.yamls` can be deployed to set up both Mutating and Admission Webhooks that will be invoked during KEDA resource creation. The logic currently is that a `metadata.label` of either `scaledjob` or `scaledobject` will be adding during resource creation by the Mutating webhook