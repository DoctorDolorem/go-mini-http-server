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
		fsu := http.FileServer(http.Dir(upload))
		http.Handle("/upload", fsu)
		log.Printf("upload directory at %s", upload)
	}

	ip := GetOutboundIP().String()

	log.Printf("listening on: %s %s", ip, port)

	log.Fatal(http.ListenAndServe(port, nil))

}
