package doc

import (
	"io"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"main_service/pkg/models"
	"main_service/internal/storage"
	"main_service/internal/pubsub"
)

func UploadDocumentHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		title := r.FormValue("doc_name")
		groupID, err := strconv.Atoi(r.FormValue("group_id"))
		if err != nil {
			http.Error(w, "Invalid group ID", http.StatusBadRequest)
			return
		}
		uploaderID, err := strconv.Atoi(r.FormValue("uploader_id"))
		if err != nil {
			http.Error(w, "Invalid uploader ID", http.StatusBadRequest)
			return
		}
		hash := r.FormValue("hash")
		if hash == "" {
			http.Error(w, "Missing document hash", http.StatusBadRequest)
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Failed to get file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}

		path, err := storage.SaveDocumentFile(handler.Filename, fileBytes)
		if err != nil {
			http.Error(w, "Failed to store file", http.StatusInternalServerError)
			return
		}

		doc := models.Document{
			DocName:    title,
			GroupID:    uint(groupID),
			UploaderID: uint(uploaderID),
			Hash:       hash,
			PrevHash:   "", 
			Path:       path,
			Version:    1,
		}

		if err := db.Create(&doc).Error; err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		notification := map[string]interface{}{
			"doc_id":     doc.ID,
			"group_id":   doc.GroupID,
			"doc_name":   doc.DocName,
			"uploader_id": doc.UploaderID,
			"path":       doc.Path,
			"hash":       doc.Hash,
		}

		redisPublisher := pubsub.RedisPublisher{}
		if err := redisPublisher.PublishNewDocument(notification); err != nil {
			http.Error(w, "Failed to publish to Redis", http.StatusInternalServerError)
			return
		}

		auditLog := models.AuditLog{
			DocID:     doc.ID,
			Action:    "upload",
			UserID:    uint(uploaderID),
			Hash:      hash,
			PrevHash:  "",
			Summary:   "Document uploaded",
			Timestamp: doc.CreatedAt,
		}
		if err := db.Create(&auditLog).Error; err != nil {
			http.Error(w, "Failed to create audit log", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Document uploaded successfully"))
	}
}
