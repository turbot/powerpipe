package resources

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/pipe-fittings/modconfig"
)

// enrich the shell category by fetching from the ResourceMapsProvider
// this is used when a category has been retrieved via a HCL reference - as cty does not serialise all properties
func enrichCategory(shellCategory *DashboardCategory, parent modconfig.HclResource, resourceMapProvider modconfig.ResourceMapsProvider) (*DashboardCategory, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	resourceMaps := resourceMapProvider.GetResourceMaps().(*ModResources)
	fullCategory, ok := resourceMaps.DashboardCategories[shellCategory.Name()]
	if !ok {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("%s contains edge %s but this has not been loaded", parent.Name(), shellCategory.Name()),
			Subject:  parent.GetDeclRange(),
		})
		return nil, diags
	}
	return fullCategory, diags
}
