package common

import (
	"bufio"
	"fmt"
	"io/ioutil" //io 工具包
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/**
 * 判断文件或目录是否存在  存在返回 true 不存在返回false
 */
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// LogFile 写入文件，filename：xxxx.txt
func LogFile(filename string, content string) {
	dirname, err := os.Getwd() //当前程序目录
	if err != nil {
		log.Fatal(err)
	}
	var logFilePath string
	if filename == "" {
		logFilePath = filepath.Join(dirname, "logs", fmt.Sprintf("log-%v.txt", time.Now().Format("20060102")))
	} else {
		// 判断 字符串str 是否拥有该后缀
		if !strings.HasSuffix(filename, `.log`) {
			filename = filename + ".txt"
		}
		if !strings.HasSuffix(filename, `.txt`) {
			filename = filename + ".txt"
		}
		logFilePath = filepath.Join(dirname, "logs", filename)
	}

	fmt.Println(logFilePath)
	file, _ := getFile(logFilePath)
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString(content + "\r\n")
	writer.Flush() //Flush将缓存的文件真正写入到文件中

}

func WriteFile(filename string, data []byte) {

	if CheckFileIsExist(filename) { //如果文件存在
		err1 := os.Remove(filename)
		check(err1)
		fmt.Printf("文件存在,已删除：%s", filename)
	} else {
		//fmt.Println("文件不存在")
	}

	err2 := ioutil.WriteFile(filename, data, 0666) //写入文件(字节数组)
	check(err2)
}

// 常用的 flag 文件处理参数：
// O_RDONLY：只读模式打开文件；
// O_WRONLY：只写模式打开文件；
// O_RDWR：读写模式打开文件；
// O_APPEND：写操作时将数据附加到文件尾部（追加）；
// O_CREATE：如果不存在将创建一个新文件；
// O_EXCL：和 O_CREATE 配合使用，文件必须不存在，否则返回一个错误；
// O_SYNC：当进行一系列写操作时，每次都要等待上次的 I/O 操作完成再进行；
// O_TRUNC：如果可能，在打开时清空文件。

//向指定文件追加内容
func FileAppend(filepath string, content string) {

	file, _ := getFile(filepath)
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString(content + "\r\n")
	writer.Flush() //Flush将缓存的文件真正写入到文件中
}

//向指定文件追加内容
func FileAppendMap(filepath string, contents map[string][]string) {

	file, _ := getFile(filepath)

	defer file.Close()

	writer := bufio.NewWriter(file)
	for key, value := range contents { //取map中的值
		str := fmt.Sprintf("%s:%s", key, strings.Join(value, ","))
		writer.WriteString(str + "\r\n")
	}
	writer.Flush() //Flush将缓存的文件真正写入到文件中
}

// 获取文件流，没有文件创建空文件
func getFile(filepath string) (file *os.File, err error) {

	// var file *os.File
	// var err error

	if !CheckFileIsExist(filepath) {
		//不存在创建
		file, err = os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	} else {
		file, err = os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND, 0666)
	}

	if err != nil {
		fmt.Printf("open file error=%v\n", err)
		return file, err
	}
	return file, nil
}

// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//返回程序目录/data
func GetDataPath() string {
	path, _ := os.Getwd()
	datapath := filepath.Join(path, "data")

	if !CheckFileIsExist(datapath) {
		err := os.Mkdir(datapath, 0777)
		if err != nil {
			fmt.Printf("文件夹：%s----创建失败[%v]\n", datapath, err)
		} else {
			fmt.Printf("文件夹：%s----创建成功\n", datapath)
		}
	}
	return datapath
}

//获取传入的指定目录/data/{}
func GetDirPath(dirname string) string {
	datapath := filepath.Join(GetDataPath(), dirname)

	if !CheckFileIsExist(datapath) {
		err := os.Mkdir(datapath, 0777)
		if err != nil {
			fmt.Printf("文件夹：%s----创建失败[%v]\n", datapath, err)
		} else {
			fmt.Printf("文件夹：%s----创建成功\n", datapath)
		}
	}
	return datapath
}

//返回/data/{默认时间}.{文件扩展名}
func GetSaveFilePath(filename string, ext string) string {

	if ext == "" {
		ext = ".txt"
	}

	if filename != "" {
		return filepath.Join(GetDataPath(), filename+ext)
	}
	return filepath.Join(GetDataPath(), time.Now().Format("20060102150405")+ext)
}
