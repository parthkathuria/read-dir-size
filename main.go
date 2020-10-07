package main

import (
	"encoding/csv"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
)

func calculateDirSize() (dirsize int64, err error) {
	err = filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Println(err)
			} else {
				//log.Println(path)
				dirsize += info.Size()
			}
			return nil
		})
	return
}

type nameSize struct {
	name string
	size int64
}

func main() {
	// read all files/directories in current directory
	dirs, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Println(err)
		return
	}
	var nameSizeArray []nameSize
	bar := progressbar.Default(int64(len(dirs)))
	for _, dir := range dirs {
		bar.Describe(dir.Name())
		// if it is not dir, continue to next dir
		if !dir.IsDir() {
			bar.Add(1)
			continue
		}

		err = os.Chdir(dir.Name())
		if err != nil {
			fmt.Println(err)
			return
		}

		size, err := calculateDirSize()
		if err != nil {
			log.Printf("%s: error - %v", dir.Name(), err)
		}
		ns := nameSize{dir.Name(), size}
		nameSizeArray = append(nameSizeArray, ns)
		bar.Add(1)
		err = os.Chdir("..")
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// sort all directories in current directory according to size
	sort.Slice(nameSizeArray, func(i, j int) bool {
		return nameSizeArray[i].size > nameSizeArray[j].size
	})

	file, err := os.Create("result.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, ns := range nameSizeArray {
		err := writer.Write([]string{ns.name, fmt.Sprintf("%d", ns.size)})
		checkError("Cannot write to file", err)
	}
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
