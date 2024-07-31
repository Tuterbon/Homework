package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"slices"
	"sort"
)

func recurseFull(out io.Writer, path string, str string) {
	dir, err := os.ReadDir(path)
	if err != nil {
		panic(err.Error())
	}
	sort.Slice(dir, func(i, j int) bool { return dir[i].Name() < dir[j].Name() })
	for i, data := range dir {
		info, err := data.Info()
		if err != nil {
			panic(err.Error())
		}
		switch i == len(dir)-1 {
		case true:
			switch data.IsDir() {
			case true:
				fmt.Fprintf(out, str+"└───%v\n", info.Name())
				recurseFull(out, path+"/"+info.Name(), str+"\t")
			case false:
				switch info.Size() == 0 {
				case true:
					fmt.Fprintf(out, str+"└───%v (empty)\n", info.Name())
				case false:
					fmt.Fprintf(out, str+"└───%v (%vb)\n", info.Name(), info.Size())
				}

			}
		case false:
			switch data.IsDir() {
			case true:
				fmt.Fprintf(out, str+"├───%v\n", info.Name())
				recurseFull(out, path+"/"+info.Name(), str+"│\t")
			case false:
				switch info.Size() == 0 {
				case true:
					fmt.Fprintf(out, str+"├───%v (empty)\n", info.Name())
				case false:
					fmt.Fprintf(out, str+"├───%v (%vb)\n", info.Name(), info.Size())
				}
			}
		}
	}
}

func recurseDir(out io.Writer, path string, str string) {
	dir, err := os.ReadDir(path)
	if err != nil {
		panic(err.Error())
	}
	sort.Slice(dir, func(i, j int) bool { return dir[i].Name() < dir[j].Name() })
	dirs := slices.DeleteFunc(dir, func(d fs.DirEntry) bool { return !d.IsDir() })
	for i, data := range dirs {
		if data == nil {
			continue
		}
		switch i == len(dirs)-1 {
		case true:
			fmt.Fprintf(out, str+"└───%v\n", data.Name())
			recurseDir(out, path+"/"+data.Name(), str+"\t")
		case false:
			fmt.Fprintf(out, str+"├───%v\n", data.Name())
			recurseDir(out, path+"/"+data.Name(), str+"│\t")
		}
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	str := ""
	switch printFiles {
	case true:
		recurseFull(out, path, str)
	case false:
		recurseDir(out, path, str)
	}
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
