package handler

import (
	"file-uploader-service/fileuploader/service"
	"fmt"
	"net/http"
	"strconv"
)

type UploadHandler struct {
	service *service.UploadService
}

func NewUploadHandler(s *service.UploadService) *UploadHandler {
	return &UploadHandler{service: s}
}

func (h *UploadHandler) StartSessionHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file_name")
	totalChunksStr := r.URL.Query().Get("total_chunks")
	totalChunks, err := strconv.Atoi(totalChunksStr)
	if err != nil || filename == "" || totalChunks <= 0 {
		http.Error(w, "Invalid file name or total chunks", http.StatusBadRequest)
		return
	}
	id := h.service.StartSession(filename, totalChunks)
	fmt.Fprintf(w, "Session ID: %s\n", id)
}
