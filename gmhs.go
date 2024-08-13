package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

var (
	port   string = "9000"
	dir    string = "current directory"
	upload string
)

func defineFlags() {
	flag.StringVar(&dir, "dir", dir, "root directory")
	flag.StringVar(&port, "port", port, "your port")
	flag.StringVar(&upload, "upload", upload, "upload folder")

	flag.Parse()
}

func validateDirectories() error {
	if dir == "current directory" {
		dir, _ = os.Getwd()
		fmt.Printf("Serving files from current directory: %s\n", dir)
	} else {
		if _, err := os.Stat(dir); errors.Is(err, fs.ErrNotExist) {
			fmt.Println("Error Validating Directories")
			fmt.Printf("Directory %s does not exist\n", dir)
			return err
		}
		fmt.Printf("Serving files from directory: %s\n", dir)
	}

	if upload != "" {
		if _, err := os.Stat(upload); errors.Is(err, fs.ErrNotExist) {
			fmt.Printf("Upload directory '%s' does not exist... creating it\n", upload)
			err := os.Mkdir(upload, 0222)
			if err != nil {
				fmt.Printf("Error creating directory %s\n", upload)
				fmt.Print("Error:", err)
				return err
			}
			fmt.Printf("upload directory at URL: %s\n", path.Join("localhost:"+port, filepath.Base(upload)))
		} else {
			fmt.Printf("Upload directory at: %s\n", filepath.Dir(upload))
		}
	}
	return nil
}

func main() {
	defineFlags()
	err := validateDirectories()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", http.FileServer(http.Dir(dir)).ServeHTTP)

	fmt.Printf("Listening on: %s\n", "http://localhost:"+port)
	http.ListenAndServe(":"+port, mux)
}
