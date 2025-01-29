package controldisplay

import (
	"fmt"

	"github.com/turbot/go-kit/files"
	"github.com/turbot/pipe-fittings/v2/constants"
	"github.com/turbot/pipe-fittings/v2/export"
	"github.com/turbot/pipe-fittings/v2/filepaths"
	"github.com/turbot/pipe-fittings/v2/modconfig"
	"github.com/turbot/powerpipe/internal/resources"
)

type FormatResolver struct {
	formatterByName map[string]Formatter
	// array of unique formatters used for export
	exportFormatters []Formatter
}

func NewFormatResolver(target modconfig.ModTreeItem) (*FormatResolver, error) {
	// TACTICAL
	// if the target is a detection or detection benchmark use a separate resolver
	var detection bool
	switch target.(type) {
	case *resources.Detection, *resources.DetectionBenchmark:
		detection = true
	default:
	}

	templates, err := loadAvailableTemplates(detection)
	if err != nil {
		return nil, err
	}

	formatters := []Formatter{
		&NullFormatter{},
		&TextFormatter{},
		&SnapshotFormatter{},
	}

	res := &FormatResolver{
		formatterByName: make(map[string]Formatter),
	}

	for _, f := range formatters {
		if err := res.registerFormatter(f); err != nil {
			return nil, err
		}
	}
	for _, t := range templates {
		f, err := NewTemplateFormatter(t)
		if err != nil {
			return nil, err
		}

		if err := res.registerFormatter(f); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (r *FormatResolver) GetFormatter(arg string) (Formatter, error) {
	if formatter, found := r.formatterByName[arg]; found {
		return formatter, nil
	}

	return nil, fmt.Errorf(" invalid output format: '%s'", arg)
}

func (r *FormatResolver) registerFormatter(f Formatter) error {
	name := f.Name()

	if _, ok := r.formatterByName[name]; ok {
		return fmt.Errorf("failed to register output formatter - duplicate format name %s", name)
	}
	r.formatterByName[name] = f
	// if the formatter has an alias, also register by alias
	if alias := f.Alias(); alias != "" {
		if _, ok := r.formatterByName[alias]; ok {
			return fmt.Errorf("failed to register output formatter - duplicate format name %s", alias)
		}
		r.formatterByName[alias] = f
	}
	// add to exportFormatters list (exclude 'None')
	if f.Name() != constants.OutputFormatNone {
		r.exportFormatters = append(r.exportFormatters, f)
	}
	return nil
}

func (r *FormatResolver) controlExporters() []export.Exporter {
	res := make([]export.Exporter, len(r.exportFormatters))
	for i, formatter := range r.exportFormatters {
		res[i] = NewControlExporter(formatter)

	}
	return res
}

func loadAvailableTemplates(detection bool) ([]*OutputTemplate, error) {
	templateDir := filepaths.EnsureControlTemplateDir()
	if detection {
		templateDir = filepaths.EnsureDetectionTemplateDir()
	}
	templateDirectories, err := files.ListFiles(templateDir, &files.ListOptions{
		Flags:   files.DirectoriesFlat | files.NotEmpty,
		Exclude: []string{"./.*"},
	})
	if err != nil {
		return nil, err
	}

	templates := []*OutputTemplate{}
	for _, templateDirectory := range templateDirectories {
		templates = append(templates, NewOutputTemplate(templateDirectory))
	}
	return templates, nil
}
