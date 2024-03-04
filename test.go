// Remove the package declaration from the code block in the selection

func receiveFile(uploadDir string) {
	http.Handle("/upload/", http.StripPrefix("/upload/", http.FileServer(http.Dir(uploadDir))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Parse the multipart form in the request
		err := r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Get a file from the multipart form
		file, handler, err := r.FormFile("myFile")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		// Create a new file in the uploads directory
		dst, err := os.Create(filepath.Join(uploadDir, handler.Filename))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// ...
	})
}
