package controldisplay

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/turbot/pipe-fittings/filepaths"
)

//go:embed templates/* detection_templates/*
var builtinTemplateFS embed.FS

type TemplateVersionFile struct {
	Version string `json:"version"`
}

// EnsureControlTemplates scans the '$STEAMPIPE_INSTALL_DIR/check/templates' directory and
// copies over the templates defined in the 'templates' package if needed.
//
// The name of the folder in the 'templates' package is used to identify
// templates in '$STEAMPIPE_INSTALL_DIR/check/templates' - where it is expected
// that a directory with the same name will exist.
//
// We re-write the templates, when there is a higher template version
// available in the 'templates' package.
func EnsureControlTemplates() error {
	slog.Debug("ensuring check export/output templates")
	dirs, err := fs.ReadDir(builtinTemplateFS, "templates")
	if err != nil {
		return err
	}
	for _, d := range dirs {
		targetDirectory := filepath.Join(filepaths.EnsureControlTemplateDir(), d.Name())
		currentVersionsFilePath := filepath.Join(targetDirectory, "version.json")
		embeddedVersionsFilePath := filepath.Join("templates", d.Name(), "version.json")

		// check if version in version.json matches with embedded template version
		if getCurrentTemplateVersion(currentVersionsFilePath) != getEmbeddedTemplateVersion(embeddedVersionsFilePath) {
			slog.Debug("versions do not match - copying updated template", "dir", d)
			if err := writeControlTemplate(d.Name(), targetDirectory); err != nil {
				slog.Debug("error copying template", "error", err)
				return err
			}
		}
	}
	return nil
}

// EnsureDetectionTemplates scans the '$STEAMPIPE_INSTALL_DIR/check/detection_templates' directory and
// copies over the templates defined in the 'templates' package if needed.
//
// The name of the folder in the 'templates' package is used to identify
// templates in '$STEAMPIPE_INSTALL_DIR/check/detection_templates' - where it is expected
// that a directory with the same name will exist.
//
// We re-write the templates, when there is a higher template version
// available in the 'templates' package.
func EnsureDetectionTemplates() error {
	slog.Debug("ensuring detection export/output templates")
	dirs, err := fs.ReadDir(builtinTemplateFS, "detection_templates")
	if err != nil {
		return err
	}
	for _, d := range dirs {
		targetDirectory := filepath.Join(filepaths.EnsureDetectionTemplateDir(), d.Name())
		currentVersionsFilePath := filepath.Join(targetDirectory, "version.json")
		embeddedVersionsFilePath := filepath.Join("detection_templates", d.Name(), "version.json")

		// check if version in version.json matches with embedded template version
		if getCurrentTemplateVersion(currentVersionsFilePath) != getEmbeddedTemplateVersion(embeddedVersionsFilePath) {
			slog.Debug("versions do not match - copying updated template", "dir", d)
			if err := writeDetectionTemplate(d.Name(), targetDirectory); err != nil {
				slog.Debug("error copying template", "error", err)
				return err
			}
		}
	}
	return nil
}

func getCurrentTemplateVersion(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Debug("template version file does not exist - install the new template")
		} else {
			slog.Debug("error reading current version file - installing the new template")
		}
		return ""
	}
	var ver TemplateVersionFile
	err = json.Unmarshal(data, &ver)
	if err != nil {
		slog.Debug("error while unmarshalling current version.json file", "error", err)
		return ""
	}
	return ver.Version
}

func getEmbeddedTemplateVersion(path string) string {
	data, err := fs.ReadFile(builtinTemplateFS, path)
	if err != nil {
		slog.Debug("error reading embedded version file - installing the new template")
		return ""
	}
	var ver TemplateVersionFile
	err = json.Unmarshal(data, &ver)
	if err != nil {
		slog.Debug("error while unmarshalling json", "error", err)
		return ""
	}
	return ver.Version
}

func writeControlTemplate(path string, target string) error {
	err := os.MkdirAll(target, 0744)
	if err != nil {
		return err
	}

	entries, err := fs.ReadDir(builtinTemplateFS, filepath.Join("templates", path))
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		sourceInEmbedFs := filepath.Join("templates", path, entry.Name())
		bytes, err := fs.ReadFile(builtinTemplateFS, sourceInEmbedFs)
		if err != nil {
			return err
		}

		//nolint: gosec // this file is safe to be read by all users
		err = os.WriteFile(filepath.Join(target, entry.Name()), bytes, 0744)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeDetectionTemplate(path string, target string) error {
	err := os.MkdirAll(target, 0744)
	if err != nil {
		return err
	}
	detectionEntries, err := fs.ReadDir(builtinTemplateFS, filepath.Join("detection_templates", path))
	if err != nil {
		return err
	}
	for _, entry := range detectionEntries {
		if entry.IsDir() {
			continue
		}
		sourceInEmbedFs := filepath.Join("detection_templates", path, entry.Name())
		bytes, err := fs.ReadFile(builtinTemplateFS, sourceInEmbedFs)
		if err != nil {
			return err
		}

		//nolint: gosec // this file is safe to be read by all users
		err = os.WriteFile(filepath.Join(target, entry.Name()), bytes, 0744)
		if err != nil {
			return err
		}
	}

	return nil
}
