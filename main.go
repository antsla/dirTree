package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
)

func getFilesLen(files []os.FileInfo, printFiles bool) (filesNum int) {
	if printFiles == true {
		filesNum = len(files)
	} else {
		filesNum = 0
		for _, f := range files {
			if f.IsDir() {
				filesNum++
			}
		}
	}
	return
}

func getSizeString(file os.FileInfo) (sizeString string) {
	if size := file.Size(); size == 0 {
		sizeString = " (empty)"
	} else {
		sizeString = fmt.Sprintf(" (%vb)", size)
	}
	return
}

func getDelimiters(index, filesNum int) (lastDelimiter, lastDelimiterDir string) {
	lastDelimiter = "├───"
	lastDelimiterDir = "│"
	if index == filesNum {
		lastDelimiter = "└───"
		lastDelimiterDir = ""
	}
	return
}

func printFilesRecursive(out io.Writer, path string, printFiles bool, subPath string) error {
	files, _ := ioutil.ReadDir(path)
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
	filesNum := getFilesLen(files, printFiles)
	index := 1
	for _, file := range files {
		lastDelimiter, lastDelimiterDir := getDelimiters(index, filesNum)
		if file.IsDir() && printFiles == false || printFiles == true {
			index++
		}

		if !file.IsDir() && printFiles == true {
			sizeString := getSizeString(file)
			fmt.Fprint(out, subPath+lastDelimiter+file.Name()+sizeString+"\n")
		}

		if file.IsDir() {
			fmt.Fprint(out, subPath+lastDelimiter+file.Name()+"\n")
			printFilesRecursive(out, path+string(os.PathSeparator)+file.Name(), printFiles, subPath+lastDelimiterDir+"\t")
		}
	}

	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	err := printFilesRecursive(out, path, printFiles, "")
	if err != nil {
		panic(err.Error())
	}

	return nil
}

func main() {
	out := os.Stdout                               // получение выходного потока, сюда мы будем писать данные
	if !(len(os.Args) == 2 || len(os.Args) == 3) { // проверка количества входных аргументов
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f" // печатать ли файли или только директории
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
