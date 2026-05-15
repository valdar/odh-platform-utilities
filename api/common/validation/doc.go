// Package validation verifies that PlatformObject implementations satisfy
// behavioral contracts the Go type system cannot enforce, such as nil
// safety, data persistence, and pointer semantics.
//
// Use [Validate] for runtime checks (webhooks, startup) and
// [ValidatePlatformObject] for test assertions.
//
// # Example (test)
//
//	func TestMyComponent_PlatformObject(t *testing.T) {
//	    obj := &v1alpha1.MyComponent{
//	        ObjectMeta: metav1.ObjectMeta{Name: "default"},
//	    }
//	    validation.ValidatePlatformObject(t, obj)
//	}
//
// # Example (runtime)
//
//	if err := validation.Validate(obj); err != nil {
//	    log.Fatalf("PlatformObject contract violated: %v", err)
//	}
package validation
