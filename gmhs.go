package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
)

var port string = "9000"
var dir string = "."
var upload string

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func DefineFlags() {
	flag.StringVar(&dir, "dir", dir, "root directory")
	flag.StringVar(&port, "port", port, "your port")
	flag.StringVar(&upload, "upload", upload, "upload folder")

	flag.Parse()
}
func receiveFile(uploadDir string) {
	http.HandleFunc("/"+uploadDir, func(w http.ResponseWriter, r *http.Request) {
		//check if the request is of method type POST
		if r.Method != "POST" {
			http.Error(w, "This page is only for file uploads", http.StatusMethodNotAllowed)
		}
		//parse the multipart form in the request
		err := r.ParseMultipartForm(20 << 20)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		log.Print("file received: ", handler.Filename)

		http.ResponseWriter.Write(w, []byte("file uploaded: "+handler.Filename))
	})
}
func main() {
	DefineFlags()

	//correctly format user supplied port number
	port = ":" + port

	//print directory being hosted
	if dir == "." {
		if directory, err := os.Getwd(); err != nil {
			log.Fatal(err)
		} else {
			log.Print("sharing content of current directory: ", directory)
		}
	}
	//check if user supplied path exists, then print path
	if dir != "." {
		if _, err := os.Stat(dir); err != nil {
			if os.IsNotExist(err) {
				log.Fatalf("Directory %s doesn't exist:", dir)
			} else {
				log.Fatalf("Encountered a problem checking the directory '%s': %s\n", dir, err)
			}
		} else {
			log.Printf("sharing content of directory: %s ", dir)
		}
	}

	//configure and start server
	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)

	//create and start sharing directory for uploads
	if len(upload) > 0 {
		err := os.Mkdir(upload, 0222)
		if err != nil {
			log.Fatal(err)
		}

		receiveFile(upload)

		log.Printf("upload directory at %s", upload)

		ip := GetOutboundIP().String()

		log.Printf("listening on: %s %s", ip, port)

	}
	//start server
	log.Fatal(http.ListenAndServe(port, nil))
}
