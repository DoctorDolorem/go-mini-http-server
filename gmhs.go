package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	switch r.Method {
	case "GET":
		w.Write([]byte("upload page"))

	case "POST":
		uploadFile(w, r)
	}
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error retrieving file")
		fmt.Println(err)
		return
	}
	defer file.Close()

	f, err := os.OpenFile(upload+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println("Error creating file:", err)
		w.Write([]byte("Error uploading file"))
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("File uploaded: " + handler.Filename + "\n"))
	log.Printf("File uploaded: %s\n", handler.Filename)
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		log.Println("Error copying file:", err)
		w.Write([]byte("Error uploading file"))
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("File uploaded successfully: " + handler.Filename + "\n"))
	log.Printf("File uploaded successfully: %s\n", handler.Filename)
}

func main() {
	defineFlags()

	err := validateShareDir()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", http.FileServer(http.Dir(dir)).ServeHTTP)

	ip, err := grabIP()
	if err != nil {
		ip = "localhost"
		fmt.Println("Error grabbing IP", err)
	}
	fmt.Printf("Available at: http://%s:%s\n", ip, port)

	if upload != "" {
		if err := validateUploadDir(); err != nil {
			log.Println("Error validating upload directory:", err)
		}
		mux.HandleFunc("/upload", uploadPage)
		fmt.Printf("Upload at: http://%s:%s/upload\n", ip, port)
	}

	go http.ListenAndServe(":"+port, mux)

	fmt.Fprint(os.Stdout, "Press Enter to exit ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input")
	}
	input = strings.TrimSpace(input)
	if input == "" {
		os.Exit(0)
	}

}
