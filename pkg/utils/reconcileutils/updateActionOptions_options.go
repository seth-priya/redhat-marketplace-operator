package reconcileutils

// Code generated by github.com/launchdarkly/go-options.  DO NOT EDIT.

import (
	"github.com/redhat-marketplace/redhat-marketplace-operator/pkg/utils/patch"
)

type ApplyUpdateActionOptionFunc func(c *updateActionOptions) error

func (f ApplyUpdateActionOptionFunc) apply(c *updateActionOptions) error {
	return f(c)
}

func newUpdateActionOptions(options ...UpdateActionOption) (updateActionOptions, error) {
	var c updateActionOptions
	err := applyUpdateActionOptionsOptions(&c, options...)
	return c, err
}

func applyUpdateActionOptionsOptions(c *updateActionOptions, options ...UpdateActionOption) error {
	for _, o := range options {
		if err := o.apply(c); err != nil {
			return err
		}
	}
	return nil
}

type UpdateActionOption interface {
	apply(*updateActionOptions) error
}

func UpdateStatusOnly(o bool) ApplyUpdateActionOptionFunc {
	return func(c *updateActionOptions) error {
		c.StatusOnly = o
		return nil
	}
}

func UpdateWithPatch(o patch.PatchAnnotator) ApplyUpdateActionOptionFunc {
	return func(c *updateActionOptions) error {
		c.WithPatch = o
		return nil
	}
}
