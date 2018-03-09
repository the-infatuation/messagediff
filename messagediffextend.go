package messagediff

import (
	"reflect"
)

type Option struct {
	name  string
	value interface{}
}

type MergedOptions struct {
	allSets bool
	fields  []FieldData
	compare Comparator
}

type ComparisonResult string

const (
	Modified = "modified"
	Matched  = "matched"
	Recurse  = "recurse"
)

type FieldData struct {
	name string
	kind reflect.Type
}

type Comparator interface {
	Compare(interface{}, interface{}) ComparisonResult
}

// Field
func Field(x interface{}, name string) *FieldData {
	return &FieldData{
		name: name,
		kind: reflect.TypeOf(x),
	}
}

// TreatRepeatedFieldsAsSets
func TreatRepeatedFieldsAsSets() *Option {
	return &Option{"all", true}
}

// TreatAsSet
func TreatAsSet(f *FieldData) *Option {
	return &Option{"field", f}
}

// CompareUsing
func CompareUsing(c Comparator) *Option {
	return &Option{"compare", c}
}

// DeepDiffWithOptions
func DeepDiffWithOptions(a, b interface{}, opts ...*Option) (*Diff, bool) {
	d := newDiff()
	var p = &MergedOptions{}
	for _, o := range opts {
		switch o.name {
		case "all":
			p.allSets = o.value.(bool)
		case "field":
			p.fields = append(p.fields, *o.value.(*FieldData))
		case "compare":
			p.compare = *o.value.(*Comparator)
		default:
		}
	}
	return d, d.diff(reflect.ValueOf(a), reflect.ValueOf(b), nil, p)
}

func contains(sv reflect.Value, val interface{}) bool {
	for i := 0; i < sv.Len(); i++ {
		if _, ok := DeepDiff(sv.Index(i).Interface(), val); ok {
			return true
		}
	}
	return false
}
