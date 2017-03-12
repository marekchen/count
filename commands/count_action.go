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
	if path == "" {
		path = getSrcFullPath()
	}
	abspath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	println("abs", abspath)
	if eachfolder {
		ignorefile := filepath.Join(abspath, ".gitignore")
		ignores := new(list.List)
		if PathExists(ignorefile) {
			ignore, _ := gitignore.NewGitIgnore(ignorefile)
			ignores.PushBack(ignore)
		}
		dirList := getDirList(abspath)
		for e := dirList.Front(); nil != e; e = e.Next() {
			path := e.Value.(string)
			fileList := getList(path, ignores)
			count := getLines(fileList, suffix)
			fmt.Printf("Dir:%s; Lines:%d\r\n", filepath.Base(path), count)
		}
	} else {
		ignores := new(list.List)
		fileList := getList(abspath, ignores)
		count := getLines(fileList, suffix)
		fmt.Printf("Dir:%s; Lines:%d\r\n", filepath.Base(abspath), count)
	}
	return nil
}

func getDirList(fullPath string) (lst list.List) {
	dir, _ := ioutil.ReadDir(fullPath)
	for _, fi := range dir {
		ignorefile := filepath.Join(fullPath, ".gitignore")
		absPath := filepath.Join(fullPath, fi.Name())
		if PathExists(ignorefile) {
			ignore, _ := gitignore.NewGitIgnore(ignorefile)
			if ignore.Match(absPath, fi.IsDir()) {
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
			lst.PushBack(absPath)
		}
	}
	return
}

func getList(fullPath string, ignores *list.List) (lst list.List) {
	dir, _ := ioutil.ReadDir(fullPath)
	for _, fi := range dir {
		ignorefile := filepath.Join(fullPath, ".gitignore")
		absPath := filepath.Join(fullPath, fi.Name())
		if PathExists(ignorefile) {
			ignore, _ := gitignore.NewGitIgnore(ignorefile)
			ignores.PushBack(ignore)
		}
		match := false
		for e := ignores.Front(); nil != e; e = e.Next() {
			ignore := e.Value.(gitignore.IgnoreMatcher)
			if ignore.Match(absPath, fi.IsDir()) {
				match = true
				break
			}
		}
		if match {
			//println("ignored:" + absPath)
			continue
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
			list2 := getList(absPath, ignores)
			lst.PushBackList(&list2)
		} else {
			lst.PushBack(absPath)
		}
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

func getSrcFullPath() (fullPath string) {
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
		//println(path)
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
