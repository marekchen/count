package commands

import (
	"github.com/codegangsta/cli"
	"fmt"
	"path/filepath"
	"os"
	"container/list"
	"bufio"
	"log"
	"github.com/monochromegane/go-gitignore"
	"io/ioutil"
	"strings"
)

func countAction(c *cli.Context) error {
	eachfolder := c.Bool("eachfolder")
	path := c.String("path")
	suffix := c.String("suffix")

	fmt.Printf("eachfolder：%t,path: %s,suffix: %s\r\n", eachfolder, path, suffix)
	abspath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	println("abspath", abspath)
	if eachfolder {
		dirList := getDirList(abspath)
		for e := dirList.Front(); nil != e; e = e.Next() {
			path := e.Value.(string)
			fileList := getList(path)
			count := getLines(fileList, suffix)
			fmt.Printf("Dir:%s; Lines:%d\r\n", filepath.Base(path), count)
		}
	} else {
		fileList := getList(abspath)
		count := getLines(fileList, suffix)
		fmt.Printf("Dir:%s; Lines:%d\r\n", filepath.Base(abspath), count)
	}
	return nil
}

func getDirList(fullPath string) (lst list.List) {
	dir, _ := ioutil.ReadDir(fullPath)
	for _, fi := range dir {
		ignorefile := fullPath + "/.gitignore"
		absPath := fullPath + "/" + fi.Name()
		if PathExists(ignorefile) {
			gitignore1, _ := gitignore.NewGitIgnore(ignorefile)
			if gitignore1.Match(absPath, fi.IsDir()) {
				continue
			}
		}
		if strings.HasPrefix(fi.Name(), ".") {
			continue
		}
		if strings.HasPrefix(fi.Name(), "gradle") {
			continue
		}
		if fi.IsDir() {
			lst.PushBack(fullPath + "/" + fi.Name())
		}
	}
	return
}

func getList(fullPath string) (lst list.List) {
	dir, _ := ioutil.ReadDir(fullPath)
	for _, fi := range dir {
		ignorefile := fullPath + "/.gitignore"
		absPath := fullPath + "/" + fi.Name()
		if PathExists(ignorefile) {
			gitignore1, _ := gitignore.NewGitIgnore(ignorefile)
			if gitignore1.Match(absPath, fi.IsDir()) {
				continue
			}
		}
		if strings.HasPrefix(fi.Name(), ".") {
			continue
		}
		if strings.HasPrefix(fi.Name(), "gradle") {
			continue
		}
		if strings.HasSuffix(fi.Name(), ".aar") {
			continue
		}
		if strings.HasSuffix(fi.Name(), ".jar") {
			continue
		}
		if strings.HasSuffix(fi.Name(), ".iml") {
			continue
		}
		if strings.HasSuffix(fi.Name(), ".png") {
			continue
		}
		if strings.HasSuffix(fi.Name(), ".jpg") {
			continue
		}
		if fi.IsDir() {
			list2 := getList(fullPath + "/" + fi.Name())
			lst.PushBackList(&list2)
		} else {
			lst.PushBack(fullPath + "/" + fi.Name())
		}
	}
	return
}

func getFilelist(path string) (lst list.List) {
	fullPath := GetSrcFullPath()
	fmt.Println("fullPath:", fullPath)
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {

		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		if filepath.Base(path) == ".gitignore" {
			print("it's ignore file\n\r")
			return nil
		}

		ignorefile := fullPath + "/.gitignore"
		absPath := fullPath + "/" + path

		if PathExists(ignorefile) {
			gitignore1, _ := gitignore.NewGitIgnore(ignorefile)
			if !gitignore1.Match(absPath, false) {
				println("path:" + path)
				lst.PushBack(path)
			}
		} else {
			println(path)
			lst.PushBack(path)
		}

		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
	return
}

func PathExists(path string) (bool) {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func GetSrcFullPath() (fullPath string) {
	args := os.Args;
	parameterLen := len(args)
	if parameterLen == 1 {
		fullPath, _ = os.Getwd()
	}
	fullPath, _ = filepath.Abs(fullPath)
	return
}

func getLines(lst list.List, suffix string) (count int) {
	for e := lst.Front(); nil != e; e = e.Next() {
		path := e.Value.(string)

		if suffix != "" {
			suffixlist := strings.Split(suffix, ",")
			length := len(suffixlist)
			for i := 0; i < length; i++ {
				if strings.HasSuffix(path, suffixlist[i]) {
					count = count + ComputeLine(path)
				}
			}
		} else {
			count = count + ComputeLine(path)
		}
	}
	return
}

func ComputeLine(path string) (code int) {
	f, err := os.Open(path)
	if nil != err {
		log.Println(err)
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		code += 1
	}
	return
}
