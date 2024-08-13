package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"net"
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
			//fmt.Printf("upload directory at URL: %s\n", path.Join("localhost:"+port, filepath.Base(upload)))
			fmt.Printf("upload directory at URL: %s\n", path.Join("localhost:"+port, "upload"))
		} else {
			fmt.Printf("Upload directory at: FULL PATH HERE%s\n", filepath.Dir(upload))
		}
	}
	return nil
}

func grabIP() (string, error) {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: net.ParseIP("1.1.1.1"), Port: 53})
	if err != nil {
		fmt.Println("Error grabbing IP")
		return "", err
	}
	ip := conn.LocalAddr().(*net.UDPAddr).IP.String()
	conn.Close()
	return ip, nil
}
func main() {
	defineFlags()

	err := validateDirectories()
	if err != nil {
		fmt.Print("Error:", err)
		os.Exit(1)
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/", http.FileServer(http.Dir(dir)).ServeHTTP)

	if upload != "" {
		mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Upload page"))
		})
	}

	ip, err := grabIP()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Available at: http://%s:%s\n", ip, port)
	http.ListenAndServe(":"+port, mux)
}
