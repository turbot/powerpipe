package controldisplay

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/turbot/pipe-fittings/v2/app_specific"
	"github.com/turbot/powerpipe/internal/controlexecute"
)

func init() {
	// Format() reads app_specific.AppVersion, which the real CLI sets during
	// its startup bootstrap (cmdconfig.SetAppSpecificConstants). A plain
	// `go test` never runs that bootstrap, so it must be set here or every
	// TemplateFormatter.Format call nil-derefs.
	app_specific.AppVersion = semver.MustParse("0.0.0-test")
}

// buildMockTree constructs a real *controlexecute.ExecutionTree by hand
// (no HCL parsing, no database) with one "ok" row and one "alarm" row, to
// exercise both the Observation-only and Observation+Finding paths of the
// oscal.json template.
func buildMockTree() *controlexecute.ExecutionTree {
	okControl := &controlexecute.ControlRun{
		FullName:    "aws_compliance.control.s3_bucket_default_encryption_enabled",
		Title:       "S3 buckets should have default encryption enabled",
		Description: "Checks that S3 buckets have default encryption configured.",
		Severity:    "high",
	}
	okControl.Rows = controlexecute.ResultRows{
		{
			Run:      okControl,
			Status:   "ok",
			Reason:   "Bucket has default encryption enabled",
			Resource: "arn:aws:s3:::my-encrypted-bucket",
		},
	}

	alarmControl := &controlexecute.ControlRun{
		FullName:    "aws_compliance.control.s3_bucket_public_access_block",
		Title:       "S3 buckets should block public access",
		Description: "Checks that S3 buckets have public access blocked.",
		Severity:    "critical",
	}
	alarmControl.Rows = controlexecute.ResultRows{
		{
			Run:      alarmControl,
			Status:   "alarm",
			Reason:   "Bucket does not have public access block enabled",
			Resource: "arn:aws:s3:::my-public-bucket",
		},
	}

	return &controlexecute.ExecutionTree{
		ControlRuns: map[string]*controlexecute.ControlRun{
			okControl.FullName:    okControl,
			alarmControl.FullName: alarmControl,
		},
	}
}

// TestOscalFormatter renders the real oscal.json template (via the real
// TemplateFormatter code path) against a hand-built ExecutionTree and writes
// the result to OSCAL_TEST_OUTPUT (if set) for out-of-process schema
// validation against the official NIST OSCAL 1.2.3 schema.
func TestOscalFormatter(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	templateDir := filepath.Join(wd, "templates", "oscal.json")
	if _, err := os.Stat(templateDir); err != nil {
		t.Fatalf("oscal.json template directory not found at %s: %v", templateDir, err)
	}

	formatter, err := NewTemplateFormatter(NewOutputTemplate(templateDir))
	if err != nil {
		t.Fatalf("NewTemplateFormatter failed: %v", err)
	}

	reader, err := formatter.Format(context.Background(), buildMockTree())
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	out, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("failed reading formatted output: %v", err)
	}

	if len(out) == 0 {
		t.Fatal("formatter produced no output")
	}

	if dest := os.Getenv("OSCAL_TEST_OUTPUT"); dest != "" {
		if err := os.WriteFile(dest, out, 0o644); err != nil {
			t.Fatalf("failed writing output to %s: %v", dest, err)
		}
		t.Logf("wrote rendered OSCAL output to %s (%d bytes)", dest, len(out))
	} else {
		t.Logf("rendered output (%d bytes):\n%s", len(out), string(out))
	}
}

// TestOscalFormatterAllOk exercises the all-"ok" (zero-finding) path: the
// "findings" key must be omitted entirely, not emitted as [], since the
// OSCAL schema's minItems:1 forbids an empty array when the key is present.
func TestOscalFormatterAllOk(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	templateDir := filepath.Join(wd, "templates", "oscal.json")

	formatter, err := NewTemplateFormatter(NewOutputTemplate(templateDir))
	if err != nil {
		t.Fatalf("NewTemplateFormatter failed: %v", err)
	}

	okControl := &controlexecute.ControlRun{
		FullName:    "aws_compliance.control.s3_bucket_default_encryption_enabled",
		Title:       "S3 buckets should have default encryption enabled",
		Description: "Checks that S3 buckets have default encryption configured.",
		Severity:    "high",
	}
	okControl.Rows = controlexecute.ResultRows{
		{Run: okControl, Status: "ok", Reason: "ok", Resource: "arn:aws:s3:::b"},
	}
	tree := &controlexecute.ExecutionTree{
		ControlRuns: map[string]*controlexecute.ControlRun{okControl.FullName: okControl},
	}

	reader, err := formatter.Format(context.Background(), tree)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}
	out, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("failed reading formatted output: %v", err)
	}

	if dest := os.Getenv("OSCAL_TEST_OUTPUT_ALLOK"); dest != "" {
		if err := os.WriteFile(dest, out, 0o644); err != nil {
			t.Fatalf("failed writing output to %s: %v", dest, err)
		}
		t.Logf("wrote all-ok output to %s (%d bytes)", dest, len(out))
	}
}
