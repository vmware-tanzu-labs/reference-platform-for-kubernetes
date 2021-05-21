// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT
package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	resultsFilename     = "sonobuoy_results.yaml"
	logFilename         = "plugin.log"
	nsName              = "tanzu-ingress-validate"
	nginxImage          = "nginx:1.19"
	testResourceName    = "validation-test"
	ingressTestTries    = 15
	ingressTestInterval = time.Second * 30
	statusPassed        = "passed"
	statusFailed        = "failed"
)

var (
	resultsDir    string
	overallStatus string
)

type Item struct {
	Name   string `json:"name" yaml:"name"`
	Status string `json:"status" yaml:"status"`
	Items  []Item `json:"items,omitempty" yaml:"items,omitempty"`
}

func main() {

	flag.StringVar(&resultsDir, "results-dir", "/tmp/results", "Directory to save test results to")
	flag.Parse()

	// Log to both stdout and to a file to include in sonobuoy results
	logFile, err := os.OpenFile(fmt.Sprintf("%s/%s", resultsDir, logFilename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	mw := io.MultiWriter(os.Stdout, logFile)
	logger := log.New(mw, "ingress: ", log.LstdFlags)

	logger.Println("ingress module validation starting...")

	resultsFilepath := fmt.Sprintf("%s/%s", resultsDir, resultsFilename)

	// Sonobuoy will provide a service account for in-cluster access
	config, err := rest.InClusterConfig()
	if err != nil {
		logger.Fatalf("Error creating in cluster config: %v\n", err)
	}
	k8sClient, err := kubernetes.NewForConfig(config)
	_, err = kubernetes.NewForConfig(config)
	if err != nil {
		logger.Fatalf("Error creating kubernetes client: %v\n", err)
	}

	// Each test item will be appended to this array
	testItems := []Item{}

	// Get ingress domain
	ingressDomain, err := getIngressDomain(k8sClient)
	if err != nil {
		// If the ingress domain cannot be determined, tests cannot be run
		overallStatus = statusFailed
		logger.Printf("Error determining ingress domain: %v\n", err)
		if err = recordResults(testItems, overallStatus, resultsFilepath); err != nil {
			logger.Printf("Error recording test results: %v\n", err)
		}
		if err = signalSonobuoy(resultsDir); err != nil {
			logger.Fatalf("Error signalling test completion: %v\n", err)
		}
		os.Exit(1)
	}

	// Create pod, service, ingress resource to test ingress layer
	workloadResult, err := testWorkloadCreate(k8sClient, ingressDomain)
	testItems = append(testItems, workloadResult)
	if err != nil {
		// If there's an error creating a test workload - no nded to run other
		// tests.  Record results, signal done and exit.
		overallStatus = statusFailed
		logger.Printf("Error creating ingress test workload: %v\n", err)
		if err = recordResults(testItems, overallStatus, resultsFilepath); err != nil {
			logger.Printf("Error recording test results: %v\n", err)
		}
		if err = signalSonobuoy(resultsDir); err != nil {
			logger.Fatalf("Error signalling test completion: %v\n", err)
		}
		os.Exit(1)
	}
	logger.Println("Ingress test workload successfully created")

	// Test ingress functionality by calling service through ingress
	ingressTestResult, err := ingressTestRequest(ingressDomain)
	testItems = append(testItems, ingressTestResult)
	if err != nil {
		overallStatus = statusFailed
		logger.Printf("Error calling ingress test workload: %v\n", err)
	} else {
		logger.Printf("Successfully called test workload through ingress")
	}

	// Overall status is passed if not previously set to failed
	if overallStatus == "" {
		overallStatus = statusPassed
	}

	// Write test results to results directory
	if err = recordResults(testItems, overallStatus, resultsFilepath); err != nil {
		logger.Printf("Error recording test results: %v\n", err)
	}

	// Clean up test resources
	if err = clean(k8sClient); err != nil {
		logger.Fatalf("Error cleaning up ingress validation resoruces: %v\n", err)
	}

	logger.Println("ingress validation complete")

	// Signal to sonbobuoy worker that test is complete
	if err = signalSonobuoy(resultsDir); err != nil {
		logger.Fatalf("Error signalling test completion: %v\n", err)
	}
}

// testWorkloadCreate creates a pod, service and ingress resource to test
// ingress functionality
func testWorkloadCreate(k8sClient *kubernetes.Clientset, ingressDomain string) (Item, error) {

	workloadCreateResult := Item{
		Name: "[ingress] Creating a pod, service and ingress resource should succeed.",
	}

	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: nsName,
		},
	}

	_, err := k8sClient.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	if err != nil {
		workloadCreateResult.Status = statusFailed
		return workloadCreateResult, err
	}

	po := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testResourceName,
			Namespace: nsName,
			Labels: map[string]string{
				"app": testResourceName,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx",
					Image: nginxImage,
					Ports: []corev1.ContainerPort{
						{
							Name:          "http",
							Protocol:      corev1.ProtocolTCP,
							ContainerPort: 80,
						},
					},
				},
			},
		},
	}

	_, err = k8sClient.CoreV1().Pods(nsName).Create(context.TODO(), po, metav1.CreateOptions{})
	if err != nil {
		workloadCreateResult.Status = statusFailed
		return workloadCreateResult, err
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testResourceName,
			Namespace: nsName,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": testResourceName,
			},
			Ports: []corev1.ServicePort{
				{
					Protocol: corev1.ProtocolTCP,
					Port:     80,
				},
			},
		},
	}

	_, err = k8sClient.CoreV1().Services(nsName).Create(context.TODO(), svc, metav1.CreateOptions{})

	ing := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testResourceName,
			Namespace: nsName,
			Annotations: map[string]string{
				"external-dns.alpha.kubernetes.io/target": fmt.Sprintf("ingress.%s", ingressDomain),
			},
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: fmt.Sprintf("%s.%s", testResourceName, ingressDomain),
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path: "/",
									Backend: networkingv1.IngressBackend{
										ServiceName: testResourceName,
										ServicePort: intstr.IntOrString{IntVal: 80},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	_, err = k8sClient.NetworkingV1beta1().Ingresses(nsName).Create(context.TODO(), ing, metav1.CreateOptions{})

	workloadCreateResult.Status = statusPassed

	return workloadCreateResult, nil
}

// ingressTestRequest makes an HTTP request to the the test workload through the
// ingress to validate ingress functionality
func ingressTestRequest(ingressDomain string) (Item, error) {

	ingressTestResult := Item{
		Name: "[ingress] Calling test workload through ingress should return a 200 response.",
	}
	ingressAttempts := 0
	var responseErr error

	for ingressAttempts < ingressTestTries {

		_, responseErr = http.Get(fmt.Sprintf("http://%s.%s", testResourceName, ingressDomain))
		if responseErr != nil {
			ingressAttempts += 1
			time.Sleep(ingressTestInterval)
		} else {
			ingressTestResult.Status = statusPassed
			return ingressTestResult, nil
		}
	}

	ingressTestResult.Status = statusFailed
	return ingressTestResult, responseErr
}

// getIngressDomain retrieves the ingress domain from the ingress envoy service
func getIngressDomain(k8sClient *kubernetes.Clientset) (string, error) {

	envoySvc, err := k8sClient.CoreV1().Services("tanzu-ingress").Get(context.TODO(), "envoy", metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	for k, v := range envoySvc.ObjectMeta.Annotations {
		if k == "external-dns.alpha.kubernetes.io/hostname" {
			return strings.Replace(v, "ingress.", "", 1), nil
		}
	}

	return "", errors.New("Unable to determine ingress domain from envoy service annotations")
}

// clean removes the namespace used for test resources
func clean(k8sClient *kubernetes.Clientset) error {

	if err := k8sClient.CoreV1().Namespaces().Delete(context.TODO(), nsName, metav1.DeleteOptions{}); err != nil {
		return err
	}

	return nil
}

// recordResults writes the results file for sonobuoy worker to collect
func recordResults(testItems []Item, overallStatus, resultsFilepath string) error {

	validationResult := Item{
		Name:   "RPK Ingress Validation",
		Status: overallStatus,
		Items: []Item{
			{
				Name:   "Validation Tests",
				Status: overallStatus,
				Items:  testItems,
			},
		},
	}

	resultsYaml, err := yaml.Marshal(validationResult)

	if _, err = os.OpenFile(resultsFilepath, os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		return err
	}
	if err = ioutil.WriteFile(resultsFilepath, resultsYaml, 0644); err != nil {
		return err
	}

	return nil
}

// signalSonobuoy creates results tarball and signals to sonbuoy worker that
// test is complete by creating a `done` file with path to results tarball
func signalSonobuoy(resultsDir string) error {

	sonobuoyResultsFile, err := os.Create(fmt.Sprintf("%s/results.tar.gz", resultsDir))
	if err != nil {
		return err
	}
	if err = archive(fmt.Sprintf("%s/", resultsDir), sonobuoyResultsFile); err != nil {
		return err
	}
	if _, err = os.OpenFile(fmt.Sprintf("%s/done", resultsDir), os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		return err
	}
	if err = ioutil.WriteFile(fmt.Sprintf("%s/done", resultsDir), []byte(fmt.Sprintf("%s/results.tar.gz", resultsDir)), 0644); err != nil {
		return err
	}

	return nil
}

// archive creates a tarball for files in a directory
func archive(src string, writers ...io.Writer) error {

	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("Unable to tar files - %v", err.Error())
	}

	mw := io.MultiWriter(writers...)

	gzw := gzip.NewWriter(mw)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	return filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(strings.Replace(file, src, "", -1), string(filepath.Separator))
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}

		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		f.Close()

		return nil
	})
}
