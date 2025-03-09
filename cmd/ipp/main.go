package main

import (
    "fmt"
    "os"
    "path/filepath"
    "regexp"
    "os/exec"
)

func main() {
    dir, get_wd_err := os.Getwd()
    if get_wd_err != nil {
        fmt.Printf("Can't get current working directory %v\n", get_wd_err)
        os.Exit(1)
    }

    py_file := []string{}
    err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            fmt.Printf("Access directory %s error %v\n", path, err)
            return err
        }

        if len(path) > 3 && path[len(path)-3:] == ".py" {
            py_file = append(py_file, path)
        }
        return nil
    })

    if err != nil {
        fmt.Printf("Go through the directory fail %v\n", err)
        os.Exit(1)
    }

    py_package_set := make(map[string]struct{})
    py_package_regex := `\b(?:import\s+(\w+)|from\s+(\w+)\s+import\b)`
    py_package_re := regexp.MustCompile(py_package_regex)

    for _, py_file_path := range py_file {
        py_file_data, read_file_err := os.ReadFile(py_file_path)
        if read_file_err != nil {
            fmt.Printf("Can't read file: %s error:%v\n", py_file_path, read_file_err)
            continue
        }

        matches := py_package_re.FindAllStringSubmatch(string(py_file_data), -1)
        for _, match := range matches {
            if match[1] != "" {
                py_package_set[match[1]] = struct{}{}
            }
            if match[2] != "" {
                py_package_set[match[2]] = struct{}{}
            }
        }
    }

    for py_package := range py_package_set {
        cmd := exec.Command("pip", "install", py_package)
        fmt.Printf("Installing %s\n", py_package)
        if run_command_err := cmd.Run(); run_command_err != nil {
            fmt.Printf("Install %s fail: %v\n", py_package, run_command_err)
        } else {
            fmt.Printf("Install Success: %s\n", py_package)
        }
    }
}