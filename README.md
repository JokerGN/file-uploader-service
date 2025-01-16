
# File Uploader Service

A simple and efficient Go-based file upload service supporting both **single file uploads** and **multi-part (chunked) uploads**. Designed with clean architecture principles and comprehensive testing.

## 🚀 Features

- **Single File Uploads**: Upload complete files in one request.
- **Multi-Part Uploads**: Upload large files in chunks and assemble them.
- **Session Management**: Tracks file upload progress for multi-part uploads.
- **Automatic Cleanup**: Removes expired sessions and temporary chunk files.
- **Automated Testing**: Fully tested with Go's `httptest` for reliability.

## 📂 Project Structure

```
file_uploader_service/
├── main.go
├── fileuploader/
│   ├── handler/
│   │   ├── upload_handler.go
│   │   └── upload_handler_test.go
│   └── service/
│       └── upload_service.go
├── uploads/
│   ├── chunks/       # Temporary storage for file chunks
│   └── <uploaded_files>
└── go.mod
```

## ⚙️ Installation

1. **Clone the Repository**
   ```bash
   git clone https://github.com/JokerGN/file-uploader-service.git
   cd file-uploader-service
   ```

2. **Install Dependencies**
   ```bash
   go mod tidy
   ```

3. **Run the Server**
   ```bash
   go run main.go
   ```

   Server starts on `http://localhost:8080`

## 📥 API Endpoints

### 1. **Single File Upload**

- **Endpoint:** `POST /upload`
- **Body:** `multipart/form-data` with the key `"file"`

**Example:**

```bash
curl -X POST -F "file=@/path/to/file.txt" http://localhost:8080/upload
```

### 2. **Multi-Part Upload**

#### **a. Start Upload Session**

- **Endpoint:** `GET /start`
- **Query Params:**
    - `file_name`: Name of the final file
    - `total_chunks`: Total number of chunks

```bash
curl "http://localhost:8080/start?file_name=largefile.zip&total_chunks=5"
```

#### **b. Upload Chunks**

```bash
curl -X POST -F "chunk=@chunk0" "http://localhost:8080/upload_chunk?session_id=<session-id>&chunk_index=0"
```

#### **c. Complete the Upload**

```bash
curl "http://localhost:8080/complete_upload?session_id=<session-id>"
```

## 🧪 Running Tests

Run automated tests:

```bash
go test -v ./...
```

## 📄 License

This project is licensed under the MIT License.
