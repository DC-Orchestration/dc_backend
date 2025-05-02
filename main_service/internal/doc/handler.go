package doc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"main_service/internal/pubsub"
	"main_service/internal/storage"
	"main_service/pkg/models"
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
func GetAllDocumentsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupID, err := strconv.Atoi(r.URL.Query().Get("group_id"))
		if err != nil {
			http.Error(w, "Invalid group ID", http.StatusBadRequest)
			return
		}

		var documents []models.Document
		if err := db.Where("group_id = ?", groupID).Find(&documents).Error; err != nil {
			http.Error(w, "Failed to fetch documents", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(documents)
	}
}

func GetDocumentByIDHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		docID, err := strconv.Atoi(r.URL.Query().Get("doc_id"))
		if err != nil {
			http.Error(w, "Invalid document ID", http.StatusBadRequest)
			return
		}

		var document models.Document
		if err := db.First(&document, docID).Error; err != nil {
			http.Error(w, "Document not found", http.StatusNotFound)
			return
		}

		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(document)
	}
}
func EditDocumentHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		docID, err := strconv.Atoi(r.URL.Query().Get("doc_id"))
		if err != nil {
			http.Error(w, "Invalid document ID", http.StatusBadRequest)
			return
		}

		var updatedDocument models.Document
		if err := json.NewDecoder(r.Body).Decode(&updatedDocument); err != nil {
			http.Error(w, "Failed to parse body", http.StatusBadRequest)
			return
		}

		var document models.Document
		if err := db.First(&document, docID).Error; err != nil {
			http.Error(w, "Document not found", http.StatusNotFound)
			return
		}
        document.Version++
		document.DocName = updatedDocument.DocName
		document.Content = updatedDocument.Content
		document.Hash = updatedDocument.Hash 
		if err := db.Save(&document).Error; err != nil {
			http.Error(w, "Failed to update document", http.StatusInternalServerError)
			return
		}
		auditLog := models.AuditLog{
			DocID:     document.ID,
			Action:    "edit",
			UserID:    document.UploaderID,
			Hash:      document.Hash,
			PrevHash:  document.PrevHash,
			Summary:   "Document edited",
			Timestamp: document.UpdatedAt,
		}//todo :send to ai here to generate summary
		if err := db.Create(&auditLog).Error; err != nil {
			http.Error(w, "Failed to create audit log", http.StatusInternalServerError)
			return

		}
		notification := map[string]interface{}{
			"doc_id":     document.ID,
			"group_id":   document.GroupID,
			"doc_name":   document.DocName,
			"uploader_id": document.UploaderID,
			"path":       document.Path,
			"hash":       document.Hash,
		}
		redisPublisher := pubsub.RedisPublisher{}
		if err := redisPublisher.PublishNewDocument(notification); err != nil {
			http.Error(w, "Failed to publish to Redis", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(document)
	}
}
func DownloadDocumentHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		docID, err := strconv.Atoi(r.URL.Query().Get("doc_id"))
		if err != nil {
			http.Error(w, "Invalid document ID", http.StatusBadRequest)
			return
		}

		var document models.Document
		if err := db.First(&document, docID).Error; err != nil {
			http.Error(w, "Document not found", http.StatusNotFound)
			return
		}

		fileData, err := storage.ReadDocumentFile(document.Path)
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}
        auditlog:=models.AuditLog{
			DocID:     document.ID,
			Action:    "download",
			UserID:    document.UploaderID,
			Hash:      document.Hash,
			PrevHash:  document.PrevHash,
			Summary:   "Document downloaded",
			Timestamp: document.UpdatedAt,	
		}
		if err := db.Create(&auditlog).Error; err != nil {
			http.Error(w, "Failed to create audit log", http.StatusInternalServerError)
			return
		}
		notification := map[string]interface{}{
			"doc_id":     document.ID,
			"group_id":   document.GroupID,
			"doc_name":   document.DocName,
			"uploader_id": document.UploaderID,
			"path":       document.Path,
			"hash":       document.Hash,
		}
		redisPublisher := pubsub.RedisPublisher{}
		if err := redisPublisher.PublishNewDocument(notification); err != nil {
			http.Error(w, "Failed to publish to Redis", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", document.DocName))
		w.WriteHeader(http.StatusOK)
		w.Write(fileData)
	}
}
func GetDocumentVersionsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		docID, err := strconv.Atoi(r.URL.Query().Get("doc_id"))
		if err != nil {
			http.Error(w, "Invalid document ID", http.StatusBadRequest)
			return
		}

		var versions []models.Document
		if err := db.Where("id = ?", docID).Order("created_at").Find(&versions).Error; err != nil {
			http.Error(w, "Failed to fetch document versions", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(versions)
	}
}
func SubmitDocumentChangesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	
		var changes map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&changes); err != nil {
			http.Error(w, "Failed to parse changes", http.StatusBadRequest)
			return
		}
		docID, ok := changes["doc_id"].(float64)
		if !ok {
			http.Error(w, "Invalid document ID", http.StatusBadRequest)
			return	
		}
		
		Changes := map[string]interface{}{
			"doc_id":     int(docID),
			"oldversion":    changes["oldversion"],
			"newversion":    changes["newversion"],
			"uploader_id": changes["uploader_id"],
		}
		redisPublisher := pubsub.RedisPublisher{}
		if err := redisPublisher.PublishIservice(Changes); err != nil {
			http.Error(w, "Failed to publish to Redis", http.StatusInternalServerError)
			return
		}
		
		

		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Changes submitted for AI analysis"))
	}
}
