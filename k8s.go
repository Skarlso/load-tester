package main

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// K8sClient wraps kubernetes dynamic client
type K8sClient struct {
	client dynamic.Interface
}

// NewK8sClient creates a new kubernetes client
func NewK8sClient() (*K8sClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			return nil, fmt.Errorf("failed to create k8s config: %w", err)
		}
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s client: %w", err)
	}

	return &K8sClient{client: client}, nil
}

// CreateObject creates a kubernetes object
func (k *K8sClient) CreateObject(ctx context.Context, obj *unstructured.Unstructured, config *Config) error {
	gvk := obj.GroupVersionKind()
	gvr := schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: config.Resource,
	}

	namespace := obj.GetNamespace()
	if namespace == "" {
		namespace = "default"
	}

	_, err := k.client.Resource(gvr).Namespace(namespace).Create(ctx, obj, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create object %s/%s: %w", obj.GetKind(), obj.GetName(), err)
	}

	return nil
}

// RunLoadTest executes the load test
func (k *K8sClient) RunLoadTest(ctx context.Context, config *Config) error {
	fmt.Printf("Starting load test: creating %d objects with %ds delay\n", config.LoadTest.Count, config.LoadTest.Delay)

	for i := 0; i < config.LoadTest.Count; i++ {
		obj, err := config.GenerateObject()
		if err != nil {
			return fmt.Errorf("failed to generate object %d: %w", i+1, err)
		}

		err = k.CreateObject(ctx, obj, config)
		if err != nil {
			return fmt.Errorf("failed to create object %d: %w", i+1, err)
		}

		fmt.Printf("Created object %d: %s/%s\n", i+1, obj.GetKind(), obj.GetName())

		if i < config.LoadTest.Count-1 && config.LoadTest.Delay > 0 {
			time.Sleep(time.Duration(config.LoadTest.Delay) * time.Second)
		}
	}

	fmt.Printf("Load test completed: %d objects created\n", config.LoadTest.Count)
	return nil
}
