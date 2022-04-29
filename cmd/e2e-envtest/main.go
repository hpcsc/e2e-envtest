package main

import (
	"context"
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
	"strings"
)

func main() {
	createdBy := flag.String("created-by", "", "value of created-by label to filter")
	flag.Parse()

	if *createdBy == "" {
		panic("created-by flag is required")
	}

	client, err := newKubeClient()
	if err != nil {
		panic(err)
	}

	namespaces, err := listNamespacesByLabel(client, "created-by", *createdBy)
	if err != nil {
		panic(err)
	}

	fmt.Printf("found %d namespaces: %v", len(namespaces), strings.Join(namespaces, ", "))
}

func listNamespacesByLabel(client dynamic.Interface, labelKey string, labelValue string) ([]string, error) {
	result, err := client.Resource(schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "namespaces",
	}).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return []string{}, nil
		}

		return nil, fmt.Errorf("failed to list namespaces: %v", err)
	}

	var namespaces []string
	for _, ns := range result.Items {
		if l, ok := ns.GetLabels()[labelKey]; ok && l == labelValue {
			namespaces = append(namespaces, ns.GetName())
		}
	}
	return namespaces, nil
}

func newKubeClient() (dynamic.Interface, error) {
	kubeConfig, err := ctrl.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load kubernetes config: %v", err)
	}

	client, err := dynamic.NewForConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client from config: %v", err)
	}

	return client, nil
}
