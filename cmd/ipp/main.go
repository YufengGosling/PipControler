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
    "log"
)

/*
负责人:Yufeng Gosling(项目创始人)
创建时间:由于项目早期无记录，不详，只有大致的日期2025/2(我创建这个项目时只是个六年级小学生啊!不可思议吧)
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

    如来保佑代码没有bug(虽然我信道，不过应该是同行吧？我觉得我该放个菩提祖师)
*/

/*
要是在main分支下的崩了，你把代码拍我脸上
记住
*/

var wg sync.WaitGroup

// 获取所有Py文件
func getPyFile(dir string) map[string]int {
	fmt.Println("Finding Python Files ...")
	pyFileMap := make(map[string]int)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if len(path) > 3 && path[len(path)-3:] == ".py" {
            fileTimeStamp, err := info.ModTime.Unix()
            if err != nil {
                log.Fatalf("Cannot get file stamp: %s\nError: %v", path, err)
            }
            pyFileMap[path] = fileTimeStamp
			return nil
		}
		return
	})

	if err != nil {
		log.Fatalf("Walk dir %s fail. \nError: %v\n", dir, err)
	}

    fmt.Println("Finded Python File.")
	return pyFileMap
}

func selectFile(input map[string]int) string[] {
    path := ".pipcontroler/fileTable"
    pyFilePath := []string
    fileTable := map[string]int
    data, err := os.ReadFile(path)
    if err != nil {
        if err == os.ErrNotExist {
            os.Create(path)
            if err != nil {
                log.Fatalf("Cannot create fileTable. \nError: %v", err)
            }
        } else {
            log.Fatalf("Connot open fileTable. \nError: %v", err)
        }
    }

    reader := strings.NewReader(string(data))
    scanner := bufio.NewScanner(reader)

    for scanner.Scan() {
        line := scanner.Text()
        fileInfo := strings.Split(line, " ")
        fileTable[fileInfo[0]] = fileInfo[1]
    }

    for fileName, fileTimeStamp := range input {
        oldFileTimeStamp, e := fileTable[fileName]
        if !e || oldFileTimeStamp != fileTimeStamp {
            pyFilePath = append(pyFilePath, fileName)
            fileTable[fileName] = fileTimeStamp
        }
    }

    for fileName, fileTimeStamp := range fileTable {
        _, e := input[fileName]
        if !e {
            delete(fileTable, fileName)
        }
    }

    rewrite := ""
    for fileName, fileTimeStamp := range fileTable {
        rewrite = rewrite + fmt.Sprintf("%s %s\n", fileName, fileTimeStamp)
    }
    fileTableFile, err := os.OpenFile(path, os.WRONLY)
    if err != nil {
        log.Fatalf("Cannot open file table file. \nError: %v", err)
    }
    _, err := fileTableFile.WriteString(rewrite)
    if err != nil {
        log.Fatalf("Connot rewrite file table. \nError: %v", err)
    }
}
            

// 读取并匹配Py文件里面的库，这里调用项目根目录scripts文件夹里面的match_lib.pl这个perl脚本
func matchLib(input chan string, output chan []string) {
	defer wg.Done()
	for pyFilePath := range input {
		cmd := exec.Command("perl", "scripts/match_lib.pl", pyFilePath)
		out, err := cmd.Output()
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				log.Printf("Error: %s", exitError.Stderr)
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
	for packageName := range packageNameCh {
		fmt.Printf("Installing %s\n", packageName)
		cmd := exec.Command("pip", "install", packageName)
		cmd.Run()
        fmt.Printf("Installed %s Done.\n", packageName)
	}
}

func main() {
    // 获取goroutine数量，当前工作目录
	numGoroutine := runtime.GOMAXPROCS(0) * 2
	dir, err := os.Getwd()
	if err != nil {
		log.Printf("Unable to retrieve the current directory. \nError: %v\n", err)
		return
	}

    // 获取python文件选择并将其装入通道
	pyFilePathMap := getPyFile(dir)
    pyFilePath := selectFile(pyFilePathMap)
    pathCh := make(chan string, len(pyFilePath))
	for _, pyFile := range pyFilePath {
		pathCh <- pyFile
	}
	close(pathCh)

    // 匹配库，然后从packSlicesCh中取出通道切片中的库名
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

    // 把库名装载进installPack，然后安装库
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

    // 结束
	fmt.Println("OK")
}
