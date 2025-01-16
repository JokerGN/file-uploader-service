package main

import (
	"file-uploader-service/fileuploader/handler"
	"file-uploader-service/fileuploader/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ensureUploadDirExists("uploads")
	ensureUploadDirExists("uploads/chunks")

	uploadService := service.NewUploadService(5 * time.Minute)
	uploadHandler := handler.NewUploadHandler(uploadService)

	http.HandleFunc("/start", uploadHandler.StartSessionHandler)
	http.HandleFunc("/upload_chunk", uploadHandler.UploadChunkHandler)
	http.HandleFunc("/complete_upload", uploadHandler.CompleteUploadHandler)
	http.HandleFunc("/upload", uploadHandler.SingleFileUploadHandler)

	go uploadService.CleanExpiredSessions()

	fmt.Println("Servier started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func ensureUploadDirExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			log.Fatalf("Failed to create upload directory: %v", err)
		}
	}
}
