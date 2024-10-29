package snapshot

import (
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/utils"
	powerpipe2 "github.com/turbot/powerpipe/internal/resources"
	"reflect"
	"testing"
)

func TestGetAsSnapshotPropertyMap(t *testing.T) {
	type args struct {
		item interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]any
	}{
		{name: "card",
			args: args{
				item: powerpipe2.DashboardChart{
					QueryProviderImpl: powerpipe2.QueryProviderImpl{
						RuntimeDependencyProviderImpl: powerpipe2.RuntimeDependencyProviderImpl{
							ModTreeItemImpl: modconfig.ModTreeItemImpl{
								HclResourceImpl: modconfig.HclResourceImpl{
									FullName:        "mod1.card.card1",
									ShortName:       "card1",
									UnqualifiedName: "card.card1",
									Description:     utils.ToStringPointer("a card"),
								},
							},
						},
						SQL: utils.ToStringPointer("select 1"),
					},
					Axes: &powerpipe2.DashboardChartAxes{
						X: &powerpipe2.DashboardChartAxesX{
							Title: &powerpipe2.DashboardChartAxisTitle{
								Value: utils.ToStringPointer("x axis"),
							},
							Min: utils.ToIntegerPointer(0),
							Max: utils.ToIntegerPointer(1000),
						},
						Y: &powerpipe2.DashboardChartAxesY{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAsSnapshotPropertyMap(tt.args.item)
			if err != nil {
				t.Fail()
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAsSnapshotPropertyMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
