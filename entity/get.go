package entity

import (
	"fmt"
	"github.com/AmadlaOrg/hery/entity/version"
	"sync"
)

func Get(entities []Entity, dest string) {
	var wg sync.WaitGroup
	wg.Add(len(entities))

	for i, entity := range entities {
		go func(i int, entity Entity) {
			defer wg.Done()
			entities[i].Version = version.Latest(
				fmt.Sprintf("https://%s/%s", entity.Origin, entity.Name),
				fmt.Sprintf("/tmp/repo%d", i))
		}(i, entity)
	}

	wg.Wait()

	for _, entity := range entities {
		fmt.Printf("Module: %s, Version: %s\n", entity.Name, entity.Version)
	}
}
