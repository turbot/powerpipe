package dashboardtypes

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/v2/schema"
	"github.com/turbot/powerpipe/internal/resources"
)

// ResolvedRuntimeDependency is a wrapper for RuntimeDependency which contains the resolved value
// we must wrap it so that we do not mutate the underlying workspace data when resolving dependency values
type ResolvedRuntimeDependency struct {
	Dependency *resources.RuntimeDependency
	valueLock  sync.Mutex
	Value      any
	// the name of the run which publishes this dependency
	publisherName string
	valueChannel  chan *ResolvedRuntimeDependencyValue
}

func NewResolvedRuntimeDependency(dep *resources.RuntimeDependency, valueChannel chan *ResolvedRuntimeDependencyValue, publisherName string) *ResolvedRuntimeDependency {
	return &ResolvedRuntimeDependency{
		Dependency:    dep,
		valueChannel:  valueChannel,
		publisherName: publisherName,
	}
}

// ScopedName returns is a unique name for the dependency by prepending the publisher name
// this is used to uniquely identify which `with` is used - for the snapshot data
func (d *ResolvedRuntimeDependency) ScopedName() string {
	return fmt.Sprintf("%s.%s", d.publisherName, d.Dependency.SourceResourceName())
}

func (d *ResolvedRuntimeDependency) IsResolved() bool {
	d.valueLock.Lock()
	defer d.valueLock.Unlock()

	return d.hasValue()
}

func (d *ResolvedRuntimeDependency) Resolve() error {
	d.valueLock.Lock()
	defer d.valueLock.Unlock()

	slog.Debug("ResolvedRuntimeDependency Resolve", "dep", d.Dependency.PropertyPath, "chan", d.valueChannel)

	// if we are already resolved, do nothing
	if d.hasValue() {
		return nil
	}

	// wait for value
	val := <-d.valueChannel

	if val.Error != nil {
		return val.Error
	}

	// set the value - apply the property path (if any
	if err := d.setValue(val.Value); err != nil {
		return err
	}

	// we should have a non nil value now
	if !d.hasValue() {
		return fmt.Errorf("nil value received for runtime dependency %s", d.Dependency.String())
	}

	// TACTICAL if the desired target value is an array, wrap in an array
	if d.Dependency.IsArray {
		d.Value = helpers.AnySliceToTypedSlice([]any{d.Value})
	}

	return nil
}

// setValue sets the value for the dependency - if a property path is provided, treat the value as a map and
// dereference the property path
func (d *ResolvedRuntimeDependency) setValue(value any) error {
	// there is special case logic for inputs
	if d.Dependency.PropertyPath.ItemType != schema.BlockTypeInput {
		// for any other runtime dependency, just set the value
		d.Value = value
	}

	if len(d.Dependency.PropertyPath.PropertyPath) == 0 || d.Dependency.PropertyPath.PropertyPath[0] != "value" {
		return fmt.Errorf("property path does not start with 'value' for input dependency %s", d.Dependency.String())
	}
	// is there any further property path, other than 'value'
	if len(d.Dependency.PropertyPath.PropertyPath) == 1 {
		// no property path - just set the value
		d.Value = value
	}

	propertyPath := d.Dependency.PropertyPath.PropertyPath[1:]
	for _, pathSegment := range propertyPath {
		// if there is a property path, we expect the value to be a map
		valueMap, ok := value.(map[string]any)
		if !ok {
			return fmt.Errorf("expected map value for input dependency %s, got %T", d.Dependency.String(), value)
		}
		value, ok = valueMap[pathSegment]
		if !ok {
			return fmt.Errorf("missing property %s for input dependency %s", pathSegment, d.Dependency.String())
		}
	}
	d.Value = value
	return nil

}

func (d *ResolvedRuntimeDependency) hasValue() bool {
	return !helpers.IsNil(d.Value)
}
