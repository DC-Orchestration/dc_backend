package audit

import (

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func RegisterAuditRoutes(r chi.Router, db *gorm.DB) {
	r.Route("/audit", func(r chi.Router) {
		r.Get("/document/{doc_id}", GetDocumentAuditLogs(db))

		r.Get("/user/{user_id}", GetUserAuditLogs(db))
	})
}
