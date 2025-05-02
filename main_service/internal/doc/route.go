package doc

import (

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func RegisterDocRoutes(r chi.Router, db *gorm.DB) {
	r.Route("/doc", func(r chi.Router) {
		r.Post("/upload", UploadDocumentHandler(db))
		r.Get("/", GetAllDocumentsHandler(db))
		r.Get("/{doc_id}", GetDocumentByIDHandler(db))
		r.Put("/{doc_id}", EditDocumentHandler(db))
		r.Get("/{doc_id}/download", DownloadDocumentHandler(db))
		r.Get("/{doc_id}/versions", GetDocumentVersionsHandler(db))
		r.Post("/submit-changes", SubmitDocumentChangesHandler(db))
	})
}
