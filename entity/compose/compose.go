package compose

import (
	"fmt"
	"os"

	"github.com/AmadlaOrg/hery/entity"
	"github.com/AmadlaOrg/hery/entity/merge"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/goccy/go-yaml"
)

// IComposer is an interface for composing entities.
type IComposer interface {
	ComposeEntity(entityArg string, printToScreen bool) error
}

// SComposer struct implements the IComposer interface.
type SComposer struct {
	Storage storage.IStorage
	Entity  entity.IEntity
}

// ComposeEntity reads all .hery files for an entity, resolves _parent chains
// via deep merge, and outputs the composed result.
func (s *SComposer) ComposeEntity(entityArg string, printToScreen bool) error {
	storagePaths, err := s.Storage.Paths()
	if err != nil {
		return fmt.Errorf("failed to get storage paths: %w", err)
	}

	// Find the entity directory
	entityDir := storagePaths.Entities
	documents, err := s.Entity.ReadAll(entityDir)
	if err != nil {
		return fmt.Errorf("failed to read entities: %w", err)
	}

	if len(documents) == 0 {
		return fmt.Errorf("no entities found")
	}

	// Filter documents matching the entity arg (by _type)
	var matched []map[string]any
	for _, doc := range documents {
		if typeVal, ok := doc["_type"].(string); ok {
			if typeVal == entityArg {
				matched = append(matched, doc)
			}
		}
	}

	if len(matched) == 0 {
		matched = documents
	}

	// Resolve _parent chains for each document
	parentLookup := func(selfURI string) (map[string]any, error) {
		for _, doc := range documents {
			if selfVal, ok := doc["_self"].(string); ok && selfVal == selfURI {
				return doc, nil
			}
		}
		return nil, fmt.Errorf("parent entity not found: %s", selfURI)
	}

	var resolved []map[string]any
	for _, doc := range matched {
		r, resolveErr := merge.ResolveParentChain(doc, parentLookup)
		if resolveErr != nil {
			return fmt.Errorf("failed to resolve _parent chain: %w", resolveErr)
		}
		resolved = append(resolved, r)
	}

	// Merge all resolved documents into one
	var composed map[string]any
	for _, doc := range resolved {
		if composed == nil {
			composed = doc
		} else {
			composed = merge.DeepMerge(composed, doc)
		}
	}

	output, err := yaml.Marshal(composed)
	if err != nil {
		return fmt.Errorf("failed to marshal composed entity: %w", err)
	}

	if printToScreen {
		fmt.Println(string(output))
	} else {
		if writeErr := os.WriteFile("composed.hery", output, 0644); writeErr != nil {
			return fmt.Errorf("failed to write composed file: %w", writeErr)
		}
	}

	return nil
}
