package validation

import (
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/opendatahub-io/odh-platform-utilities/api/common"
)

// --- negative tests: broken implementations that violate the contract ----

func TestCheckGetStatus_NilStatus(t *testing.T) {
	t.Parallel()

	obj := &brokenNilStatus{
		ObjectMeta: metav1.ObjectMeta{Name: "broken"},
	}

	err := checkGetStatus(obj)
	require.ErrorIs(t, err, ErrNilStatus)
}

func TestCheckConditionsRoundTrip_NotStored(t *testing.T) {
	t.Parallel()

	obj := &brokenConditionsNotStored{
		ObjectMeta: metav1.ObjectMeta{Name: "broken"},
	}

	err := checkConditionsRoundTrip(obj)
	require.ErrorIs(t, err, ErrCondRoundTrip)
}

func TestCheckMandatoryConditionTypes_NotStored(t *testing.T) {
	t.Parallel()

	obj := &brokenConditionsNotStored{
		ObjectMeta: metav1.ObjectMeta{Name: "broken"},
	}

	err := checkMandatoryConditionTypes(obj)
	require.ErrorIs(t, err, ErrMandatoryConds)
}

func TestCheckReleaseStatusRoundTrip_NilRelease(t *testing.T) {
	t.Parallel()

	obj := &brokenNilReleaseStatus{
		ObjectMeta: metav1.ObjectMeta{Name: "broken"},
	}

	err := checkReleaseStatusRoundTrip(obj)
	require.ErrorIs(t, err, ErrNilReleaseStatus)
}

func TestCheckReleaseStatusRoundTrip_NotStored(t *testing.T) {
	t.Parallel()

	obj := &brokenReleasesNotStored{
		ObjectMeta: metav1.ObjectMeta{Name: "broken"},
	}

	err := checkReleaseStatusRoundTrip(obj)
	require.ErrorIs(t, err, ErrRelRoundTrip)
}

func TestCheckPhaseValues_NotWritable(t *testing.T) {
	t.Parallel()

	obj := &brokenPhaseNotWritable{
		ObjectMeta: metav1.ObjectMeta{Name: "broken"},
	}

	err := checkPhaseValues(obj)
	require.ErrorIs(t, err, ErrPhaseField)
}

func TestValidate_ReturnsAllViolations(t *testing.T) {
	t.Parallel()

	obj := &brokenNilStatus{
		ObjectMeta: metav1.ObjectMeta{Name: "broken"},
	}

	err := Validate(obj)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNilStatus)
	require.ErrorIs(t, err, ErrCondRoundTrip)
	require.ErrorIs(t, err, ErrMandatoryConds)
	require.ErrorIs(t, err, ErrNilReleaseStatus)
}

// --- broken implementations ----------------------------------------------

type brokenNilStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
}

func (b *brokenNilStatus) GetStatus() *common.Status { return nil }
func (b *brokenNilStatus) GetConditions() []common.Condition {
	return nil
}
func (b *brokenNilStatus) SetConditions(_ []common.Condition) {}
func (b *brokenNilStatus) GetReleaseStatus() *common.ComponentReleaseStatus {
	return nil
}
func (b *brokenNilStatus) SetReleaseStatus(
	_ common.ComponentReleaseStatus,
) {
}
func (b *brokenNilStatus) DeepCopyObject() runtime.Object {
	return b
}

// --- brokenNilReleaseStatus: release status accessor returns nil ---------

type brokenNilReleaseStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	status common.Status
}

func (b *brokenNilReleaseStatus) GetStatus() *common.Status {
	return &b.status
}
func (b *brokenNilReleaseStatus) GetConditions() []common.Condition {
	return b.status.Conditions
}
func (b *brokenNilReleaseStatus) SetConditions(
	c []common.Condition,
) {
	b.status.Conditions = c
}
func (b *brokenNilReleaseStatus) GetReleaseStatus() *common.ComponentReleaseStatus {
	return nil
}
func (b *brokenNilReleaseStatus) SetReleaseStatus(
	_ common.ComponentReleaseStatus,
) {
}
func (b *brokenNilReleaseStatus) DeepCopyObject() runtime.Object {
	return b
}

// --- brokenConditionsNotStored: condition setter is a no-op --------------

type brokenConditionsNotStored struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	releases common.ComponentReleaseStatus
	status   common.Status
}

func (b *brokenConditionsNotStored) GetStatus() *common.Status {
	return &b.status
}
func (b *brokenConditionsNotStored) GetConditions() []common.Condition {
	return nil
}
func (b *brokenConditionsNotStored) SetConditions(
	_ []common.Condition,
) {
}
func (b *brokenConditionsNotStored) GetReleaseStatus() *common.ComponentReleaseStatus {
	return &b.releases
}
func (b *brokenConditionsNotStored) SetReleaseStatus(
	s common.ComponentReleaseStatus,
) {
	b.releases = s
}
func (b *brokenConditionsNotStored) DeepCopyObject() runtime.Object {
	return b
}

// --- brokenReleasesNotStored: release setter is a no-op ------------------

type brokenReleasesNotStored struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	releases common.ComponentReleaseStatus
	status   common.Status
}

func (b *brokenReleasesNotStored) GetStatus() *common.Status {
	return &b.status
}
func (b *brokenReleasesNotStored) GetConditions() []common.Condition {
	return b.status.Conditions
}
func (b *brokenReleasesNotStored) SetConditions(
	c []common.Condition,
) {
	b.status.Conditions = c
}
func (b *brokenReleasesNotStored) GetReleaseStatus() *common.ComponentReleaseStatus {
	return &b.releases
}
func (b *brokenReleasesNotStored) SetReleaseStatus(
	_ common.ComponentReleaseStatus,
) {
}
func (b *brokenReleasesNotStored) DeepCopyObject() runtime.Object {
	return b
}

// --- brokenPhaseNotWritable: status accessor returns a copy --------------

type brokenPhaseNotWritable struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	releases common.ComponentReleaseStatus
	status   common.Status
}

func (b *brokenPhaseNotWritable) GetStatus() *common.Status {
	cp := b.status

	return &cp
}
func (b *brokenPhaseNotWritable) GetConditions() []common.Condition {
	return b.status.Conditions
}
func (b *brokenPhaseNotWritable) SetConditions(
	c []common.Condition,
) {
	b.status.Conditions = c
}
func (b *brokenPhaseNotWritable) GetReleaseStatus() *common.ComponentReleaseStatus {
	return &b.releases
}
func (b *brokenPhaseNotWritable) SetReleaseStatus(
	s common.ComponentReleaseStatus,
) {
	b.releases = s
}
func (b *brokenPhaseNotWritable) DeepCopyObject() runtime.Object {
	return b
}
