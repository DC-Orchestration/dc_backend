package group

import (
	"encoding/json"
	"fmt"
	"main_service/pkg/models"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

func CreateGroupHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var group models.Group
		if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
			http.Error(w, "Invalid group data", http.StatusBadRequest)
			return
		}

		if group.Name == "" {
			http.Error(w, "Group name is required", http.StatusBadRequest)
			return
		}

		if err := db.Create(&group).Error; err != nil {
			http.Error(w, "Failed to create group", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(group)
	}
}
func GetGroupByIDHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupID := r.URL.Query().Get("group_id")
		if groupID == "" {
			http.Error(w, "Group ID is required", http.StatusBadRequest)
			return
		}

		var group models.Group
		if err := db.First(&group, groupID).Error; err != nil {
			http.Error(w, "Group not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(group)
	}
}
func GetAllGroupsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var groups []models.Group
		if err := db.Find(&groups).Error; err != nil {
			http.Error(w, "Failed to fetch groups", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(groups)
	}
}
func UpdateGroupHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupID := r.URL.Query().Get("group_id")
		if groupID == "" {
			http.Error(w, "Group ID is required", http.StatusBadRequest)
			return
		}

		var updatedGroup models.Group
		if err := json.NewDecoder(r.Body).Decode(&updatedGroup); err != nil {
			http.Error(w, "Invalid group data", http.StatusBadRequest)
			return
		}

		if err := db.Model(&models.Group{}).Where("id = ?", groupID).Updates(updatedGroup).Error; err != nil {
			http.Error(w, "Failed to update group", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedGroup)
	}
}
func ListGroupMembersHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupID, err := strconv.Atoi(r.URL.Query().Get("group_id"))
		if err != nil {
			http.Error(w, "Invalid group ID", http.StatusBadRequest)
			return
		}

		var group models.Group
		if err := db.Preload("Members").First(&group, groupID).Error; err != nil {
			http.Error(w, "Group not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(group.Members)
	}
}
func InviteUserToGroupHandler(db *gorm.DB, emailSender func(to, subject, body string) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var invite struct {
			GroupID uint   `json:"group_id"`
			Email   string `json:"email"`
		}

		if err := json.NewDecoder(r.Body).Decode(&invite); err != nil {
			http.Error(w, "Invalid request data", http.StatusBadRequest)
			return
		}

		var group models.Group
		if err := db.First(&group, invite.GroupID).Error; err != nil {
			http.Error(w, "Group not found", http.StatusNotFound)
			return
		}

		var user models.User
		if err := db.Where("email = ?", invite.Email).First(&user).Error; err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if err := db.Model(&group).Association("Members").Append(&user); err != nil {
			http.Error(w, "Failed to add user to group", http.StatusInternalServerError)
			return
		}

		subject := "You’ve been invited to a group!"
		body := fmt.Sprintf("Hello %s,\n\nYou've been added to the group: .\nLogin to view and collaborate on the documents in the group.\n\n– Team",  group.Name)

		if err := emailSender(user.Email, subject, body); err != nil {
			http.Error(w, "Failed to send invitation email", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User invited and email sent"))
	}
}

