package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func main() {
	fmt.Println(os.Getppid())
	// 参数解析
	path := flag.String("path", "", "安装包下载地址")
	pkgname := flag.String("pkgname", "", "安装包名")
	version := flag.String("version", "", "当前版本号")
	temp := flag.String("temp", "", "下载目录")

	flag.Parse()

	if *path == "" || *pkgname == "" || *version == "" || *temp == "" {
		log.Fatal("parameter is null")
		return
	}

	// 判断下载目录是否存在
	c, err := CreateDir(*temp)
	if !c {
		log.Fatal("Create dir failed: ", err)
		return
	}

	// 获取最新发布版本号
	errCode, errMsg := Download(*path, *temp, "database.txt")
	if errCode != Success {
		log.Fatal("download file failed: ", errMsg)
		return
	}

	database := filepath.Join(*temp, "database.txt")
	name, md5 := FileParse(database, *pkgname)
	if name == "" || md5 == "" {
		log.Fatal("database file illegal")
		return
	}

	// 版本比较
	names := strings.SplitN(name, "-", 2)
	index := strings.LastIndex(names[1], ".")
	newVersion := names[1][:index]
	// newVersions := strings.Split(newVersion, "-")
	// versions := strings.Split(*version, "-")
	result := strings.Compare(newVersion, *version)
	if result <= 0 {
		log.Fatal("No new version found")
		return
	}

	// 下载最新发布版本
	errCode, errMsg = Download(*path, *temp, name)
	if errCode != Success {
		log.Fatal("download release package failed: ", errMsg)
		return
	}

	installName := filepath.Join(*temp, name)
	b := CheckFileMd5(installName, md5)
	if !b {
		log.Fatal("check file md5 failed")
		return
	}

	// 获取父进程的 ID
	ppid := uint32(syscall.Getppid())
	fmt.Println("parent process ID:", ppid)

	cmd := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", ppid))
	err = cmd.Run()
	if err != nil {
		fmt.Println("发送退出信号失败:", err)
		return
	}

	time.Sleep(5 * time.Second)

	// 执行静默安装
	cmd = exec.Command(installName, "/silent")
	err = cmd.Run()
	if err != nil {
		fmt.Println("静默安装执行失败:", err)
		return
	}
}
