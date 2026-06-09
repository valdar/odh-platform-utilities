package kustomize_test

import (
	"context"
	"testing"

	common "github.com/opendatahub-io/odh-platform-utilities/api/common"
	"github.com/opendatahub-io/odh-platform-utilities/framework/api"
	"github.com/opendatahub-io/odh-platform-utilities/framework/controller/actions/render/kustomize"
	"github.com/opendatahub-io/odh-platform-utilities/framework/controller/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	k8stypes "k8s.io/apimachinery/pkg/types"

	. "github.com/onsi/gomega"
)

type fakeInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	status api.Status
}

func (f *fakeInstance) GetStatus() *api.Status {
	return &f.status
}

func (f *fakeInstance) GetConditions() []api.Condition {
	return f.status.Conditions
}

func (f *fakeInstance) SetConditions(c []api.Condition) {
	f.status.Conditions = c
}

func (f *fakeInstance) GetReleaseStatus() *common.ComponentReleaseStatus {
	return nil
}

func (f *fakeInstance) SetReleaseStatus(_ common.ComponentReleaseStatus) {}

func (f *fakeInstance) DeepCopyObject() runtime.Object {
	o := *f
	return &o
}

func minimalInstance() api.PlatformObject {
	return &fakeInstance{
		TypeMeta:   metav1.TypeMeta{APIVersion: "test/v1", Kind: "Fake"},
		ObjectMeta: metav1.ObjectMeta{Name: "test-instance", UID: k8stypes.UID("uid-1234"), Generation: 1},
	}
}

type renderCall struct {
	path string
	opts []kustomize.RenderOpt
}

type mockEngine struct {
	calls []renderCall
}

func (m *mockEngine) Render(path string, opts ...kustomize.RenderOpt) ([]unstructured.Unstructured, error) {
	m.calls = append(m.calls, renderCall{path: path, opts: opts})

	obj := unstructured.Unstructured{}
	obj.SetKind("ConfigMap")
	obj.SetAPIVersion("v1")
	obj.SetName("rendered-from-" + path)

	return []unstructured.Unstructured{obj}, nil
}

func optsLen(opts []kustomize.RenderOpt) int {
	return len(opts)
}

func TestRenderPerManifestNamespaceOverride(t *testing.T) {
	t.Parallel()

	g := NewWithT(t)
	ctx := context.Background()

	engine := &mockEngine{}
	defaultNS := "default-ns"
	overrideNS := "module-ns"

	action := kustomize.NewAction(
		kustomize.WithEngine(engine),
		kustomize.WithCache(false),
		kustomize.WithNamespaceFn(func(_ context.Context, _ *types.ReconciliationRequest) (string, error) {
			return defaultNS, nil
		}),
	)

	rr := &types.ReconciliationRequest{
		Instance: minimalInstance(),
		Manifests: []types.ManifestInfo{
			{Path: "manifest-a"},
			{Path: "manifest-b", Namespace: overrideNS},
			{Path: "manifest-c"},
		},
	}

	err := action(ctx, rr)
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(rr.Resources).Should(HaveLen(3))

	g.Expect(engine.calls).Should(HaveLen(3))

	g.Expect(optsLen(engine.calls[0].opts)).Should(Equal(1),
		"manifest-a should receive only the default namespace opt")

	g.Expect(optsLen(engine.calls[1].opts)).Should(Equal(2),
		"manifest-b should receive both default and override namespace opts (override wins)")

	g.Expect(optsLen(engine.calls[2].opts)).Should(Equal(1),
		"manifest-c should receive only the default namespace opt")
}

func TestRenderAllManifestsUseDefaultNamespace(t *testing.T) {
	t.Parallel()

	g := NewWithT(t)
	ctx := context.Background()

	engine := &mockEngine{}
	defaultNS := "default-ns"

	action := kustomize.NewAction(
		kustomize.WithEngine(engine),
		kustomize.WithCache(false),
		kustomize.WithNamespaceFn(func(_ context.Context, _ *types.ReconciliationRequest) (string, error) {
			return defaultNS, nil
		}),
	)

	rr := &types.ReconciliationRequest{
		Instance: minimalInstance(),
		Manifests: []types.ManifestInfo{
			{Path: "manifest-a"},
			{Path: "manifest-b"},
		},
	}

	err := action(ctx, rr)
	g.Expect(err).ShouldNot(HaveOccurred())

	g.Expect(engine.calls).Should(HaveLen(2))
	g.Expect(optsLen(engine.calls[0].opts)).Should(Equal(1))
	g.Expect(optsLen(engine.calls[1].opts)).Should(Equal(1))
}

func TestRenderNoNamespaceFnWithOverride(t *testing.T) {
	t.Parallel()

	g := NewWithT(t)
	ctx := context.Background()

	engine := &mockEngine{}
	overrideNS := "module-ns"

	action := kustomize.NewAction(
		kustomize.WithEngine(engine),
		kustomize.WithCache(false),
	)

	rr := &types.ReconciliationRequest{
		Instance: minimalInstance(),
		Manifests: []types.ManifestInfo{
			{Path: "manifest-a"},
			{Path: "manifest-b", Namespace: overrideNS},
		},
	}

	err := action(ctx, rr)
	g.Expect(err).ShouldNot(HaveOccurred())

	g.Expect(engine.calls).Should(HaveLen(2))
	g.Expect(optsLen(engine.calls[0].opts)).Should(Equal(0),
		"manifest-a should have no namespace opts when namespaceFn is nil and no override")

	g.Expect(optsLen(engine.calls[1].opts)).Should(Equal(1),
		"manifest-b should have the override namespace opt even without namespaceFn")
}
