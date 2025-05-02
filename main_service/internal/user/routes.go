package user


import (

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func RegisterUserRoutes(r chi.Router, db *gorm.DB) {
	r.Route("/user", func(r chi.Router) {
		r.Post("/login", LoginHandler(db))
		r.Post("/login/signature", LoginWithSignatureHandler(db))
		r.Post("/accept-invite", AcceptInviteHandler(db))
	})
}
