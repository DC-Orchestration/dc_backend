package audit


	import (
		"encoding/json"
		"net/http"
		"strconv"
	
		"gorm.io/gorm"
		"github.com/gorilla/mux"
		"main_service/pkg/models"
	)
	
	func GetDocumentAuditLogs(db *gorm.DB) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			docID, err := strconv.Atoi(mux.Vars(r)["doc_id"])
			if err != nil {
				http.Error(w, "Invalid document ID", http.StatusBadRequest)
				return
			}
	
			var logs []models.AuditLog
			if err := db.Where("doc_id = ?", docID).Order("timestamp asc").Find(&logs).Error; err != nil {
				http.Error(w, "Failed to fetch audit logs", http.StatusInternalServerError)
				return
			}
	
			json.NewEncoder(w).Encode(logs)
		}
	}
	
	func GetUserAuditLogs(db *gorm.DB) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			userID, err := strconv.Atoi(mux.Vars(r)["user_id"])
			if err != nil {
				http.Error(w, "Invalid user ID", http.StatusBadRequest)
				return
			}
	
			var logs []models.AuditLog
			if err := db.Where("user_id = ?", userID).Order("timestamp desc").Find(&logs).Error; err != nil {
				http.Error(w, "Failed to fetch user audit logs", http.StatusInternalServerError)
				return
			}
	
			json.NewEncoder(w).Encode(logs)
		}
	}
	