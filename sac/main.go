package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Printf("Error reading kubeconfig: %v\n", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("--- STARTING CLUSTER SECURITY AUDIT ---")

	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	vulnerabilitiesFound := 0

	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			// Check 1: Running as ROOT (UID 0)
			isRoot := true
			if container.SecurityContext != nil && container.SecurityContext.RunAsNonRoot != nil {
				isRoot = !*container.SecurityContext.RunAsNonRoot
			}

			// Check 2: Privileged mode (grants host access)
			isPrivileged := false
			if container.SecurityContext != nil && container.SecurityContext.Privileged != nil {
				isPrivileged = *container.SecurityContext.Privileged
			}

			if isRoot || isPrivileged {
				vulnerabilitiesFound++
				fmt.Printf("\n[!] VULNERABILITY in pod: %s (Namespace: %s)\n", pod.Name, pod.Namespace)
				fmt.Printf("    Container: %s\n", container.Name)
				if isRoot {
					fmt.Println("    - Running as ROOT (Potential privilege escalation)")
				}
				if isPrivileged {
					fmt.Println("    - PRIVILEGED MODE (Container can escape to host!)")
				}

			}
			fmt.Println("    - NOT VULNERABLE: running as non-root and not privileged")

		}
	}

	fmt.Printf("\n--- AUDIT COMPLETE. Vulnerabilities found: %d ---\n", vulnerabilitiesFound)
}
