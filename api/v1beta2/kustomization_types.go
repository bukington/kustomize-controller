/*
Copyright 2021 The Flux authors

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

package v1beta2

import (
	"time"

	"github.com/fluxcd/pkg/apis/kustomize"
	"github.com/fluxcd/pkg/apis/meta"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	KustomizationKind         = "Kustomization"
	KustomizationFinalizer    = "finalizers.fluxcd.io"
	MaxConditionMessageLength = 20000
	EnabledValue              = "enabled"
	DisabledValue             = "disabled"
	MergeValue                = "merge"
)

// KustomizationSpec defines the configuration to calculate the desired state from a Source using Kustomize.
type KustomizationSpec struct {
	// CommonMetadata specifies the common labels and annotations that are applied to all resources.
	// Any existing label or annotation will be overridden if its key matches a common one.
	// +optional
	CommonMetadata *CommonMetadata `json:"commonMetadata,omitempty"`

	// DependsOn may contain a meta.NamespacedObjectReference slice
	// with references to Kustomization resources that must be ready before this
	// Kustomization can be reconciled.
	// +optional
	DependsOn []meta.NamespacedObjectReference `json:"dependsOn,omitempty"`

	// Decrypt Kubernetes secrets before applying them on the cluster.
	// +optional
	Decryption *Decryption `json:"decryption,omitempty"`

	// The interval at which to reconcile the Kustomization.
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ms|s|m|h))+$"
	// +required
	Interval metav1.Duration `json:"interval"`

	// The interval at which to retry a previously failed reconciliation.
	// When not specified, the controller uses the KustomizationSpec.Interval
	// value to retry failures.
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ms|s|m|h))+$"
	// +optional
	RetryInterval *metav1.Duration `json:"retryInterval,omitempty"`

	// The KubeConfig for reconciling the Kustomization on a remote cluster.
	// When used in combination with KustomizationSpec.ServiceAccountName,
	// forces the controller to act on behalf of that Service Account at the
	// target cluster.
	// If the --default-service-account flag is set, its value will be used as
	// a controller level fallback for when KustomizationSpec.ServiceAccountName
	// is empty.
	// +optional
	KubeConfig *meta.KubeConfigReference `json:"kubeConfig,omitempty"`

	// Path to the directory containing the kustomization.yaml file, or the
	// set of plain YAMLs a kustomization.yaml should be generated for.
	// Defaults to 'None', which translates to the root path of the SourceRef.
	// +optional
	Path string `json:"path,omitempty"`

	// PostBuild describes which actions to perform on the YAML manifest
	// generated by building the kustomize overlay.
	// +optional
	PostBuild *PostBuild `json:"postBuild,omitempty"`

	// Prune enables garbage collection.
	// +required
	Prune bool `json:"prune"`

	// A list of resources to be included in the health assessment.
	// +optional
	HealthChecks []meta.NamespacedObjectKindReference `json:"healthChecks,omitempty"`

	// NamePrefix sets or overrides the namePrefix in the
	// kustomization.yaml file.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=252
	// +kubebuilder:validation:Optional
	// +optional
	NamePrefix string `json:"namePrefix,omitempty"`

	// NameSuffix sets or overrides the nameSuffix in the
	// kustomization.yaml file.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=252
	// +kubebuilder:validation:Optional
	// +optional
	NameSuffix string `json:"nameSuffix,omitempty"`

	// Strategic merge and JSON patches, defined as inline YAML objects,
	// capable of targeting objects based on kind, label and annotation selectors.
	// +optional
	Patches []kustomize.Patch `json:"patches,omitempty"`

	// Strategic merge patches, defined as inline YAML objects.
	// Deprecated: Use Patches instead.
	// +optional
	PatchesStrategicMerge []apiextensionsv1.JSON `json:"patchesStrategicMerge,omitempty"`

	// JSON 6902 patches, defined as inline YAML objects.
	// Deprecated: Use Patches instead.
	// +optional
	PatchesJSON6902 []kustomize.JSON6902Patch `json:"patchesJson6902,omitempty"`

	// Images is a list of (image name, new name, new tag or digest)
	// for changing image names, tags or digests. This can also be achieved with a
	// patch, but this operator is simpler to specify.
	// +optional
	Images []kustomize.Image `json:"images,omitempty"`

	// The name of the Kubernetes service account to impersonate
	// when reconciling this Kustomization.
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`

	// Reference of the source where the kustomization file is.
	// +required
	SourceRef CrossNamespaceSourceReference `json:"sourceRef"`

	// This flag tells the controller to suspend subsequent kustomize executions,
	// it does not apply to already started executions. Defaults to false.
	// +optional
	Suspend bool `json:"suspend,omitempty"`

	// TargetNamespace sets or overrides the namespace in the
	// kustomization.yaml file.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Optional
	// +optional
	TargetNamespace string `json:"targetNamespace,omitempty"`

	// Timeout for validation, apply and health checking operations.
	// Defaults to 'Interval' duration.
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ms|s|m|h))+$"
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// Force instructs the controller to recreate resources
	// when patching fails due to an immutable field change.
	// +kubebuilder:default:=false
	// +optional
	Force bool `json:"force,omitempty"`

	// Wait instructs the controller to check the health of all the reconciled resources.
	// When enabled, the HealthChecks are ignored. Defaults to false.
	// +optional
	Wait bool `json:"wait,omitempty"`

	// Components specifies relative paths to specifications of other Components.
	// +optional
	Components []string `json:"components,omitempty"`

	// Deprecated: Not used in v1beta2.
	// +kubebuilder:validation:Enum=none;client;server
	// +optional
	Validation string `json:"validation,omitempty"`
}

// CommonMetadata defines the common labels and annotations.
type CommonMetadata struct {
	// Annotations to be added to the object's metadata.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// Labels to be added to the object's metadata.
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
}

// Decryption defines how decryption is handled for Kubernetes manifests.
type Decryption struct {
	// Provider is the name of the decryption engine.
	// +kubebuilder:validation:Enum=sops
	// +required
	Provider string `json:"provider"`

	// The secret name containing the private OpenPGP keys used for decryption.
	// +optional
	SecretRef *meta.LocalObjectReference `json:"secretRef,omitempty"`
}

// PostBuild describes which actions to perform on the YAML manifest
// generated by building the kustomize overlay.
type PostBuild struct {
	// Substitute holds a map of key/value pairs.
	// The variables defined in your YAML manifests
	// that match any of the keys defined in the map
	// will be substituted with the set value.
	// Includes support for bash string replacement functions
	// e.g. ${var:=default}, ${var:position} and ${var/substring/replacement}.
	// +optional
	Substitute map[string]string `json:"substitute,omitempty"`

	// SubstituteFrom holds references to ConfigMaps and Secrets containing
	// the variables and their values to be substituted in the YAML manifests.
	// The ConfigMap and the Secret data keys represent the var names and they
	// must match the vars declared in the manifests for the substitution to happen.
	// +optional
	SubstituteFrom []SubstituteReference `json:"substituteFrom,omitempty"`
}

// SubstituteReference contains a reference to a resource containing
// the variables name and value.
type SubstituteReference struct {
	// Kind of the values referent, valid values are ('Secret', 'ConfigMap').
	// +kubebuilder:validation:Enum=Secret;ConfigMap
	// +required
	Kind string `json:"kind"`

	// Name of the values referent. Should reside in the same namespace as the
	// referring resource.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +required
	Name string `json:"name"`

	// Optional indicates whether the referenced resource must exist, or whether to
	// tolerate its absence. If true and the referenced resource is absent, proceed
	// as if the resource was present but empty, without any variables defined.
	// +kubebuilder:default:=false
	// +optional
	Optional bool `json:"optional,omitempty"`
}

// KustomizationStatus defines the observed state of a kustomization.
type KustomizationStatus struct {
	meta.ReconcileRequestStatus `json:",inline"`

	// ObservedGeneration is the last reconciled generation.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// The last successfully applied revision.
	// Equals the Revision of the applied Artifact from the referenced Source.
	// +optional
	LastAppliedRevision string `json:"lastAppliedRevision,omitempty"`

	// LastAttemptedRevision is the revision of the last reconciliation attempt.
	// +optional
	LastAttemptedRevision string `json:"lastAttemptedRevision,omitempty"`

	// Inventory contains the list of Kubernetes resource object references that have been successfully applied.
	// +optional
	Inventory *ResourceInventory `json:"inventory,omitempty"`
}

// GetTimeout returns the timeout with default.
func (in Kustomization) GetTimeout() time.Duration {
	duration := in.Spec.Interval.Duration - 30*time.Second
	if in.Spec.Timeout != nil {
		duration = in.Spec.Timeout.Duration
	}
	if duration < 30*time.Second {
		return 30 * time.Second
	}
	return duration
}

// GetRetryInterval returns the retry interval
func (in Kustomization) GetRetryInterval() time.Duration {
	if in.Spec.RetryInterval != nil {
		return in.Spec.RetryInterval.Duration
	}
	return in.GetRequeueAfter()
}

// GetRequeueAfter returns the duration after which the Kustomization must be
// reconciled again.
func (in Kustomization) GetRequeueAfter() time.Duration {
	return in.Spec.Interval.Duration
}

// GetDependsOn returns the list of dependencies across-namespaces.
func (in Kustomization) GetDependsOn() []meta.NamespacedObjectReference {
	return in.Spec.DependsOn
}

// GetConditions returns the status conditions of the object.
func (in Kustomization) GetConditions() []metav1.Condition {
	return in.Status.Conditions
}

// SetConditions sets the status conditions on the object.
func (in *Kustomization) SetConditions(conditions []metav1.Condition) {
	in.Status.Conditions = conditions
}

// GetStatusConditions returns a pointer to the Status.Conditions slice.
// Deprecated: use GetConditions instead.
func (in *Kustomization) GetStatusConditions() *[]metav1.Condition {
	return &in.Status.Conditions
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=ks
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description=""
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].status",description=""
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].message",description=""
// +kubebuilder:deprecatedversion:warning="v1beta2 Kustomization is deprecated, upgrade to v1"

// Kustomization is the Schema for the kustomizations API.
type Kustomization struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec KustomizationSpec `json:"spec,omitempty"`
	// +kubebuilder:default:={"observedGeneration":-1}
	Status KustomizationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KustomizationList contains a list of kustomizations.
type KustomizationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Kustomization `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Kustomization{}, &KustomizationList{})
}
