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

	return l.decodeResourceBlock(ctx, entry, block)
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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

	// Map of block types to factory functions
	factoryFuncs := map[string]func(*hcl.Block, *modconfig.Mod, string) modconfig.HclResource{
		schema.BlockTypeBenchmark: resources.NewBenchmark,
		schema.BlockTypeCard:      resources.NewDashboardCard,
		schema.BlockTypeCategory:  resources.NewDashboardCategory,
		schema.BlockTypeContainer: resources.NewDashboardContainer,
		schema.BlockTypeChart:     resources.NewDashboardChart,
		schema.BlockTypeControl:   resources.NewControl,
		schema.BlockTypeDashboard: resources.NewDashboard,
		schema.BlockTypeDetection: resources.NewDetection,
		schema.BlockTypeEdge:      resources.NewDashboardEdge,
		schema.BlockTypeFlow:      resources.NewDashboardFlow,
		schema.BlockTypeGraph:     resources.NewDashboardGraph,
		schema.BlockTypeHierarchy: resources.NewDashboardHierarchy,
		schema.BlockTypeImage:     resources.NewDashboardImage,
		schema.BlockTypeInput:     resources.NewDashboardInput,
		schema.BlockTypeNode:      resources.NewDashboardNode,
		schema.BlockTypeQuery:     resources.NewQuery,
		schema.BlockTypeTable:     resources.NewDashboardTable,
		schema.BlockTypeText:      resources.NewDashboardText,
		schema.BlockTypeWith:      resources.NewDashboardWith,
	}

	factoryFunc, ok := factoryFuncs[block.Type]
	if !ok {
		return nil, fmt.Errorf("unsupported block type: %s", block.Type)
	}

	return factoryFunc(block, mod, entry.ShortName), nil
}

// decodeNestedBlocks decodes nested blocks (children) for dashboards and containers.
func (l *Loader) decodeNestedBlocks(ctx context.Context, parent modconfig.HclResource, body *hclsyntax.Body, evalCtx *hcl.EvalContext, entry *resourceindex.IndexEntry) error {
	// Check if this resource type supports nested blocks
	parentType := parent.GetBlockType()
	if parentType != schema.BlockTypeDashboard && parentType != schema.BlockTypeContainer {
		return nil
	}

	// Track child index for generating anonymous names
	childCounts := make(map[string]int)

	for _, b := range body.Blocks {
		block := b.AsHCLBlock()

		// Skip non-dashboard child blocks
		if !isDashboardChildType(block.Type) {
			continue
		}

		// Generate name for anonymous blocks
		childCounts[block.Type]++
		shortName := ""
		if len(block.Labels) > 0 {
			shortName = block.Labels[0]
		} else {
			// Anonymous block - generate a name based on parent and index
			shortName = fmt.Sprintf("%s_%d", block.Type, childCounts[block.Type])
		}

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

		// Recursively decode nested blocks in this child
		if childBody, ok := block.Body.(*hclsyntax.Body); ok {
			if err := l.decodeNestedBlocks(ctx, childResource, childBody, evalCtx, entry); err != nil {
				continue
			}
		}

		// Add child to parent based on parent type
		if childItem, ok := childResource.(modconfig.ModTreeItem); ok {
			switch p := parent.(type) {
			case *resources.Dashboard:
				p.AddChild(childItem)
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
		schema.BlockTypeText:
		return true
	default:
		return false
	}
}
