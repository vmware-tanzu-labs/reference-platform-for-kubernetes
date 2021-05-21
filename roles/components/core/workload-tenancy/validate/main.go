// Copyright 2006-2021 VMware, Inc.
// # SPDX-License-Identifier: MIT
package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	resultsFilename   = "sonobuoy_results.yaml"
	logFilename       = "plugin.log"
	nsName            = "tanzu-workload-tenancy-validate"
	resourceQuotaName = "tanzu-resource-quota"
	limitRangeName    = "tanzu-limit-range"
	statusPassed      = "passed"
	statusFailed      = "failed"
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
	logger := log.New(mw, "workload-tenancy: ", log.LstdFlags)

	logger.Println("workload-tenancy module validation starting...")

	resultsFilepath := fmt.Sprintf("%s/%s", resultsDir, resultsFilename)

	// Sonobuoy will provide a service account for in-cluster access
	config, err := rest.InClusterConfig()
	if err != nil {
		logger.Fatalf("Error creating in cluster config: %v\n", err)
	}
	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Fatalf("Error creating kubernetes client: %v\n", err)
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		logger.Fatalf("Error creating dynamic client: %v\n", err)
	}

	// Each test item will be appended to this array
	testItems := []Item{}

	// Test creation of TanzuNamespace resource
	tanzuNSResult, err := testTanzuNsCreate(dynamicClient)
	testItems = append(testItems, tanzuNSResult)
	if err != nil {
		// TanzuNamespace creation failed - no point in running other tests
		// record result, signal done and exit out
		overallStatus = statusFailed
		logger.Printf("Error creating tanzu namespace: %v\n", err)
		if err = recordResults(testItems, overallStatus, resultsFilepath); err != nil {
			logger.Printf("Error recording test results: %v\n", err)
		}
		if err = signalSonobuoy(resultsDir); err != nil {
			logger.Fatalf("Error signalling test completion: %v\n", err)
		}
		os.Exit(1)
	}
	logger.Printf("Tanzu namespace %s successfully created", nsName)

	// Pause for resource creation
	time.Sleep(time.Second * 10)

	// Check that namespace has been created
	checkNSResult, err := checkNSCreated(k8sClient)
	testItems = append(testItems, checkNSResult)
	if err != nil {
		overallStatus = statusFailed
		logger.Printf("Error gettting namespaces: %v\n", err)
	} else {
		logger.Printf("Kubernetes namespace %s found", nsName)
	}

	// Check that resource quota has been created
	checkRQResult, err := checkRQCreated(k8sClient)
	testItems = append(testItems, checkRQResult)
	if err != nil {
		overallStatus = statusFailed
		logger.Printf("Error gettting resource quota: %v\n", err)
	} else {
		logger.Printf("Resource quota %s found", resourceQuotaName)
	}

	// Check that limit range has been created
	checkLRResult, err := checkLRCreated(k8sClient)
	testItems = append(testItems, checkLRResult)
	if err != nil {
		overallStatus = statusFailed
		logger.Printf("Error gettting limit range: %v\n", err)
	} else {
		logger.Printf("Limit range %s found", limitRangeName)
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
	if err = clean(dynamicClient); err != nil {
		logger.Fatalf("Error cleaning up tanzu namespace: %v\n", err)
	}

	logger.Println("workload-tenancy validation complete")

	// Signal to sonbobuoy worker that test is complete
	if err = signalSonobuoy(resultsDir); err != nil {
		logger.Fatalf("Error signalling test completion: %v\n", err)
	}
}

// testTanzuNsCreate tests the creation of a TanzuNamespace resource
func testTanzuNsCreate(dynamicClient dynamic.Interface) (Item, error) {

	tanzuNSResult := Item{
		Name: "[workload-tenancy] Creating a TanzuNamespace resource should succeed.",
	}

	tanzuNSRes := schema.GroupVersionResource{
		Group:    "tenancy.platform.cnr.vmware.com",
		Version:  "v1alpha1",
		Resource: "tanzunamespaces",
	}

	tanzuNS := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "tenancy.platform.cnr.vmware.com/v1alpha1",
			"kind":       "TanzuNamespace",
			"metadata": map[string]interface{}{
				"name": nsName,
			},
			"spec": map[string]interface{}{
				"tanzuNamespaceName": nsName,
			},
		},
	}

	_, err := dynamicClient.Resource(tanzuNSRes).Create(context.TODO(), tanzuNS, metav1.CreateOptions{})
	if err != nil {
		tanzuNSResult.Status = statusFailed
		return tanzuNSResult, err
	}
	tanzuNSResult.Status = statusPassed

	return tanzuNSResult, nil
}

// checkNSCreated checks to make sure resource quota was created
func checkNSCreated(k8sClient *kubernetes.Clientset) (Item, error) {

	checkNSResult := Item{
		Name: "[workload-tenancy] Creating a TanzuNamespace resource should create a Kubernetes namespace.",
	}

	_, err := k8sClient.CoreV1().Namespaces().Get(context.TODO(), nsName, metav1.GetOptions{})
	if err != nil {
		checkNSResult.Status = statusFailed
		return checkNSResult, err
	} else {
		checkNSResult.Status = statusPassed
	}

	return checkNSResult, nil
}

// checkRQCreated checks to make sure resource quota was created
func checkRQCreated(k8sClient *kubernetes.Clientset) (Item, error) {

	checkRQResult := Item{
		Name: "[workload-tenancy] Creating a TanzuNamespace resource should create a resource quota.",
	}

	_, err := k8sClient.CoreV1().ResourceQuotas(nsName).Get(context.TODO(), resourceQuotaName, metav1.GetOptions{})
	if err != nil {
		checkRQResult.Status = statusFailed
		return checkRQResult, err
	} else {
		checkRQResult.Status = statusPassed
	}

	return checkRQResult, nil
}

// checkLRCreated checks to make sure resource quota was created
func checkLRCreated(k8sClient *kubernetes.Clientset) (Item, error) {

	checkLRResult := Item{
		Name: "[workload-tenancy] Creating a TanzuNamespace resource should create a limit range.",
	}

	_, err := k8sClient.CoreV1().LimitRanges(nsName).Get(context.TODO(), limitRangeName, metav1.GetOptions{})
	if err != nil {
		checkLRResult.Status = statusFailed
		return checkLRResult, err
	} else {
		checkLRResult.Status = statusPassed
	}

	return checkLRResult, nil
}

// clean removes the test TanzuNamespace
func clean(dynamicClient dynamic.Interface) error {

	tanzuNSRes := schema.GroupVersionResource{
		Group:    "tenancy.platform.cnr.vmware.com",
		Version:  "v1alpha1",
		Resource: "tanzunamespaces",
	}

	if err := dynamicClient.Resource(tanzuNSRes).Delete(context.TODO(), nsName, metav1.DeleteOptions{}); err != nil {
		return err
	}

	return nil
}

// recordResults writes the results file for sonobuoy worker to collect
func recordResults(testItems []Item, overallStatus, resultsFilepath string) error {

	validationResult := Item{
		Name:   "RPK Workload Tenancy Validation",
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
