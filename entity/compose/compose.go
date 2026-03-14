package compose

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AmadlaOrg/hery/entity"
	"github.com/AmadlaOrg/hery/entity/merge"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/goccy/go-yaml"
)

// Composer is an interface for composing entities.
type Composer interface {
	ComposeEntity(entityArg string, printToScreen bool) error
}

// composer struct implements the Composer interface.
type composer struct {
	Storage storage.Storage
	Entity  entity.Service
}

// ComposeEntity reads all .hery files for an entity, resolves _extends chains
// via deep merge, and outputs the composed result.
func (s *composer) ComposeEntity(entityArg string, printToScreen bool) error {
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

	// Resolve _extends chains by reading from the extended entity's directory in the entity cache.
	// _extends is a git path (e.g., github.com/some-org/base-configs/webserver) that maps
	// to a directory under the entities cache.
	extendsLookup := func(extendsURI string) (map[string]any, error) {
		extendsDir := filepath.Join(storagePaths.Entities, extendsURI)
		extendsDocs, readErr := s.Entity.ReadAll(extendsDir)
		if readErr != nil {
			return nil, fmt.Errorf("extended entity not found at %s: %w", extendsURI, readErr)
		}
		if len(extendsDocs) == 0 {
			return nil, fmt.Errorf("no .hery files found for extended entity: %s", extendsURI)
		}
		// Merge all documents from the extended entity directory into one
		result := extendsDocs[0]
		for _, doc := range extendsDocs[1:] {
			result = merge.DeepMerge(result, doc)
		}
		return result, nil
	}

	var resolved []map[string]any
	for _, doc := range matched {
		r, resolveErr := merge.ResolveExtendsChain(doc, extendsLookup)
		if resolveErr != nil {
			return fmt.Errorf("failed to resolve _extends chain: %w", resolveErr)
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
