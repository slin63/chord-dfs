// Helper functions for file system operations
package filesys

import (
	"io/ioutil"
	"log"

	"../spec"
)

// Destructive file writer
func Write(filename string, data []byte) int {
	err := ioutil.WriteFile(spec.Filedir+filename, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return len(data)
}
