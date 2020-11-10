package file

import (
	"gin-template-app/pkg/utils"
	"strings"
)

// GetSubPath get parent directory
func GetSubPath(dirctory string, level int) string {
	for i := 0; i < level; i++ {
		dirctory = utils.Substring(dirctory, 0, strings.LastIndex(dirctory, "/"))
	}
	return dirctory
}
