package entity

import (
	"fmt"
	"net/http"
)

func Download(uri string) error {
	if uri == "" {
		return fmt.Errorf("uri is empty")
	}

	// TODO: How does Go do it? SSH, HTTPS, some other way?
	if _, err := http.Get(fmt.Sprintf("https://%s", uri)); err != nil {
		return err
	}
	return nil
}
