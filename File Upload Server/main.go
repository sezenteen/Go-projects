package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Write([]byte("Use POST method"))
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	dst, err := os.Create("./uploads/" + handler.Filename)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s", handler.Filename)
	r.ParseMultipartForm(10 << 20)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	html := `
	<html>
	<body>
		<h2>Upload File</h2>
		<form enctype="multipart/form-data" action="/upload" method="post">
			<input type="file" name="file" />
			<input type="submit" value="upload" />
		</form>
	</body>
	</html>`
	w.Write([]byte(html))
}

func main() {
	os.MkdirAll("./uploads", os.ModePerm)

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/upload", uploadHandler)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}