module github.com/vmware-tanzu-labs/reference-platform-for-kubernetes/roles/components/core/autoscaling/validate

go 1.15

require (
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/autoscaler/vertical-pod-autoscaler v0.9.0
	k8s.io/client-go v0.19.2
)
