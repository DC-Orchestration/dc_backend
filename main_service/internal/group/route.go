package group

import (
	

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func RegisterGroupRoutes(r chi.Router, db *gorm.DB, emailSender func(to, subject, body string) error) {
	r.Route("/group", func(r chi.Router) {
		r.Post("/", CreateGroupHandler(db))

		r.Get("/{group_id}", GetGroupByIDHandler(db))

		r.Get("/", GetAllGroupsHandler(db))

		r.Put("/{group_id}", UpdateGroupHandler(db))

		r.Get("/{group_id}/members", ListGroupMembersHandler(db))

		r.Post("/invite", InviteUserToGroupHandler(db, emailSender))
	})
}
