package main

import (
    "fmt"
    "os"
    "path/filepath"
    "os/exec"
    "runtime"
    "sync"
    "slices"
    "strings"
)

/*
负责人:Yufeng Gosling(项目创始人)
创建时间:由于项目早期无记录，不详，只有大致的日期2025/2(我创建这个项目时只是个六年级小学生啊!)
电子邮箱:yufeng_gosling_work@outlook.com
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
                         
               如来保佑代码没有bug(虽然我信道，不过应该是同吧？)
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

// 读取并匹配Py文件里面的库，这里调用项目根目录scripts文件夹里面的match_lib.pl这个perl脚本
func match_lib(input chan string, output chan []string) {
    defer wg.Done()
    for py_file_path := range input {
        cmd := exec.Command("perl", "scripts/match_lib.pl", py_file_path)
        out, err := cmd.Output()
        if err != nil {
            if exitError, ok := err.(*exec.ExitError); ok {
                fmt.Printf("Error: %s", exitError.Stderr)
            }
            return
        }
        lib := string(out)
        output <- strings.Split(lib, "\n")
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

    pack_slices_ch := make(chan []string, len(py_file_path))
    for i := 0; i < num_goroutine; i += 1 {
        wg.Add(1)
        go match_lib(path_ch, pack_slices_ch)
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



