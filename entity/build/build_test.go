package build

import (
	"github.com/AmadlaOrg/hery/entity"
	"reflect"
	"testing"
)

func TestMetaFromRemoteWithoutVersion(t *testing.T) {
	entityBuildService := NewEntityBuildService()

	tests := []struct {
		name         string
		entityUri    string
		expectEntity entity.Entity
		expectError  error
	}{
		{
			name: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := entityBuildService.metaFromRemoteWithoutVersion(tt.entityUri)
			if err != nil {
				return
			}

			if !reflect.DeepEqual(entity, tt.expectEntity) {
				t.Errorf("got %v, want %v", entity, tt.expectEntity)
			}
		})
	}
}
