package conf

import "fmt"

func defaultConfPath(appName string) string {
	if appName == "" {
		return "."
	}
	return fmt.Sprintf("/etc/%s/", appName)
}
