package mod

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var errBuiltinModuleNotFound = errors.New("builtin module not found")

func checkModulesBuiltin(modulesPath, moduleName string) error {
	f, err := os.Open(filepath.Join(modulesPath, "modules.builtin"))
	if err != nil {
		return err
	}
	defer f.Close()

	moduleName = strings.TrimSuffix(moduleName, ".ko")

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSuffix(line, ".ko")
		if strings.HasSuffix(line, "/"+moduleName) {
			return nil
		}
	}

	return fmt.Errorf("%w: %s", errBuiltinModuleNotFound, moduleName)
}
