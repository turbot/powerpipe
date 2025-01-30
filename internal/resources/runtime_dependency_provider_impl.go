package resources

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/pipe-fittings/v2/modconfig"
)

type RuntimeDependencyProviderImpl struct {
	modconfig.ModTreeItemImpl
	// required to allow partial decoding
	RuntimeDependencyProviderRemain hcl.Body `hcl:",remain" json:"-"`

	runtimeDependencies map[string]*RuntimeDependency
}

func (b *RuntimeDependencyProviderImpl) AddRuntimeDependencies(dependencies []*RuntimeDependency) {
	if b.runtimeDependencies == nil {
		b.runtimeDependencies = make(map[string]*RuntimeDependency)
	}
	for _, dependency := range dependencies {
		// set the dependency provider (this is used if this resource is inherited via base)
		dependency.Provider = b
		b.runtimeDependencies[dependency.String()] = dependency
	}
}

func (b *RuntimeDependencyProviderImpl) GetRuntimeDependencies() map[string]*RuntimeDependency {
	return b.runtimeDependencies
}

func (b *RuntimeDependencyProviderImpl) GetNestedStructs() []modconfig.CtyValueProvider {
	// return all nested structs - this is used to get the nested structs for the cty serialisation
	// we return ourselves and our base structs
	return append([]modconfig.CtyValueProvider{b}, b.ModTreeItemImpl.GetNestedStructs()...)
}
