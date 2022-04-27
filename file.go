package ginCore

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
)

type FileNode struct {
	Name      string      `json:"name"`
	Parent    string      `json:"parent"`
	Root      string      `json:"root"`
	Relative  string      `json:"relative"`
	Path      string      `json:"path"`
	FileNodes []*FileNode `json:"children"`
}

func GetFileTree(rootpath string) (FileNode, error) {
	rootpath = filepath.Join(rootpath, "")
	root := FileNode{"Root", rootpath, rootpath, "", rootpath, []*FileNode{}}
	fileInfo, _ := os.Lstat(rootpath)
	walk(rootpath, rootpath, fileInfo, &root, `.md`)
	return root, nil
}

func walk(path string, root string, info os.FileInfo, node *FileNode, filter string) {
	// 列出当前目录下的所有目录、文件
	files := listFiles(path)

	// 遍历这些文件
	for _, filename := range files {

		// 拼接全路径
		fpath := filepath.Join(path, filename)
		if IsDir(fpath) || (!IsDir(fpath) && IsMatch(fpath, `.*.md`)) {
			// 构造文件结构
			fio, _ := os.Lstat(fpath)
			relative_path := fpath[len(root):]
			if relative_path[0] == '\\' {
				relative_path = relative_path[1:]
			}

			// 将当前文件作为子节点添加到目录下
			child := FileNode{filename, path, root, fpath[len(root):], fpath, []*FileNode{}}
			node.FileNodes = append(node.FileNodes, &child)

			// 如果遍历的当前文件是个目录，则进入该目录进行递归
			if fio.IsDir() {
				walk(fpath, root, fio, &child, filter)
			}
		}
	}

	return
}

func listFiles(dirname string) []string {
	f, _ := os.Open(dirname)
	names, _ := f.Readdirnames(-1)
	f.Close()
	sort.Strings(names)
	return names
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

func Write(file, mode, content string) (*string, error) {
	var f *os.File
	var err error
	if mode == "w" {
		f, err = os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			log.Fatal(err)
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {

			}
		}(f)
		writeString, err := f.WriteString(content)
		if err != nil {
			return nil, err
		}
		fmt.Println(writeString)
	} else if mode == "a" {
		f, err = os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			return nil, err
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {

			}
		}(f)
		writeString, err := f.WriteString(content)
		if err != nil {
			return nil, err
		}
		fmt.Println(writeString)
	}
	str := "success"
	return &str, nil
}

func Read(filepath string) (*string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var result = ""
	var buf [256]byte
	var content []byte
	for {
		n, err := file.Read(buf[:])
		if err == io.EOF {
			break
		}
		if err != nil {
			return &result, nil
		} else if err == io.EOF {
			content = append(content, buf[:n]...)
			break
		} else {
			content = append(content, buf[:n]...)
		}
	}
	result = string(content)
	return &result, nil
}
