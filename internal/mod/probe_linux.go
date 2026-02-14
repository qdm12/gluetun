package mod

import (
	"fmt"
)

// Probe loads the given kernel module and its dependencies.
func Probe(moduleName string) error {
	modulesInfo, err := getModulesInfo()
	if err != nil {
		return fmt.Errorf("getting modules information: %w", err)
	}

	modulePath, err := findModulePath(moduleName, modulesInfo)
	if err != nil {
		return fmt.Errorf("finding module path: %w", err)
	}

	info := modulesInfo[modulePath]
	if info.state == builtin || info.state == loaded {
		return nil
	}

	info.state = loading
	for _, dependencyModulePath := range info.dependencyPaths {
		err = initDependencies(dependencyModulePath, modulesInfo)
		if err != nil {
			return fmt.Errorf("init dependencies: %w", err)
		}
	}

	err = initModule(modulePath)
	if err != nil {
		return fmt.Errorf("init module: %w", err)
	}
	return nil
}
