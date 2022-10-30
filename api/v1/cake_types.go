/*
Copyright 2022.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CakeSpec defines the desired state of Cake
type CakeSpec struct {
	// Number of replicas for the Nginx Pods
	ReplicaCount int32 `json:"replicaCount"`
	// Exposed port for the Nginx server
	Port int32 `json:"port"`
	//COLOUR can be one of "white" or "colour"
	COLOUR string `json:"COLOUR"`
	//Decoration can be one of "ghost" or "heart"
	DECORATION string `json:"DECORATION"`
	MESSAGE    string `json:"MESSAGE"`
}

// CakeStatus defines the observed state of Cake
type CakeStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Cake is the Schema for the cakes API
type Cake struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CakeSpec   `json:"spec,omitempty"`
	Status CakeStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CakeList contains a list of Cake
type CakeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cake `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Cake{}, &CakeList{})
}
