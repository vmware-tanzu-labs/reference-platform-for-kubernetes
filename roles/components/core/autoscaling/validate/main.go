// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT
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
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	vpav1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1"
	vpa_clientset "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	resultsFilename  = "sonobuoy_results.yaml"
	logFilename      = "plugin.log"
	nsName           = "tanzu-autoscaling-validate"
	vpaName          = "tanzu-vpa"
	testResourceName = "validation-test"
	statusPassed     = "passed"
	statusFailed     = "failed"
	nginxImage       = "nginx:1.19"
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
	logger := log.New(mw, "autoscaling: ", log.LstdFlags)

	logger.Println("autoscaling module validation starting...")

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
	vpaClient, err := vpa_clientset.NewForConfig(config)
	if err != nil {
		logger.Fatalf("Error creating VPA client: %v\n", err)
	}

	// Each test item will be appended to this array
	testItems := []Item{}

	// Create pod and service resources to test vpa
	workloadResult, err := testWorkloadCreate(k8sClient)
	testItems = append(testItems, workloadResult)
	if err != nil {
		// If there's an error creating a test workload - no nded to run other
		// tests.  Record results, signal done and exit.
		overallStatus = statusFailed
		logger.Printf("Error creating autoscaling test workload: %v\n", err)
		if err = recordResults(testItems, overallStatus, resultsFilepath); err != nil {
			logger.Printf("Error recording test results: %v\n", err)
		}
		if err = signalSonobuoy(resultsDir); err != nil {
			logger.Fatalf("Error signalling test completion: %v\n", err)
		}
		os.Exit(1)
	}
	logger.Println("Autoscaling test workload successfully created")

	// Create vertical pod autoscaler for test workload
	vpaResult, err := testVpaCreate(vpaClient)
	testItems = append(testItems, vpaResult)
	if err != nil {
		overallStatus = statusFailed
		logger.Printf("Error creating vertical pod autoscaler: %v\n", err)
	} else {
		logger.Printf("Vertical pod autoscaler %s created", vpaName)
	}

	// Pause for a few seconds to allow VPA to generate a recommendation
	time.Sleep(time.Second * 10)

	// Check VPA for a recommendation
	vpaRecommendation, err := checkVpaRecommendation(vpaClient)
	testItems = append(testItems, vpaRecommendation)
	if err != nil {
		overallStatus = statusFailed
		logger.Printf("Vertical pod autoscaler %s contains no resource recommendations: %v\n", vpaName, err)
	} else {
		logger.Printf("Recommendation found in vertical pod autoscaler %s", vpaName)
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
		logger.Fatalf("Error cleaning up tanzu namespace: %v\n", err)
	}

	logger.Println("workload-tenancy validation complete")

	// Signal to sonbobuoy worker that test is complete
	if err = signalSonobuoy(resultsDir); err != nil {
		logger.Fatalf("Error signalling test completion: %v\n", err)
	}
}

// testWorkloadCreate creates a pod and service resources to test vertical pod
// autoscaler functionality
func testWorkloadCreate(k8sClient *kubernetes.Clientset) (Item, error) {

	workloadCreateResult := Item{
		Name: "[autoscaling] Creating pod and service resources should succeed.",
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

	workloadCreateResult.Status = statusPassed

	return workloadCreateResult, nil
}

// testVpaCreate tests the creation of a vertical pod autoscaler
func testVpaCreate(vpaClient *vpa_clientset.Clientset) (Item, error) {

	vpaCreateResult := Item{
		Name: "[autoscaling] Creating a vertical pod autoscaler resource should succeed.",
	}

	objectRef := autoscalingv1.CrossVersionObjectReference{
		Kind:       "Deployment",
		Name:       testResourceName,
		APIVersion: "apps/v1",
	}

	updateMode := vpav1.UpdateModeOff

	updatePolicy := vpav1.PodUpdatePolicy{
		UpdateMode: &updateMode,
	}

	vpa := &vpav1.VerticalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      vpaName,
			Namespace: nsName,
		},
		Spec: vpav1.VerticalPodAutoscalerSpec{
			TargetRef:    &objectRef,
			UpdatePolicy: &updatePolicy,
		},
	}

	_, err := vpaClient.AutoscalingV1().VerticalPodAutoscalers(nsName).Create(context.TODO(), vpa, metav1.CreateOptions{})
	if err != nil {
		vpaCreateResult.Status = statusFailed
		return vpaCreateResult, err
	}

	vpaCreateResult.Status = statusPassed

	return vpaCreateResult, nil
}

// checkVpaRecommendation checks that the VPA includes a container resource
// recommendation for the workload
func checkVpaRecommendation(vpaClient *vpa_clientset.Clientset) (Item, error) {

	vpaRecommendationResult := Item{
		Name: "[autoscaling] The vertical pod autoscaler should contain a resource recommendation.",
	}

	vpa, err := vpaClient.AutoscalingV1().VerticalPodAutoscalers(nsName).Get(context.TODO(), vpaName, metav1.GetOptions{})
	if err != nil {
		vpaRecommendationResult.Status = statusFailed
		return vpaRecommendationResult, err
	}

	//if len(vpa.Status.Recommendation.ContainerRecommendations) < 1 {
	//if vpa.Status.Recommendation != nil {
	emptyRecommendation := &vpav1.RecommendedPodResources{}
	if vpa.Status.Recommendation == emptyRecommendation {
		vpaRecommendationResult.Status = statusFailed
		return vpaRecommendationResult, err
	}

	vpaRecommendationResult.Status = statusPassed

	return vpaRecommendationResult, nil
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
		Name:   "RPK Autoscaling Validation",
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
