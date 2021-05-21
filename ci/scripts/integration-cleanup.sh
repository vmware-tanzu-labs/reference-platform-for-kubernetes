#!/bin/bash
# Copyright 2006-2021 VMware, Inc.
# SPDX-License-Identifier: MIT

# dashboard cleanup
# DONE - moved to make clean.role
kubectl delete tanzunamespace tanzu-dashboard
kubectl delete ns tanzu-dashboard --wait=false
kubectl delete clusterrole tanzu-dashboard-clusterrole
kubectl delete clusterrolebinding tanzu-dashboard-clusterrolebinding

# application pipeline cleanup
# DONE - moved to make clean.role
kubectl delete tanzunamespace tanzu-app-pipeline tanzu-app-pipeline-sit tanzu-app-pipeline-prod
kubectl delete ns tanzu-app-pipeline tanzu-app-pipeline-sit tanzu-app-pipeline-prod --wait=false

# container registry cleanup
# DONE - moved to make clean.role
kubectl delete tanzunamespace tanzu-container-registry
kubectl delete ns tanzu-container-registry --wait=false

# networking cleanup
# DONE - moved to make clean.role
kubectl delete globalnetworkpolicies.crd.projectcalico.org tanzu-global-deny-all
kubectl delete tanzunamespace tanzu-networking
kubectl delete ns tanzu-networking --wait=false

# admission control cleanup
# DONE - moved to make clean.role
kubectl delete crds configs.config.gatekeeper.sh constraintpodstatuses.status.gatekeeper.sh constrainttemplatepodstatuses.status.gatekeeper.sh constrainttemplates.templates.gatekeeper.sh
kubectl delete validatingwebhookconfigurations.admissionregistration.k8s.io gatekeeper-validating-webhook-configuration
kubectl delete mutatingwebhookconfigurations.admissionregistration.k8s.io mutating-webhook-configuration
kubectl delete tanzunamespace tanzu-admission
kubectl delete ns tanzu-admission --wait=false

# storage cleanup
# DONE - moved to make clean.role
kubectl delete sc bronze gold platinum silver ephemeral
kubectl delete tanzunamespace tanzu-storage
kubectl delete ns tanzu-storage --wait=false
kubectl delete clusterrolebinding csi-attacher-role csi-provisioner-role csi-resizer-role csi-snapshotter-role
kubectl delete clusterrole external-attacher-runner external-provisioner-runner external-resizer-runner external-snapshotter-runner
kubectl delete crd volumesnapshotclasses.snapshot.storage.k8s.io volumesnapshotcontents.snapshot.storage.k8s.io volumesnapshots.snapshot.storage.k8s.io
kubectl delete apiservice v1beta1.snapshot.storage.k8s.io
kubectl delete csidriver hostpath.csi.k8s.io

# ingress cleanup
# DONE - moved to make clean.role
kubectl delete tanzunamespace tanzu-ingress
kubectl delete ns tanzu-ingress --wait=false
kubectl delete mutatingwebhookconfigurations.admissionregistration.k8s.io  cert-manager-webhook
kubectl delete crd tlscertificatedelegations.contour.heptio.com tlscertificatedelegations.projectcontour.io httpproxies.projectcontour.io ingressroutes.contour.heptio.com

# security cleanup
# NOTE: only cleanup worker cluster cert-manager.  although we shouldn't be doing this in mgmt clusters
#       things happen.
if [[ $1 == 'worker' ]]; then
  kubectl delete apiservice v1beta1.webhook.cert-manager.io v1alpha2.acme.cert-manager.io
  kubectl delete crd certificaterequests.cert-manager.io certificates.cert-manager.io challenges.acme.cert-manager.io clusterissuers.cert-manager.io issuers.cert-manager.io orders.acme.cert-manager.io
  kubectl delete crd challenges.acme.cert-manager.io
  kubectl delete validatingwebhookconfigurations.admissionregistration.k8s.io cert-manager-webhook
  kubectl delete clusterrole cert-manager-certificates cert-manager-challenges cert-manager-clusterissuers cert-manager-ingress-shim cert-manager-issuers cert-manager-orders cert-manager-webhook:webhook-requester
  kubectl delete clusterrolebinding cert-manager-certificates cert-manager-challenges cert-manager-clusterissuers cert-manager-ingress-shim cert-manager-issuers cert-manager-orders cert-manager-webhook:auth-delegator
  kubectl delete rolebinding cert-manager-webhook:webhook-authentication-reader -n kube-system
fi
kubectl delete clusterrole tanzu-ca-jobs
kubectl delete clusterrolebinding tanzu-ca-jobs
kubectl delete psp tanzu-ca-jobs
kubectl delete tanzunamespace tanzu-security
kubectl delete ns tanzu-security --wait=false
kubectl delete clusterissuers letsencrypt-prod letsencrypt-stage

# monitoring
# DONE - moved to make clean.role
kubectl delete crd alertmanagers.monitoring.coreos.com podmonitors.monitoring.coreos.com prometheuses.monitoring.coreos.com prometheusrules.monitoring.coreos.com servicemonitors.monitoring.coreos.com thanosrulers.monitoring.coreos.com
kubectl delete tanzunamespace tanzu-monitoring
kubectl delete ns tanzu-monitoring --wait=false
kubectl delete apiservice v1beta1.metrics.k8s.io
kubectl delete clusterrolebinding grafana-clusterrolebinding grafana-rolebinding
kubectl delete clusterrole grafana-clusterrole grafana-role

# logging
# DONE - moved to make clean.role
kubectl delete tanzunamespace tanzu-logging
kubectl delete ns tanzu-logging --wait=false
kubectl delete crd apmservers.apm.k8s.elastic.co kibanas.kibana.k8s.elastic.co elasticsearches.elasticsearch.k8s.elastic.co
kubectl delete apiservice v1beta1.apm.k8s.elastic.co v1beta1.elasticsearch.k8s.elastic.co v1beta1.kibana.k8s.elastic.co
kubectl delete apiservice v1.apm.k8s.elastic.co v1.elasticsearch.k8s.elastic.co v1.kibana.k8s.elastic.co

# application-stack cleanup
# DONE - moved to make clean.role
kubectl delete tanzunamespace tanzu-app-stack-department tanzu-app-stack-employee tanzu-app-stack-gateway tanzu-app-stack-mongodb tanzu-app-stack-monitoring tanzu-app-stack-organization
kubectl delete ns tanzu-app-stack-department tanzu-app-stack-employee tanzu-app-stack-gateway tanzu-app-stack-mongodb tanzu-app-stack-monitoring tanzu-app-stack-organization --wait=false

# application-pipeline cleanup
kubectl delete tanzunamespace tanzu-app-pipeline tanzu-app-pipeline-sit tanzu-app-pipeline-prod
kubectl delete ns tanzu-app-pipeline tanzu-app-pipeline-sit tanzu-app-pipeline-prod --wait=false

# secrets cleanup
# DONE - moved to make clean.role
kubectl delete tanzunamespace tanzu-secrets
kubectl delete ns tanzu-secrets --wait=false
kubectl delete clusterrolebinding hashicorp-vault-agent-injector-binding hashicorp-vault-server-binding
kubectl delete clusterrole hashicorp-vault-agent-injector-clusterrole
kubectl delete mutatingwebhookconfigurations.admissionregistration.k8s.io hashicorp-vault-agent-injector-cfg

# identity cleanup
# DONE - moved to make clean.role
kubectl delete tanzunamespace tanzu-identity
kubectl delete ns tanzu-identity --wait=false
kubectl delete clusterrolebinding janedoe@example.com johndoe@example.com

# workload-tenancy cleanup
# DONE - moved to make clean.role
kubectl delete ns tanzu-workload-tenancy --wait=false
kubectl delete limitrange tanzu-default-limit-range
kubectl delete resourcequota tanzu-default-resource-quota
kubectl delete crd tanzunamespaces.tenancy.platform.cnr.vmware.com
kubectl delete clusterrolebinding namespace-operator-rolebinding
kubectl delete clusterrole namespace-operator-role

# mesh cleanup
# DONE - moved to make clean.role
kubectl delete validatingwebhookconfigurations.admissionregistration.k8s.io istiod-tanzu-mesh istiod-tanzu-mesh-system
kubectl delete mutatingwebhookconfigurations.admissionregistration.k8s.io istio-sidecar-injector-1-7-0-tanzu-mesh
kubectl delete deploy istio-operator
kubectl patch istiooperator/istiocontrolplane -n tanzu-mesh -p '{"metadata":{"finalizers":[]}}' --type=merge
kubectl delete istiooperator istiocontrolplane -n tanzu-mesh --wait=false
kubectl patch istiooperator/istiocontrolplane -n tanzu-mesh -p '{"metadata":{"finalizers":[]}}' --type=merge
kubectl delete istiooperator istiocontrolplane -n tanzu-mesh --wait=false
kubectl delete tanzunamespaces.tenancy.platform.cnr.vmware.com tanzu-mesh tanzu-mesh-system
kubectl delete ns tanzu-mesh tanzu-mesh-system --wait=false
kubectl delete apiservices.apiregistration.k8s.io v1alpha1.install.istio.io v1alpha2.config.istio.io  v1alpha3.networking.istio.io v1beta1.networking.istio.io  v1beta1.security.istio.io
kubectl delete crd adapters.config.istio.io attributemanifests.config.istio.io authorizationpolicies.security.istio.io destinationrules.networking.istio.io envoyfilters.networking.istio.io gateways.networking.istio.io handlers.config.istio.io httpapispecbindings.config.istio.io httpapispecs.config.istio.io instances.config.istio.io peerauthentications.security.istio.io quotaspecbindings.config.istio.io quotaspecs.config.istio.io requestauthentications.security.istio.io rules.config.istio.io serviceentries.networking.istio.io sidecars.networking.istio.io templates.config.istio.io virtualservices.networking.istio.io workloadentries.networking.istio.io

# loop until no namespaces are terminating
while kubectl get ns | grep Terminating; do
  echo "namespaces still terminating; retrying..."
  sleep 5
done
