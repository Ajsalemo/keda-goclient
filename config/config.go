package config

import (
	"flag"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

var (
	kedaGVRJobs = schema.GroupVersionResource{
		Group:    "keda.sh",
		Version:  "v1alpha1",
		Resource: "scaledjobs",
	}
)

var (
	kedaGVRObject = schema.GroupVersionResource{
		Group:    "keda.sh",
		Version:  "v1alpha1",
		Resource: "scaledobjects",
	}
)

func kedaGVR(clientType string) schema.GroupVersionResource {
	if clientType == "scaledObject" {
		return kedaGVRObject
	} else {
		return kedaGVRJobs
	}
}

func KubeConfig(clientType string) dynamic.NamespaceableResourceInterface {
	var kubeconfig *string
	// See: https://maxchadwick.xyz/blog/preventing-flag-conflicts-in-go
	// NewFlagSet is added to prevent `panic: flag redefined: kubeconfig` - which will happen when we're importing this function into controllers and defining the client type there (eg. scaledjob or scaledobject)
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return dynamicClient.Resource(kedaGVR(clientType))
}
