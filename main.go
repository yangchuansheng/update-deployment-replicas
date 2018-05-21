/*
Copyright 2017 The Kubernetes Authors.

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

package main

import (
        "os"
	"flag"
	"fmt"
	//apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/rest"
)

func main() {

	var ns string
	flag.StringVar(&ns, "namespace", "default", "namespace name")
	deploymentName := flag.String("deployment", "", "deployment name")
	imageName := flag.String("image", "", "new image name")
	replicasNum := flag.Int("replicas",  0, "number of replicas")
	flag.Parse()

	if *deploymentName == "" {
		fmt.Println("You must specify the deployment name.")
		os.Exit(0)
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	deploymentsClient := clientset.AppsV1beta1().Deployments(ns)

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {

		var updateErr error

		result, err := deploymentsClient.Get(*deploymentName, metav1.GetOptions{})
		if err != nil {
			panic(err.Error())
		}
		if errors.IsNotFound(err) {
			fmt.Printf("Deployment not found\n")
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting deployment%v\n", statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found deployment\n")
			name := result.GetName()
			fmt.Println("name ->", name)
			oldreplicas := *result.Spec.Replicas
			fmt.Println("old replicas ->", oldreplicas)
			fmt.Println("new replicas ->", *replicasNum)
			*result.Spec.Replicas = int32(*replicasNum)
			if *imageName != "" {
				result.Spec.Template.Spec.Containers[0].Image = *imageName
			}
			_, updateErr = deploymentsClient.Update(result)
			println(updateErr)
		}
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	fmt.Println("Updated deployment...")
}
