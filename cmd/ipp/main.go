package main

import (
    "fmt"
    "os"
    "path/filepath"
    "regexp"
    "os/exec"
    "runtime"
    "sync"
    "slices"
)

/*
负责人:Yufeng Gosling
*/

/*
                             _ooOoo_
                            o8888888o
                            88" . "88
                            (| -_- |)
                            O\  =  /O
                         ____/`---'\____
                       .'  \\|     |//  `.
                      /  \\|||  :  |||//  \
                     /  _||||| -:- |||||-  \
                     |   | \\\  -  /// |   |
                     | \_|  ''\---/''  |   |
                     \  .-\__  `-`  ___/-. /
                   ___`. .'  /--.--\  `. . __
                ."" '<  `.___\_<|>_/___.'  >'"".
               | | :  `- \`.;`\ _ /`;.`/ - ` : | |
               \  \ `-.   \_ __\ /__ _/   .-` /  /
          ======`-.____`-.___\_____/___.-`____.-'======
                             `=---='
                         
                      如来保佑代码没有bug
*/

var wg sync.WaitGroup

// 获取所有Py文件
func get_py_file(dir string) ([]string, error) {
    py_file_lis := []string{}
    err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if len(path) > 3 && path[len(path) - 3:] == ".py" {
            py_file_lis = append(py_file_lis, path)
            return nil
        }
    return err})
    
    if err != nil {
	    fmt.Printf("Walk dir %s fail. \nError: %v\n", dir, err)
    }

    return py_file_lis, err
}

// 读取Py文件
func read_file(input chan string, output chan string) error {
    defer wg.Done()
    for file_name := range input {
        file_data, err := os.ReadFile(file_name)
        if err != nil {
            fmt.Printf("Unable to read file %s. \nError: %v\n", file_name, err)
            return err
        }
        output <- string(file_data)
    }
    return nil
}

// 匹配代码的库
func match_lib(file_data_ch chan string, output chan []string, re *regexp.Regexp) {
    defer wg.Done()
    for file_data := range file_data_ch {
        output <- re.FindAllString(file_data, -1)
    }
}

// 安装库
func install_lib(package_name_ch chan string) {
    defer wg.Done()
    for package_name := range package_name_ch {
        fmt.Printf("Installing %s\n", package_name)
        cmd := exec.Command("pip", "install", package_name)
        cmd.Run()
    }
}
    
func main() {
    num_goroutine := runtime.GOMAXPROCS(0) * 2
    dir, err := os.Getwd()
    if err != nil {
        fmt.Printf("Unable to retrieve the current directory. \nError: %v\n", err)
        return
    }

    py_file_path, err := get_py_file(dir)
    if err != nil {
        return
    }

    path_ch := make(chan string, len(py_file_path))
    for _, py_file := range py_file_path {
        path_ch <- py_file
    }
    close(path_ch)
    
    file_data := make(chan string, len(path_ch))
    for i := 0; i < num_goroutine; i += 1 {
        wg.Add(1)
        go read_file(path_ch, file_data)
    }
    wg.Wait()
    close(file_data)

    reg := `\b(?:import\s+(\w+)|from\s+(\w+)\s+import\b)`
    re := regexp.MustCompile(reg)
    pack_slices_ch := make(chan []string, len(file_data))
    for i := 0; i < num_goroutine; i += 1 {
        wg.Add(1)
        go match_lib(file_data, pack_slices_ch, re)
    }
    wg.Wait()
    close(pack_slices_ch)

    pack_lis := []string{}
    for pack_slices := range pack_slices_ch {
        for _, pack := range pack_slices {
            if !slices.Contains(pack_lis, pack) {
                pack_lis = append(pack_lis, pack)
            }
        }
    }

    install_pack := make(chan string, len(pack_lis))
    for _, pack := range pack_lis {
        install_pack <- pack
    }
    close(install_pack)

    for i := 0; i < num_goroutine; i += 1 {
        wg.Add(1)
        go install_lib(install_pack)
    }
    wg.Wait()
}



