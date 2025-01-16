package handler

import (
	"bytes"
	"file-uploader-service/fileuploader/service"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestSinglfileUploadHandler(t *testing.T) {
	uploadService := service.NewUploadService(5 * time.Minute)
	uploadHandler := NewUploadHandler(uploadService)

	fileContent := "This is a test file."
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "testfile.txt")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	_, err = io.Copy(part, strings.NewReader(fileContent))
	if err != nil {
		t.Fatalf("Failed to write file content: %v", err)
	}
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	uploadHandler.SingleFileUploadHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Single file uploaded successfully") {
		t.Errorf("Expected upload success message, got %v", rr.Body.String())
	}
	os.Remove("../../uploads/testfile.txt")
}

func TestMultiPartUploadHandler(t *testing.T) {
	uploadService := service.NewUploadService(5 * time.Minute)
	uploadHandler := NewUploadHandler(uploadService)

	reqStart := httptest.NewRequest(http.MethodGet, "/start?file_name=testfile.txt&total_chunks=2", nil)
	rrStart := httptest.NewRecorder()
	uploadHandler.StartSessionHandler(rrStart, reqStart)

	sessionID := strings.TrimPrefix(rrStart.Body.String(), "Session ID: ")
	sessionID = strings.TrimSpace(sessionID)

	uploadChunk(t, uploadHandler, sessionID, 0, "This is chunk 1.")

	uploadChunk(t, uploadHandler, sessionID, 1, "This is chunk 2.")

	reqComplete := httptest.NewRequest(http.MethodGet, "/complete_upload?session_id="+sessionID, nil)
	rrComplete := httptest.NewRecorder()
	uploadHandler.CompleteUploadHandler(rrComplete, reqComplete)

	if rrComplete.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rrComplete.Code)
	}
	if !strings.Contains(rrComplete.Body.String(), "File assembled successfully") {
		t.Errorf("Expected file assembled success message, got %v", rrComplete.Body.String())
	}
	os.Remove("../../uploads/testfile.txt")
}

func uploadChunk(t *testing.T, handler *UploadHandler, sessionID string, chunkIndex int, content string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("chunk", "chunk")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	_, err = io.Copy(part, strings.NewReader(content))
	if err != nil {
		t.Fatalf("Failed to write chunk content: %v", err)
	}
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload_chunk?session_id="+sessionID+"&chunk_index="+strconv.Itoa(chunkIndex), body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	handler.UploadChunkHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK for chunk %d, got %v", chunkIndex, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Chunk") {
		t.Errorf("Expected chunk upload success message, got %v", rr.Body.String())
	}
}
