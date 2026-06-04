package api

import (
	"github.com/blang/semver/v4"
	common "github.com/opendatahub-io/odh-platform-utilities/api/common"
)

type Condition = common.Condition
type ConditionSeverity = common.ConditionSeverity
type ConditionType = common.ConditionType
type Status = common.Status
type Phase = common.Phase
type ConditionsAccessor = common.ConditionsAccessor
type WithStatus = common.WithStatus
type PlatformObject = common.PlatformObject

const (
	ConditionSeverityError = common.ConditionSeverityError
	ConditionSeverityInfo  = common.ConditionSeverityInfo
	ConditionTypeReady     = common.ConditionTypeReady
)

const ConditionReasonError = "Error"

// Platform identifies the operator distribution.
type Platform string

// Release includes information on operator version and platform.
// +kubebuilder:object:generate=true
type Release struct {
	Name    Platform       `json:"name,omitempty"`
	Version semver.Version `json:"version,omitempty"`
}

func (in *Release) DeepCopyInto(out *Release) {
	*out = *in
	out.Version = in.Version
}

func (in *Release) DeepCopy() *Release {
	if in == nil {
		return nil
	}
	out := new(Release)
	in.DeepCopyInto(out)
	return out
}
