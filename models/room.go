package models

import (
	"errors"
	"strconv"
	"webchat/database"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type RoomManager struct{}

type Room struct {
	gorm.Model
	Name        string `form:"name"`
	ManagerID   string
	Description string `form:"description"`
	Childrens   RelationUsers `sql:"TYPE:json"`
}

func NewRoomManager() *RoomManager {
	database.DB.AutoMigrate(&Room{})
	return &RoomManager{}
}

func (*RoomManager) CreateRoom(userID string, r *Room) error {
	r.ManagerID = userID
	if err := database.DB.Create(r).Error; err != nil {
		return err
	}
	return nil
}

func (*RoomManager) ListRooms(userID string) ([]Room, error) {
	var roomes []Room
	if err := database.DB.Where("manager_id = ?", userID).Find(&roomes).Error; err != nil {
		return nil, err
	}
	return roomes, nil
}

func (*RoomManager) UpdateRoom(ctx *gin.Context) error {
	var room *Room
	if err := ctx.ShouldBind(room); err != nil {
		return err
	}

	if err := database.DB.Model(&room).Update("description", room.Description).Error; err != nil {
		return err
	}
	return nil
}

func (*RoomManager) GetRoom(roomID string) (*Room, error) {
	var room Room
	if err := database.DB.Where("id = ?", roomID).Find(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (*RoomManager) DeleteRoom(id string) error {
	var room Room
	if err := database.DB.First(&room, id).Error; err != nil {
		return err
	}

	if room.ManagerID == id {
		database.DB.Unscoped().Delete(&room)
	} else {
		return errors.New("this user is not allow to delet this room")
	}
	return nil
}


func (r *Room) filterChilds (id string) {
	for i, user := range r.Childrens {
		if strconv.Itoa(int(user.ID)) == id {
			r.Childrens =  append(r.Childrens[:i], r.Childrens[i+1:]...)
		}
	}
}