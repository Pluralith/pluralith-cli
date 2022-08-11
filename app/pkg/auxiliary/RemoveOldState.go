package auxiliary

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func RemoveOldState() error {
	functionName := "RemoveOldState"

	for _, file := range []string{"pluralith.plan", "pluralith.state.json", "pluralith.state.hashed"} {
		oldFilePath := filepath.Join(StateInstance.WorkingPath, ".pluralith", file)

		_, existErr := os.Stat(oldFilePath)       // Check if old state exists
		if !errors.Is(existErr, os.ErrNotExist) { // If it exists -> delete
			if removeErr := os.Remove(oldFilePath); removeErr != nil {
				return fmt.Errorf("deleting old plan failed -> %v: %w", functionName, removeErr)
			}
		}
	}

	return nil
}
