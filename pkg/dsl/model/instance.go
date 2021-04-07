/*
Copyright 2021 The KubeVela Authors.

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

package model

import (
	"strings"

	"github.com/pkg/errors"

	"cuelang.org/go/cue/build"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/format"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/klog"

	"github.com/oam-dev/kubevela/pkg/dsl/model/sets"
)

// Instance defines Model Interface
type Instance interface {
	String() string
	Unstructured() (*unstructured.Unstructured, error)
	IsBase() bool
	Unify(other Instance) error
	Compile() ([]byte, error)
}

type instance struct {
	v    string
	base bool
}

// String return instance's cue format string
func (inst *instance) String() string {
	return inst.v
}

// IsBase indicate whether the instance is base model
func (inst *instance) IsBase() bool {
	return inst.base
}

func (inst *instance) Compile() ([]byte, error) {
	bi := build.NewContext().NewInstance("", nil)
	err := bi.AddFile("-", inst.v)
	if err != nil {
		return nil, err
	}
	var r cue.Runtime
	it, err := r.Build(bi)
	if err != nil {
		return nil, err
	}
	// compiled object should be final and concrete value
	if err := it.Value().Validate(cue.Concrete(true), cue.Final()); err != nil {
		return nil, it.Err
	}
	return it.Value().MarshalJSON()
}

// Unstructured convert cue values to unstructured.Unstructured
// TODO(wonderflow): will it be better if we try to decode it to concrete object(such as K8s Deployment) by using runtime.Schema?
func (inst *instance) Unstructured() (*unstructured.Unstructured, error) {
	jsonv, err := inst.Compile()
	if err != nil {
		return nil, err
	}
	o := &unstructured.Unstructured{}
	if err := o.UnmarshalJSON(jsonv); err != nil {
		return nil, err
	}
	return o, nil

}

// Unify implement unity operations between instances
func (inst *instance) Unify(other Instance) error {
	pv, err := sets.StrategyUnify(inst.v, other.String())
	if err != nil {
		return err
	}
	klog.Infof("XXX patch value %v XXX\n", pv)
	if strings.Contains(pv, "_|_") {
		return errors.New(IndexMatchLine(pv, "_|_"))
	}
	inst.v = pv
	return nil
}

// NewBase create a base instance
func NewBase(v cue.Value) (Instance, error) {
	vs, err := openPrint(v)
	if err != nil {
		return nil, err
	}
	return &instance{
		v:    vs,
		base: true,
	}, nil
}

// NewOther create a non-base instance
func NewOther(v cue.Value) (Instance, error) {
	vs, err := openPrint(v)
	if err != nil {
		return nil, err
	}
	return &instance{
		v: vs,
	}, nil
}

func openPrint(v cue.Value) (string, error) {
	sysopts := []cue.Option{cue.All(), cue.DisallowCycles(true), cue.ResolveReferences(true), cue.Docs(true)}
	f, err := sets.ToFile(v.Syntax(sysopts...))
	if err != nil {
		return "", err
	}
	for _, decl := range f.Decls {
		listOpen(decl)
	}

	ret, err := format.Node(f)
	if err != nil {
		return "", err
	}
	if strings.Contains(string(ret), "_|_") {
		return "", errors.New(IndexMatchLine(string(ret), "_|_"))
	}
	return string(ret), nil
}

// IndexMatchLine will index and extract the line contains the pattern.
func IndexMatchLine(res string, pattern string) string {
	idx := strings.Index(res, pattern)
	if idx < 0 || idx >= len(res) {
		// should not happen if we check contains first.
		return ""
	}
	var start = -1
	var end = len(res)
	for i := idx; i >= 0; i-- {
		if res[i] == '\n' {
			start = i
			break
		}
	}
	for i := idx; i < len(res); i++ {
		if res[i] == '\n' {
			end = i
			break
		}
	}
	return res[start+1 : end]
}

func listOpen(expr ast.Node) {
	switch v := expr.(type) {
	case *ast.Field:
		listOpen(v.Value)
	case *ast.StructLit:
		for _, elt := range v.Elts {
			listOpen(elt)
		}
	case *ast.BinaryExpr:
		listOpen(v.X)
		listOpen(v.Y)
	case *ast.EmbedDecl:
		listOpen(v.Expr)
	case *ast.Comprehension:
		listOpen(v.Value)
	case *ast.ListLit:
		for _, elt := range v.Elts {
			listOpen(elt)
		}
		if len(v.Elts) > 0 {
			if _, ok := v.Elts[len(v.Elts)-1].(*ast.Ellipsis); !ok {
				v.Elts = append(v.Elts, &ast.Ellipsis{})
			}
		}
	}
}
