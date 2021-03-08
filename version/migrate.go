package version

import "fmt"

func MigrateFrom(prevVersion string) bool {
	if prevVersion != VERSION {
		fmt.Printf("Updated from version %s to %s\n", prevVersion, VERSION)
	}
	return false
}
