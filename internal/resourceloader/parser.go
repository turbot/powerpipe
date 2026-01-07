package resourceloader

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/turbot/pipe-fittings/v2/funcs"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/pipe-fittings/v2/parse"
	"github.com/turbot/pipe-fittings/v2/schema"
	pparse "github.com/turbot/powerpipe/internal/parse"
	"github.com/turbot/powerpipe/internal/resourceindex"
	"github.com/turbot/powerpipe/internal/resources"
)

// parseResource parses a single resource from its file using the index entry metadata.
func (l *Loader) parseResource(ctx context.Context, entry *resourceindex.IndexEntry) (modconfig.HclResource, error) {
	blockContent, err := l.readResourceBlock(entry)
	if err != nil {
		return nil, fmt.Errorf("reading block: %w", err)
	}

	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(blockContent, entry.FileName)
	if diags.HasErrors() {
		return nil, fmt.Errorf("parsing HCL: %s", diags.Error())
	}

	// Extract the block from the parsed file
	content, _, diags := file.Body.PartialContent(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: entry.Type, LabelNames: []string{"name"}},
		},
	})
	if diags.HasErrors() {
		return nil, fmt.Errorf("extracting block: %s", diags.Error())
	}

	if len(content.Blocks) == 0 {
		return nil, fmt.Errorf("no block found for %s", entry.Name)
	}

	block := content.Blocks[0]

	resource, err := l.decodeResourceBlock(ctx, entry, block)
	if err != nil {
		return nil, err
	}

	// Set source definition on the resource metadata.
	// The blockContent we read is the HCL source for this resource.
	if rwm, ok := resource.(modconfig.ResourceWithMetadata); ok {
		if meta := rwm.GetMetadata(); meta != nil {
			meta.SetSourceDefinition(strings.TrimSpace(string(blockContent)))
		}
	}

	return resource, nil
}

// readResourceBlock reads just the bytes for a single resource block from the file.
func (l *Loader) readResourceBlock(entry *resourceindex.IndexEntry) ([]byte, error) {
	file, err := os.Open(entry.FileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Use byte offset if available for efficient seeking
	if entry.ByteOffset > 0 && entry.ByteLength > 0 {
		if _, err := file.Seek(entry.ByteOffset, io.SeekStart); err != nil {
			return nil, fmt.Errorf("seeking: %w", err)
		}
		content := make([]byte, entry.ByteLength)
		if _, err := io.ReadFull(file, content); err != nil {
			return nil, fmt.Errorf("reading: %w", err)
		}
		return content, nil
	}

	// Fallback to line-based reading
	return l.readByLines(entry)
}

// readByLines reads the resource block line by line using start/end line numbers.
func (l *Loader) readByLines(entry *resourceindex.IndexEntry) ([]byte, error) {
	file, err := os.Open(entry.FileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var content strings.Builder
	lineNum := 0

	for {
		line, err := reader.ReadString('\n')
		lineNum++

		if lineNum >= entry.StartLine && lineNum <= entry.EndLine {
			content.WriteString(line)
		}

		if lineNum >= entry.EndLine || err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}

	return []byte(content.String()), nil
}

// decodeResourceBlock decodes an HCL block into the appropriate resource type.
func (l *Loader) decodeResourceBlock(ctx context.Context, entry *resourceindex.IndexEntry, block *hcl.Block) (modconfig.HclResource, error) {
	// Create the resource using the factory function
	resource, err := l.createResource(block, entry)
	if err != nil {
		return nil, err
	}

	// Create an eval context with standard HCL functions.
	// The root path is the mod root directory, which is needed for file() function.
	rootPath := entry.ModRoot
	if rootPath == "" {
		// Fallback to file directory if mod root not set
		rootPath = filepath.Dir(entry.FileName)
	}
	evalCtx := &hcl.EvalContext{
		Functions: funcs.ContextFunctions(rootPath),
	}

	l.mu.RLock()
	provider := l.resourceProvider
	l.mu.RUnlock()

	diags := parse.DecodeHclBody(block.Body, evalCtx, provider, resource)
	if diags.HasErrors() {
		// For lazy loading, we're more tolerant of decode errors.
		// Variables/locals might not be available, but we can still use the resource
		// for its basic structure (title, children, queries).
		// Only fail on truly critical errors.
		for _, diag := range diags {
			if diag.Severity == hcl.DiagError && !isNonCriticalDecodeError(diag) {
				return nil, fmt.Errorf("decoding resource: %s", diags.Error())
			}
		}
	}

	// Decode nested blocks for dashboards and containers
	if body, ok := block.Body.(*hclsyntax.Body); ok {
		if err := l.decodeNestedBlocks(ctx, resource, body, evalCtx, entry); err != nil {
			return nil, err
		}
	}

	// Initialize inputs map for dashboards after nested blocks are decoded
	// This is needed so that runtime dependencies on inputs can be resolved
	if dashboard, ok := resource.(*resources.Dashboard); ok {
		_ = dashboard.InitInputs()
	}

	// Call OnDecoded to allow the resource to perform post-decode initialization.
	// This is important for NodeAndEdgeProviders (graph/flow/hierarchy) which
	// need to populate NodeNames/EdgeNames from their Nodes/Edges lists.
	if od, ok := resource.(interface {
		OnDecoded(*hcl.Block, modconfig.ModResourcesProvider) hcl.Diagnostics
	}); ok {
		// Pass nil for ModResourcesProvider since we don't have full workspace context
		// during lazy loading. The OnDecoded methods we care about don't require it.
		diags := od.OnDecoded(block, nil)
		// Ignore non-critical diagnostics from OnDecoded
		for _, diag := range diags {
			if diag.Severity == hcl.DiagError && !isNonCriticalDecodeError(diag) {
				return nil, fmt.Errorf("OnDecoded: %s", diag.Summary)
			}
		}
	}

	// Clear the Remain field to save memory
	if clearer, ok := resource.(remainClearer); ok {
		clearer.ClearRemain()
	}

	return resource, nil
}

// isNonCriticalDecodeError checks if a diagnostic is a non-critical error that
// can be ignored for lazy loading. Variable and function reference errors are
// non-critical because the resource structure can still be used.
func isNonCriticalDecodeError(diag *hcl.Diagnostic) bool {
	if diag == nil || diag.Summary == "" {
		return false
	}
	// Variable/local reference errors
	if strings.Contains(diag.Summary, "Variables not allowed") ||
		strings.Contains(diag.Summary, "Unknown variable") ||
		strings.Contains(diag.Summary, "Unsupported attribute") {
		return true
	}
	// Function errors (missing context, etc.)
	if strings.Contains(diag.Summary, "Function calls not allowed") ||
		strings.Contains(diag.Summary, "Call to unknown function") {
		return true
	}
	// Type errors that cascade from variable/function errors
	if strings.Contains(diag.Summary, "Unsuitable value type") ||
		strings.Contains(diag.Summary, "Incorrect attribute value type") {
		return true
	}
	return false
}

// remainClearer interface for resources that can clear their Remain field
type remainClearer interface {
	ClearRemain()
}

// createResource creates a new resource of the appropriate type using factory functions.
func (l *Loader) createResource(block *hcl.Block, entry *resourceindex.IndexEntry) (modconfig.HclResource, error) {
	mod := l.mod

	// If the entry is from a different mod (dependency), create a mod with the correct short name.
	// This ensures resources from dependency mods have the correct FullName.
	if entry.ModName != "" && entry.ModName != l.mod.ShortName {
		mod = modconfig.NewMod(entry.ModName, entry.ModRoot, hcl.Range{})
	}

	// Use shared factory functions registry to ensure consistency with eager loading
	factoryFuncs := pparse.GetResourceFactoryFuncs()
	factoryFunc, ok := factoryFuncs[block.Type]
	if !ok {
		return nil, fmt.Errorf("unsupported block type: %s", block.Type)
	}

	return factoryFunc(block, mod, entry.ShortName), nil
}

// decodeNestedBlocks decodes nested blocks (children) for dashboards and containers.
// The isParentDashboard parameter indicates if the immediate parent is a dashboard
// (as opposed to a container), which determines if children should be considered top-level.
func (l *Loader) decodeNestedBlocks(ctx context.Context, parent modconfig.HclResource, body *hclsyntax.Body, evalCtx *hcl.EvalContext, entry *resourceindex.IndexEntry) error {
	// Check if this resource type supports nested blocks
	parentType := parent.GetBlockType()
	if parentType != schema.BlockTypeDashboard && parentType != schema.BlockTypeContainer {
		return nil
	}

	// Direct children of a dashboard are considered "top-level" for validation purposes
	// (they don't require SQL). Children inside containers are NOT top-level (require SQL).
	isParentDashboard := parentType == schema.BlockTypeDashboard

	// Track child index for generating anonymous names
	childCounts := make(map[string]int)

	for _, b := range body.Blocks {
		block := b.AsHCLBlock()

		// Skip non-dashboard child blocks
		if !isDashboardChildType(block.Type) {
			continue
		}

		// Generate name for anonymous blocks
		shortName := ""
		if len(block.Labels) > 0 {
			shortName = block.Labels[0]
		} else {
			// Anonymous block - use shared utility for consistent naming with eager loading
			shortName = modconfig.AnonymousBlockName(parent.GetUnqualifiedName(), block.Type, childCounts[block.Type])
		}
		childCounts[block.Type]++

		// Create child resource
		childEntry := &resourceindex.IndexEntry{
			Type:      block.Type,
			ShortName: shortName,
			ModRoot:   entry.ModRoot,
			FileName:  entry.FileName,
		}

		childResource, err := l.createResource(block, childEntry)
		if err != nil {
			continue // Skip children that can't be created
		}

		// Set top-level status: direct children of dashboard are "top-level" for validation
		// (they don't require SQL), but children inside containers are not top-level
		childResource.SetTopLevel(isParentDashboard)

		// Decode child's attributes
		l.mu.RLock()
		provider := l.resourceProvider
		l.mu.RUnlock()

		diags := parse.DecodeHclBody(block.Body, evalCtx, provider, childResource)
		// Ignore non-critical errors
		if diags.HasErrors() {
			hasCritical := false
			for _, diag := range diags {
				if diag.Severity == hcl.DiagError && !isNonCriticalDecodeError(diag) {
					hasCritical = true
					break
				}
			}
			if hasCritical {
				continue // Skip children with critical errors
			}
		}

		// Extract runtime dependencies from args attribute
		if syntaxBody, ok := block.Body.(*hclsyntax.Body); ok {
			l.extractRuntimeDependencies(syntaxBody, childResource, evalCtx)
		}

		// Try to resolve base reference if present
		if syntaxBody, ok := block.Body.(*hclsyntax.Body); ok {
			l.resolveBaseReference(ctx, syntaxBody, childResource)
		}

		// Recursively decode nested blocks in this child
		if childBody, ok := block.Body.(*hclsyntax.Body); ok {
			if err := l.decodeNestedBlocks(ctx, childResource, childBody, evalCtx, entry); err != nil {
				return err
			}
			// For flow/graph/hierarchy, also decode node and edge blocks
			if nep, ok := childResource.(resources.NodeAndEdgeProvider); ok {
				if err := l.decodeNodeAndEdgeBlocks(ctx, nep, childBody, evalCtx, entry); err != nil {
					return err
				}
			}
		}

		// Call OnDecoded on child resource for post-decode initialization
		if od, ok := childResource.(interface {
			OnDecoded(*hcl.Block, modconfig.ModResourcesProvider) hcl.Diagnostics
		}); ok {
			_ = od.OnDecoded(block, nil)
		}

		// Validate the child resource has required SQL/query
		if err := validateResource(childResource); err != nil {
			return err
		}

		// Add child to parent based on parent type
		if childItem, ok := childResource.(modconfig.ModTreeItem); ok {
			switch p := parent.(type) {
			case *resources.Dashboard:
				_ = p.AddChild(childItem)
			case *resources.DashboardContainer:
				p.AddChild(childItem)
			}
		}
	}

	return nil
}

// isDashboardChildType returns true if the block type can be a child of a dashboard or container.
func isDashboardChildType(blockType string) bool {
	switch blockType {
	case schema.BlockTypeContainer,
		schema.BlockTypeCard,
		schema.BlockTypeChart,
		schema.BlockTypeFlow,
		schema.BlockTypeGraph,
		schema.BlockTypeHierarchy,
		schema.BlockTypeImage,
		schema.BlockTypeInput,
		schema.BlockTypeTable,
		schema.BlockTypeText,
		schema.BlockTypeWith:
		return true
	default:
		return false
	}
}

// decodeNodeAndEdgeBlocks decodes node and edge blocks inside flow/graph/hierarchy resources.
func (l *Loader) decodeNodeAndEdgeBlocks(ctx context.Context, nep resources.NodeAndEdgeProvider, body *hclsyntax.Body, evalCtx *hcl.EvalContext, entry *resourceindex.IndexEntry) error {
	parentResource := nep.(modconfig.HclResource)
	nodeCount := 0
	edgeCount := 0
	var nodes resources.DashboardNodeList
	var edges resources.DashboardEdgeList

	for _, b := range body.Blocks {
		block := b.AsHCLBlock()

		// Only process node and edge blocks
		if block.Type != schema.BlockTypeNode && block.Type != schema.BlockTypeEdge {
			continue
		}

		// Generate name for this node/edge
		shortName := ""
		if len(block.Labels) > 0 {
			shortName = block.Labels[0]
		} else {
			// Anonymous - generate a name
			if block.Type == schema.BlockTypeNode {
				shortName = fmt.Sprintf("%s_node_%d", parentResource.GetUnqualifiedName(), nodeCount)
				nodeCount++
			} else {
				shortName = fmt.Sprintf("%s_edge_%d", parentResource.GetUnqualifiedName(), edgeCount)
				edgeCount++
			}
		}

		// Create node/edge resource
		childEntry := &resourceindex.IndexEntry{
			Type:      block.Type,
			ShortName: shortName,
			ModRoot:   entry.ModRoot,
			FileName:  entry.FileName,
		}

		childResource, err := l.createResource(block, childEntry)
		if err != nil {
			continue
		}

		// Decode child's attributes
		l.mu.RLock()
		provider := l.resourceProvider
		l.mu.RUnlock()

		_ = parse.DecodeHclBody(block.Body, evalCtx, provider, childResource)

		// Extract runtime dependencies from args attribute
		if syntaxBody, ok := block.Body.(*hclsyntax.Body); ok {
			l.extractRuntimeDependencies(syntaxBody, childResource, evalCtx)
		}

		// Try to resolve base reference if present
		if syntaxBody, ok := block.Body.(*hclsyntax.Body); ok {
			l.resolveBaseReference(ctx, syntaxBody, childResource)
		}

		// Collect nodes and edges
		switch block.Type {
		case schema.BlockTypeNode:
			if node, ok := childResource.(*resources.DashboardNode); ok {
				nodes = append(nodes, node)
			}
		case schema.BlockTypeEdge:
			if edge, ok := childResource.(*resources.DashboardEdge); ok {
				edges = append(edges, edge)
			}
		}
	}

	// Set nodes and edges on the parent
	if len(nodes) > 0 {
		nep.SetNodes(nodes)
	}
	if len(edges) > 0 {
		nep.SetEdges(edges)
	}

	return nil
}

// validateResource validates a resource after decoding.
// This ensures that nested resources have required SQL/query definitions.
func validateResource(resource modconfig.HclResource) error {
	// Check NodeAndEdgeProvider FIRST - flow, graph, hierarchy can have either SQL OR edges/nodes
	// This is more permissive than QueryProvider validation
	if nep, ok := resource.(resources.NodeAndEdgeProvider); ok {
		diags := validateNodeAndEdgeProvider(nep)
		if diags.HasErrors() {
			return fmt.Errorf("%s", diags.Error())
		}
		// Don't also validate as QueryProvider since we've already validated
		return nil
	}

	// Validate QueryProvider resources (chart, table, etc.)
	// Call the resource's own ValidateQuery() method which may have
	// type-specific overrides (e.g., DashboardImage, DashboardCard don't require SQL)
	if qp, ok := resource.(resources.QueryProvider); ok {
		diags := qp.ValidateQuery()
		if diags.HasErrors() {
			return fmt.Errorf("%s", diags.Error())
		}
	}

	return nil
}

// validateNodeAndEdgeProvider validates flow/graph/hierarchy resources.
func validateNodeAndEdgeProvider(nep resources.NodeAndEdgeProvider) hcl.Diagnostics {
	var diags hcl.Diagnostics
	containsEdgesOrNodes := len(nep.GetEdges())+len(nep.GetNodes()) > 0
	definesQuery := nep.GetSQL() != nil || nep.GetQuery() != nil

	// cannot declare both edges/nodes AND sql/query
	if definesQuery && containsEdgesOrNodes {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("%s contains edges/nodes AND has a query", nep.Name()),
			Subject:  nep.GetDeclRange(),
		})
	}

	// if resource is NOT top level must have either edges/nodes OR sql/query
	if !nep.IsTopLevel() && !definesQuery && !containsEdgesOrNodes {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("%s does not define a query or SQL, and has no edges/nodes", nep.Name()),
			Subject:  nep.GetDeclRange(),
		})
	}

	return diags
}

// extractRuntimeDependencies extracts runtime dependencies from the args attribute
// and adds them to the resource. This is needed for lazy loading because DecodeHclBody
// doesn't call DecodeArgs which extracts runtime dependencies.
func (l *Loader) extractRuntimeDependencies(body *hclsyntax.Body, resource modconfig.HclResource, evalCtx *hcl.EvalContext) {
	if body == nil {
		return
	}

	// Check if resource is a QueryProvider (only they have args)
	queryProvider, ok := resource.(resources.QueryProvider)
	if !ok {
		return
	}

	// Look for args attribute
	argsAttr, ok := body.Attributes[schema.AttributeTypeArgs]
	if !ok {
		return
	}

	// Decode args to extract runtime dependencies
	args, runtimeDependencies, diags := pparse.DecodeArgs(argsAttr.AsHCLAttribute(), evalCtx, queryProvider)
	if !diags.HasErrors() {
		// Set args and add runtime dependencies
		queryProvider.SetArgs(args)
		queryProvider.AddRuntimeDependencies(runtimeDependencies)
	}
}

// resolveBaseReference attempts to resolve a `base` attribute reference for a resource.
// This is needed for lazy loading because base references can't be resolved during HCL decode
// without the referenced resource being already loaded.
func (l *Loader) resolveBaseReference(ctx context.Context, body *hclsyntax.Body, resource modconfig.HclResource) {
	// Only nodes and edges support base references in the context of lazy loading
	if body == nil {
		return
	}

	// Check if there's a base attribute
	baseAttr, ok := body.Attributes["base"]
	if !ok {
		return
	}

	// Try to extract the reference from the base attribute expression
	// The expression should be a scope traversal like `node.chaos_cache_check_top`
	expr := baseAttr.Expr
	traversal, ok := expr.(*hclsyntax.ScopeTraversalExpr)
	if !ok || len(traversal.Traversal) < 2 {
		return
	}

	// Extract resource type and name from traversal
	// e.g., node.chaos_cache_check_top -> type="node", name="chaos_cache_check_top"
	resourceType := traversal.Traversal.RootName()
	if len(traversal.Traversal) < 2 {
		return
	}
	attrTraversal, ok := traversal.Traversal[1].(hcl.TraverseAttr)
	if !ok {
		return
	}
	resourceName := attrTraversal.Name

	// Build the full resource name for lookup
	// The mod name should be the same as the current resource's mod
	fullName := fmt.Sprintf("%s.%s.%s", l.mod.ShortName, resourceType, resourceName)

	// Try to load the base resource
	baseResource, err := l.Load(ctx, fullName)
	if err != nil || baseResource == nil {
		return
	}

	// Set the base on the resource based on type
	switch r := resource.(type) {
	case *resources.DashboardNode:
		if baseNode, ok := baseResource.(*resources.DashboardNode); ok {
			r.Base = baseNode
			r.SetBaseProperties()
		}
	case *resources.DashboardEdge:
		if baseEdge, ok := baseResource.(*resources.DashboardEdge); ok {
			r.Base = baseEdge
			r.SetBaseProperties()
		}
	case *resources.DashboardChart:
		if baseChart, ok := baseResource.(*resources.DashboardChart); ok {
			r.Base = baseChart
			r.SetBaseProperties()
		}
	case *resources.DashboardCard:
		if baseCard, ok := baseResource.(*resources.DashboardCard); ok {
			r.Base = baseCard
			r.SetBaseProperties()
		}
	case *resources.DashboardTable:
		if baseTable, ok := baseResource.(*resources.DashboardTable); ok {
			r.Base = baseTable
			r.SetBaseProperties()
		}
	case *resources.DashboardInput:
		if baseInput, ok := baseResource.(*resources.DashboardInput); ok {
			r.Base = baseInput
			r.SetBaseProperties()
		}
	}
}
