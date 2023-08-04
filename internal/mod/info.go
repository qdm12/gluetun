package mod

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/sys/unix"
)

type state uint8

const (
	unloaded state = iota
	loading
	loaded
	builtin
)

type moduleInfo struct {
	state           state
	dependencyPaths []string
}

var (
	ErrModulesDirectoryNotFound = errors.New("modules directory not found")
)

func getModulesInfo() (modulesInfo map[string]moduleInfo, err error) {
	var utsName unix.Utsname
	err = unix.Uname(&utsName)
	if err != nil {
		return nil, fmt.Errorf("getting unix uname release: %w", err)
	}
	release := unix.ByteSliceToString(utsName.Release[:])
	release = strings.TrimSpace(release)

	modulePaths := []string{
		filepath.Join("/lib/modules", release),
		filepath.Join("/usr/lib/modules", release),
	}

	var modulesPath string
	var found bool
	for _, modulesPath = range modulePaths {
		info, err := os.Stat(modulesPath)
		if err == nil && info.IsDir() {
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("%w: %s are not valid existing directories"+
			"; have you bind mounted the /lib/modules directory?",
			ErrModulesDirectoryNotFound, strings.Join(modulePaths, ", "))
	}

	dependencyFilepath := filepath.Join(modulesPath, "modules.dep")
	dependencyFile, err := os.Open(dependencyFilepath)
	if err != nil {
		return nil, fmt.Errorf("opening dependency file: %w", err)
	}

	modulesInfo = make(map[string]moduleInfo)
	scanner := bufio.NewScanner(dependencyFile)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		path := filepath.Join(modulesPath, strings.TrimSpace(parts[0]))
		dependenciesString := strings.TrimSpace(parts[1])

		if dependenciesString == "" {
			modulesInfo[path] = moduleInfo{}
			continue
		}

		dependencyNames := strings.Split(dependenciesString, " ")
		dependencies := make([]string, len(dependencyNames))
		for i := range dependencyNames {
			dependencies[i] = filepath.Join(modulesPath, dependencyNames[i])
		}
		modulesInfo[path] = moduleInfo{dependencyPaths: dependencies}
	}

	err = scanner.Err()
	if err != nil {
		_ = dependencyFile.Close()
		return nil, fmt.Errorf("modules dependency file scanning: %w", err)
	}

	err = dependencyFile.Close()
	if err != nil {
		return nil, fmt.Errorf("closing dependency file: %w", err)
	}

	err = getBuiltinModules(modulesPath, modulesInfo)
	if err != nil {
		return nil, fmt.Errorf("getting builtin modules: %w", err)
	}

	err = getLoadedModules(modulesInfo)
	if err != nil {
		return nil, fmt.Errorf("getting loaded modules: %w", err)
	}

	return modulesInfo, nil
}

func getBuiltinModules(modulesDirPath string, modulesInfo map[string]moduleInfo) error {
	file, err := os.Open(filepath.Join(modulesDirPath, "modules.builtin"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("opening builtin modules file: %w", err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txt := scanner.Text()
		path := filepath.Join(modulesDirPath, strings.TrimSpace(txt))
		info := modulesInfo[path]
		info.state = builtin
		modulesInfo[path] = info
	}

	err = scanner.Err()
	if err != nil {
		_ = file.Close()
		return fmt.Errorf("scanning builtin modules file: %w", err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("closing builtin modules file: %w", err)
	}
	return nil
}

func getLoadedModules(modulesInfo map[string]moduleInfo) (err error) {
	file, err := os.Open("/proc/modules")
	if err != nil {
		// File cannot be opened, so assume no module is loaded
		return nil //nolint:nilerr
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		name := parts[0]
		path, err := findModulePath(name, modulesInfo)
		if err != nil {
			_ = file.Close()
			return fmt.Errorf("finding module path: %w", err)
		}
		info := modulesInfo[path]
		info.state = loaded
		modulesInfo[path] = info
	}

	err = scanner.Err()
	if err != nil {
		_ = file.Close()
		return fmt.Errorf("scanning modules: %w", err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("closing process modules file: %w", err)
	}

	return nil
}

var (
	ErrModulePathNotFound = errors.New("module path not found")
)

func findModulePath(moduleName string, modulesInfo map[string]moduleInfo) (modulePath string, err error) {
	// Kernel module names can have underscores or hyphens in their names,
	// but only one or the other in one particular name.
	nameHyphensOnly := strings.ReplaceAll(moduleName, "_", "-")
	nameUnderscoresOnly := strings.ReplaceAll(moduleName, "-", "_")

	validModuleExtensions := []string{".ko", ".ko.gz", ".ko.xz", ".ko.zst"}
	const nameVariants = 2
	validFilenames := make(map[string]struct{}, nameVariants*len(validModuleExtensions))
	for _, ext := range validModuleExtensions {
		validFilenames[nameHyphensOnly+ext] = struct{}{}
		validFilenames[nameUnderscoresOnly+ext] = struct{}{}
	}

	for modulePath := range modulesInfo {
		moduleFileName := path.Base(modulePath)
		_, valid := validFilenames[moduleFileName]
		if valid {
			return modulePath, nil
		}
	}

	return "", fmt.Errorf("%w: for %q", ErrModulePathNotFound, moduleName)
}
