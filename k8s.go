/*
Copyright 2016 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

func testK8s() {
	_ = k8sClient()
}

func k8sClient() *kubernetes.Clientset {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	var err1 error
	if err != nil {
		kubeconfig := "kubeconfig"
		config, err1 = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	if err1 != nil {
		log.Println("read cluster config:", err.Error())
		log.Fatal(err1)
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset

}

func gethpa(name string, namespace string) ([]byte, error) {
	clientset := k8sClient()
	hpaobj, err := clientset.AutoscalingV1().HorizontalPodAutoscalers(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("HPA %s in namespace %s not found\n", name, namespace)
		return nil, err
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting hpa %s in namespace %s: %v\n",
			name, namespace, statusError.ErrStatus.Message)
		return nil, err
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found HPA %s in namespace %s\n", name, namespace)
		annotations := hpaobj.GetAnnotations()
		if event, ok := annotations["hpa-event"]; ok {
			return []byte(event), nil
		}
		return nil, nil

	}

}

func getsecret(secretloc string) (map[string][]byte, error) {
	sec := strings.Split(secretloc, "/")
	clientset := k8sClient()
	secret, err := clientset.CoreV1().Secrets(sec[0]).Get(context.TODO(), sec[1], metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Secret %s in namespace %s not found\n", sec[1], sec[0])
		return nil, err
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting hpa %s in namespace %s: %v\n",
			sec[1], sec[0], statusError.ErrStatus.Message)
		return nil, err
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found Secret %s in namespace %s\n", sec[1], sec[0])
		return secret.Data, nil
	}

}

func checkhpa(name string, namespace string) (secret string, Error error) {
	clientset := k8sClient()
	hpaobj, err := clientset.AutoscalingV1().HorizontalPodAutoscalers(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("HPA %s in namespace %s not found\n", name, namespace)
		return "", err
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting hpa %s in namespace %s: %v\n",
			name, namespace, statusError.ErrStatus.Message)
		return "", err
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found HPA %s in namespace %s\n", name, namespace)
		annotations := hpaobj.GetAnnotations()
		if _, ok := annotations["hpa-event"]; ok {
			ret := annotations["hpa-event"]
			return ret, nil
		}
		return "", nil

	}
}

// SecretData contains the data to send to the webhook for that particular HPA
type SecretData struct {
	Body    map[string]string `json:"body"`
	Headers map[string]string `json:"headers"`
	URL     string            `json:"url"`
}
