package handler

import (
	"file-uploader-service/fileuploader/service"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

func (h *UploadHandler) UploadChunkHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	chunkIndexStr := r.URL.Query().Get("chunk_index")
	chunkIndex, err := strconv.Atoi(chunkIndexStr)
	if err != nil || !h.service.AddChunk(sessionID, chunkIndex) {
		http.Error(w, "Invalid session or chunk index", http.StatusUnauthorized)
		return
	}

	file, _, err := r.FormFile("chunk")
	defer file.Close()
	chunkPath := filepath.Join("../../uploads/chunks/", fmt.Sprintf("%s_%d", sessionID, chunkIndex))
	fmt.Println(chunkPath)
	out, _ := os.Create(chunkPath)
	defer out.Close()
	io.Copy(out, file)
	fmt.Fprintf(w, "Chunk %d uploaded\n", chunkIndex)
}

func (h *UploadHandler) CompleteUploadHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")

	session, exists := h.service.GetSession(sessionID)
	if !exists || !h.service.IsUploadComplete(sessionID) {
		http.Error(w, "Upload incomplete or invalid session", http.StatusBadRequest)
		return
	}

	finalFilePath := filepath.Join("../../uploads", session.FileName)
	fmt.Println(finalFilePath)
	finalFile, err := os.Create(finalFilePath)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to create final file", http.StatusInternalServerError)
		return
	}
	defer finalFile.Close()

	for i := 0; i < session.TotalChunks; i++ {
		chunkPath := filepath.Join("../../uploads/chunks", fmt.Sprintf("%s_%d", sessionID, i))
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			fmt.Println(err)
			http.Error(w, fmt.Sprintf("Failed to read chunk %d", i), http.StatusInternalServerError)
			return
		}
		io.Copy(finalFile, chunkFile)
		chunkFile.Close()
		os.Remove(chunkPath)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File assembled successfully: %s\n", session.FileName)
}

func (h *UploadHandler) SingleFileUploadHandler(w http.ResponseWriter, r *http.Request) {
	file, header, _ := r.FormFile("file")
	defer file.Close()
	filePath := filepath.Join("../../uploads", header.Filename)
	out, _ := os.Create(filePath)
	defer out.Close()
	io.Copy(out, file)
	fmt.Fprintf(w, "Single file uploaded successfully: %s\n", header.Filename)
}
