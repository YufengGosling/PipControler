package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
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
func getPyFile(dir string) ([]string, error) {
	fmt.Println("Finding Python Files ...")
	pyFileLis := []string{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if len(path) > 3 && path[len(path)-3:] == ".py" {
			pyFileLis = append(pyFileLis, path)
			return nil
		}
		return err
	})

	if err != nil {
		fmt.Printf("Walk dir %s fail. \nError: %v\n", dir, err)
	}

	return pyFileLis, err
	fmt.Println("Finded Python File.")
}

// 读取并匹配Py文件里面的库，这里调用项目根目录scripts文件夹里面的match_lib.pl这个perl脚本
func matchLib(input chan string, output chan []string) {
	defer wg.Done()
	for pyFilePath := range input {
		cmd := exec.Command("perl", "scripts/match_lib.pl", pyFilePath)
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
func installLib(packageNameCh chan string) {
	defer wg.Done()
	for package_name := range package_name_ch {
		fmt.Printf("Installing %s\n", package_name)
		cmd := exec.Command("pip", "install", package_name)
		cmd.Run()
		fmt.Printf("Installed %s Done.\n", package_name)
	}
}

func main() {
	numGoroutine := runtime.GOMAXPROCS(0) * 2
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Unable to retrieve the current directory. \nError: %v\n", err)
		return
	}

	pyFilePath, err := getPyFile(dir)
	if err != nil {
		return
	}

	pathCh := make(chan string, len(pyFilePath))
	for _, pyFile := range pyFilePath {
		pathCh <- pyFile
	}
	close(pathCh)

	packSlicesCh := make(chan []string, len(pyFilePath))
	for i := 0; i < numGoroutine; i += 1 {
		wg.Add(1)
		go matchLib(pathCh, packSlicesCh)
	}
	wg.Wait()
	close(packSlicesCh)

	packLis := []string{}
	for packSlices := range packSlicesCh {
		for _, pack := range packSlices {
			if !slices.Contains(packLis, pack) {
				packLis = append(packLis, pack)
			}
		}
	}

	installPack := make(chan string, len(packLis))
	for _, pack := range packLis {
		installPack <- pack
	}
	close(installPack)

	for i := 0; i < numGoroutine; i += 1 {
		wg.Add(1)
		go installLib(installPack)
	}
	wg.Wait()

	fmt.Println("OK")
}
