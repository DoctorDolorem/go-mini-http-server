package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
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

func validateShareDir() error {
	if dir == "current directory" {
		dir, _ = os.Getwd()
		fmt.Printf("Serving files from current directory: %s\n", dir)
	} else {
		if _, err := os.Stat(dir); errors.Is(err, fs.ErrNotExist) {
			fmt.Printf("ABORT: Error validating share directory\n")
			fmt.Printf("Directory %s does not exist\n", dir)
			return err
		}
		fmt.Printf("Serving files from directory: %s\n", dir)
	}
	return nil
}

func validateUploadDir() error {
	if upload != "" {
		if _, err := os.Stat(upload); errors.Is(err, fs.ErrNotExist) {
			fmt.Printf("Upload directory '%s' does not exist... creating it\n", upload)
			err := os.Mkdir(upload, 0222)
			if err != nil {
				fmt.Printf("Error creating directory %s\n", upload)
				return err
			}
		} else {
			path, err := filepath.Abs(upload)
			if err != nil {
				fmt.Printf("Error getting absolute path of upload directory %s\n", upload)
				return err
			}
			fmt.Printf("Upload directory at: %s\n", path)
		}
	}
	return nil
}

func grabIP() (string, error) {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: net.ParseIP("192.168.1.1"), Port: 53})
	if err != nil {
		fmt.Println("Error grabbing IP")
		return "", err
	}
	ip := conn.LocalAddr().(*net.UDPAddr).IP.String()
	conn.Close()
	return ip, nil
}
func uploadPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Upload page"))
}

func main() {
	defineFlags()

	if err := validateShareDir(); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", http.FileServer(http.Dir(dir)).ServeHTTP)

	if upload != "" {
		if err := validateUploadDir(); err != nil {
			log.Printf("Error validating upload directory: %w\n", err)
		}
		mux.HandleFunc("/upload", uploadPage)
	}

	ip, err := grabIP()
	if err != nil {
		ip = "localhost"
		log.Printf("Error grabbing IP %w\n", err)
	}

	fmt.Printf("Available at: http://%s:%s\n", ip, port)
	fmt.Printf("Upload at: http://%s:%s/upload\n", ip, port)
	http.ListenAndServe(":"+port, mux)
}
