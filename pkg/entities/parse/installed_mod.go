package parse

import (
	"github.com/Masterminds/semver/v3"
	"github.com/turbot/powerpipe/pkg/entities"
)

type InstalledMod struct {
	Mod     *entities.Mod
	Version *semver.Version
}
