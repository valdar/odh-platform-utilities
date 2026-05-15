package validation_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/opendatahub-io/odh-platform-utilities/api/common"
	"github.com/opendatahub-io/odh-platform-utilities/api/common/validation"
)

func TestValidatePlatformObject_ValidImplementation(t *testing.T) {
	t.Parallel()

	obj := &fakeModule{
		ObjectMeta: metav1.ObjectMeta{Name: "test-module"},
	}

	validation.ValidatePlatformObject(t, obj)
}

func TestValidate_ValidImplementation(t *testing.T) {
	t.Parallel()

	obj := &fakeModule{
		ObjectMeta: metav1.ObjectMeta{Name: "test-module"},
	}

	err := validation.Validate(obj)
	require.NoError(t, err)
}

func TestValidate_NilObject(t *testing.T) {
	t.Parallel()

	err := validation.Validate(nil)
	require.Error(t, err)
}

// --- fakeModule: correct PlatformObject implementation -------------------

type fakeModuleStatus struct {
	common.ComponentReleaseStatus `json:",inline"`
	common.Status                 `json:",inline"`
}

type fakeModule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Status fakeModuleStatus `json:"status,omitempty"`
}

func (f *fakeModule) GetStatus() *common.Status {
	return &f.Status.Status
}

func (f *fakeModule) GetConditions() []common.Condition {
	return f.Status.Conditions
}

func (f *fakeModule) SetConditions(conditions []common.Condition) {
	f.Status.Conditions = conditions
}

func (f *fakeModule) GetReleaseStatus() *common.ComponentReleaseStatus {
	return &f.Status.ComponentReleaseStatus
}

func (f *fakeModule) SetReleaseStatus(
	status common.ComponentReleaseStatus,
) {
	f.Status.ComponentReleaseStatus = status
}

func (f *fakeModule) DeepCopyObject() runtime.Object {
	cp := &fakeModule{
		TypeMeta: f.TypeMeta,
	}
	f.DeepCopyInto(&cp.ObjectMeta)
	f.Status.Status.DeepCopyInto(&cp.Status.Status)
	f.Status.ComponentReleaseStatus.DeepCopyInto(
		&cp.Status.ComponentReleaseStatus,
	)

	return cp
}
