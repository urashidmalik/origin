package policy

import (
	"fmt"

	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/registry/generic"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/util/fielderrors"

	authorizationapi "github.com/openshift/origin/pkg/authorization/api"
	"github.com/openshift/origin/pkg/authorization/api/validation"
)

// strategy implements behavior for nodes
type strategy struct {
	runtime.ObjectTyper
}

// Strategy is the default logic that applies when creating and updating Policy objects.
var Strategy = strategy{kapi.Scheme}

// NamespaceScoped is true for policies.
func (strategy) NamespaceScoped() bool {
	return true
}

// AllowCreateOnUpdate is false for policies.
func (strategy) AllowCreateOnUpdate() bool {
	return false
}

func (strategy) AllowUnconditionalUpdate() bool {
	return false
}

func (strategy) GenerateName(base string) string {
	return base
}

// PrepareForCreate clears fields that are not allowed to be set by end users on creation.
func (strategy) PrepareForCreate(obj runtime.Object) {
	policy := obj.(*authorizationapi.Policy)

	policy.Name = authorizationapi.PolicyName
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (strategy) PrepareForUpdate(obj, old runtime.Object) {
	_ = obj.(*authorizationapi.Policy)
}

// Validate validates a new policy.
func (strategy) Validate(ctx kapi.Context, obj runtime.Object) fielderrors.ValidationErrorList {
	return validation.ValidateLocalPolicy(obj.(*authorizationapi.Policy))
}

// ValidateUpdate is the default update validation for an end user.
func (strategy) ValidateUpdate(ctx kapi.Context, obj, old runtime.Object) fielderrors.ValidationErrorList {
	return validation.ValidateLocalPolicyUpdate(obj.(*authorizationapi.Policy), old.(*authorizationapi.Policy))
}

// Matcher returns a generic matcher for a given label and field selector.
func Matcher(label labels.Selector, field fields.Selector) generic.Matcher {
	return &generic.SelectionPredicate{
		Label: label,
		Field: field,
		GetAttrs: func(obj runtime.Object) (labels.Set, fields.Set, error) {
			policy, ok := obj.(*authorizationapi.Policy)
			if !ok {
				return nil, nil, fmt.Errorf("not a policy")
			}
			return labels.Set(policy.ObjectMeta.Labels), authorizationapi.PolicyToSelectableFields(policy), nil
		},
	}
}
