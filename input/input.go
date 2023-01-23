package input

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func ForDay(day int) string {
	filePath := path.Join("resources", fmt.Sprintf("Day%d.txt", day))
	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("%s: error %s", filePath, err))
	}
	return strings.TrimSpace(string(content))
}
