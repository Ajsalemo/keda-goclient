package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
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

func kedaGVR() schema.GroupVersionResource {
	// Mock function. Testing to see if we can switch between object types / different clients
	scaledObject := false
	if scaledObject {
		return kedaGVRObject
	} else {
		return kedaGVRJobs
	}
}

func GetScaledObjectByName(dynamicClient *dynamic.DynamicClient, name string) {
	scaledObjectClient := dynamicClient.Resource(kedaGVR())
	scaledObject, err := scaledObjectClient.Namespace("default").Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("Error retrieving scaledjobs: %s", err)
	} else {
		fmt.Printf("Got ScaledObject: %v", scaledObject)
	}
}

func main() {
	var kubeconfig *string
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

	GetScaledObjectByName(dynamicClient, "github-runner")
}
