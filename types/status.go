/*
Copyright 2019 The MayaData Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TODO(@amitkumardas):
// Move status & conditions logic to individual controller/xyz
// packages. Let the constants be available in this package/file.

// ConditionType is a custom datatype that
// refers to various conditions supported in this operator
type ConditionType string

const (
	// CStorClusterConfigReconcileErrorCondition is used to
	// indicate presence or absence of error while reconciling
	// CStorClusterConfig
	CStorClusterConfigReconcileErrorCondition ConditionType = "CStorClusterConfigReconcileError"

	// CStorClusterPlanReconcileErrorCondition is used to
	// indicate presence or absence of error while reconciling
	// CStorClusterConfigPlan
	CStorClusterPlanReconcileErrorCondition ConditionType = "CStorClusterPlanReconcileError"

	// CStorClusterStorageSetReconcileErrorCondition is used to
	// indicate presence or absence of error while reconciling
	// CStorClusterStorageSet
	CStorClusterStorageSetReconcileErrorCondition ConditionType = "CStorClusterStorageSetReconcileError"

	// StorageToBlockDeviceAssociationErrorCondition is used to
	// indicate presence or absence of error while reconciling
	// the association of Storage with corresponding BlockDevice
	StorageToBlockDeviceAssociationErrorCondition ConditionType = "StorageToBlockDeviceAssociationError"

	// CStorClusterPlanCSPCApplyErrorCondition is used to indicate
	// presence or absence of error while reconciling
	// the CStorClusterPlan with CStorPoolCluster
	CStorClusterPlanCSPCApplyErrorCondition ConditionType = "CStorClusterPlanCSPCApplyError"
)

// ConditionState is a custom datatype that
// refers to presence or absence of any condition
type ConditionState string

const (
	// ConditionIsPresent refers to presence of any condition
	ConditionIsPresent ConditionState = "True"

	// ConditionIsAbsent refers to absence of any condition
	ConditionIsAbsent ConditionState = "False"
)

// StatusPhase refers to various phases found in a resource'
// status
type StatusPhase string

const (
	// StatusPhaseError refers to a generic error status phase
	StatusPhaseError StatusPhase = "Error"
)

// DeviceClaimState defines the observed claim state of BlockDevice
type DeviceClaimState string

const (
	// BlockDeviceUnclaimed represents that the block device is
	// not bound to any BDC, all cleanup jobs have been completed
	// and is available for claiming.
	BlockDeviceUnclaimed DeviceClaimState = "Unclaimed"

	// BlockDeviceClaimed represents that the block device is
	// bound to any BDC currently
	BlockDeviceClaimed DeviceClaimState = "Claimed"
)

// DeviceState defines the observed state of BlockDevice
type DeviceState string

const (
	// BlockDeviceActive represents a Active BlockDevice
	BlockDeviceActive DeviceState = "Active"

	// BlockDeviceInactive represents an Inactive BlockDevice
	BlockDeviceInactive DeviceState = "Inactive"
)

// now returns the current time in following format
// 2006-01-02 15:04:05.000000
func now() string {
	return metav1.Now().Format("2006-01-02 15:04:05.000000")
}

// MakeCStorClusterConfigReconcileErrCond builds a new
// CStorClusterConfigConditionReconcileError condition
// suitable to be used in API status.conditions
func MakeCStorClusterConfigReconcileErrCond(err error) map[string]interface{} {
	return map[string]interface{}{
		"type":             CStorClusterConfigReconcileErrorCondition,
		"status":           ConditionIsPresent,
		"reason":           err.Error(),
		"lastObservedTime": now(),
	}
}

// MakeCStorClusterPlanReconcileErrCond builds a new
// CStorClusterPlanConditionReconcileError condition
// suitable to be used in API status.conditions
func MakeCStorClusterPlanReconcileErrCond(err error) map[string]interface{} {
	return map[string]interface{}{
		"type":             CStorClusterPlanReconcileErrorCondition,
		"status":           ConditionIsPresent,
		"reason":           err.Error(),
		"lastObservedTime": now(),
	}
}

// MakeCStorClusterStorageSetReconcileErrCond builds a new
// CStorClusterStorageSetConditionReconcileError condition
// suitable to be used in API status.conditions
func MakeCStorClusterStorageSetReconcileErrCond(err error) map[string]interface{} {
	return map[string]interface{}{
		"type":             CStorClusterStorageSetReconcileErrorCondition,
		"status":           ConditionIsPresent,
		"reason":           err.Error(),
		"lastObservedTime": now(),
	}
}

// MakeCStorClusterPlanCSPCApplyErrCond builds a new
// CStorPoolClusterApplyErrorCondition suitable to be
// used in API status.conditions
func MakeCStorClusterPlanCSPCApplyErrCond(err error) map[string]interface{} {
	return map[string]interface{}{
		"type":             CStorClusterPlanCSPCApplyErrorCondition,
		"status":           ConditionIsPresent,
		"reason":           err.Error(),
		"lastObservedTime": now(),
	}
}

// MakeStorageToBlockDeviceAssociationErrCond builds a new
// StorageToBlockDeviceAssociationErrorCondition suitable to
// be used in API status.conditions
func MakeStorageToBlockDeviceAssociationErrCond(err error) map[string]interface{} {
	return map[string]interface{}{
		"type":             StorageToBlockDeviceAssociationErrorCondition,
		"status":           ConditionIsPresent,
		"reason":           err.Error(),
		"lastObservedTime": now(),
	}
}

// MakeNoCStorClusterConfigReconcileErrCond builds a new no
// CStorClusterConfigConditionReconcileError condition. This
// should be used in such a way that it voids previous occurrence of
// this error if any.
func MakeNoCStorClusterConfigReconcileErrCond() map[string]interface{} {
	return map[string]interface{}{
		"type":             string(CStorClusterConfigReconcileErrorCondition),
		"status":           string(ConditionIsAbsent),
		"lastObservedTime": now(),
	}
}

// MergeNoReconcileErrorOnCStorClusterConfig sets
// CStorClusterConfigConditionReconcileError condition to false.
func MergeNoReconcileErrorOnCStorClusterConfig(obj *CStorClusterConfig) {
	noErrCond := CStorClusterConfigStatusCondition{
		Type:             CStorClusterConfigReconcileErrorCondition,
		Status:           ConditionIsAbsent,
		LastObservedTime: now(),
	}
	var newConds []CStorClusterConfigStatusCondition
	for _, old := range obj.Status.Conditions {
		if old.Type == CStorClusterConfigReconcileErrorCondition {
			// ignore previous occurrence of ReconcileError
			continue
		}
		newConds = append(newConds, old)
	}
	newConds = append(newConds, noErrCond)
	obj.Status.Conditions = newConds
}

// MergeNoReconcileErrorOnCStorClusterPlan sets
// CStorClusterPlanConditionReconcileError condition to false.
func MergeNoReconcileErrorOnCStorClusterPlan(obj *CStorClusterPlan) {
	noErrCond := CStorClusterPlanStatusCondition{
		Type:             CStorClusterPlanReconcileErrorCondition,
		Status:           ConditionIsAbsent,
		LastObservedTime: now(),
	}
	var newConds []CStorClusterPlanStatusCondition
	for _, old := range obj.Status.Conditions {
		if old.Type == CStorClusterPlanReconcileErrorCondition {
			// ignore previous occurrence of ReconcileError
			continue
		}
		newConds = append(newConds, old)
	}
	newConds = append(newConds, noErrCond)
	obj.Status.Conditions = newConds
}

// MergeNoCSPCApplyErrorOnCStorClusterPlan sets
// CStorClusterPlanConditionReconcileError condition to false.
func MergeNoCSPCApplyErrorOnCStorClusterPlan(obj *CStorClusterPlan) {
	noErrCond := CStorClusterPlanStatusCondition{
		Type:             CStorClusterPlanCSPCApplyErrorCondition,
		Status:           ConditionIsAbsent,
		LastObservedTime: now(),
	}
	var newConds []CStorClusterPlanStatusCondition
	for _, old := range obj.Status.Conditions {
		if old.Type == CStorClusterPlanCSPCApplyErrorCondition {
			// ignore previous occurrence of ReconcileError
			continue
		}
		newConds = append(newConds, old)
	}
	newConds = append(newConds, noErrCond)
	obj.Status.Conditions = newConds
}

// MergeNoReconcileErrorOnCStorClusterStorageSet sets
// CStorClusterStorageSetConditionReconcileError condition to false.
func MergeNoReconcileErrorOnCStorClusterStorageSet(obj *CStorClusterStorageSet) {
	noErrCond := CStorClusterStorageSetStatusCondition{
		Type:             CStorClusterStorageSetReconcileErrorCondition,
		Status:           ConditionIsAbsent,
		LastObservedTime: now(),
	}
	var newConds []CStorClusterStorageSetStatusCondition
	for _, old := range obj.Status.Conditions {
		if old.Type == CStorClusterStorageSetReconcileErrorCondition {
			// ignore previous occurrence of ReconcileError
			continue
		}
		newConds = append(newConds, old)
	}
	newConds = append(newConds, noErrCond)
	obj.Status.Conditions = newConds
}

// MakeCStorClusterPlanToOnlineWithNoReconcileErr sets the given
// CStorClusterPlan status to online and returns this newly
// formed status object.
func MakeCStorClusterPlanToOnlineWithNoReconcileErr(obj *CStorClusterPlan) map[string]interface{} {
	MergeNoReconcileErrorOnCStorClusterPlan(obj)
	obj.Status.Phase = CStorClusterPlanStatusPhaseOnline
	return map[string]interface{}{
		"phase":      CStorClusterPlanStatusPhaseOnline,
		"conditions": obj.Status.Conditions,
	}
}

// MakeCStorClusterStorageSetToOnlineWithNoReconcileErr sets the given
// CStorClusterStorageSet status to online and returns this newly
// formed status object.
func MakeCStorClusterStorageSetToOnlineWithNoReconcileErr(obj *CStorClusterStorageSet) map[string]interface{} {
	MergeNoReconcileErrorOnCStorClusterStorageSet(obj)
	obj.Status.Phase = CStorClusterStorageSetStatusPhaseOnline
	return map[string]interface{}{
		"phase":      CStorClusterStorageSetStatusPhaseOnline,
		"conditions": obj.Status.Conditions,
	}
}
