package model

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"main/database"
	"main/utils"
)

type FishAddress struct {
	gorm.Model
	Prefix string  `json:"prefix"`
	Suffix string  `json:"suffix"`
	Addr   string  `json:"addr"`
	Prov   string  `json:"prov"`
	IsUsed int     `json:"is_used"`
	UsedAt int     `json:"used_at"`
	Amount float32 `json:"amount"`
}

func GetFishAddresses(c *fiber.Ctx) error {
	db := database.DBConn
	var fishAddresses []FishAddress
	db.Find(&fishAddresses)
	return c.JSON(fishAddresses)
}

func GetFishAddressByPage(p *utils.Pagination) (fishes []FishAddress, err error) {
	db := database.DBConn
	err = db.Model(&FishAddress{}).Scopes(p.GormPaginate()).Find(&fishes).Error
	if err != nil {
		return nil, err
	}
	var total int64
	db.Model(&FishAddress{}).Count(&total)
	p.Total = cast.ToInt(total)
	return
}

func GetFishAddress(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var fishAddresses FishAddress
	db.Find(&fishAddresses, id)
	return c.JSON(fishAddresses)
}

func NewFishAddress(c *fiber.Ctx) error {
	db := database.DBConn
	fishAddress := new(FishAddress)
	if err := c.BodyParser(fishAddress); err != nil {
		return c.Status(503).SendString(err.Error())
	}
	db.Create(&fishAddress)
	return c.JSON(fishAddress)
}

func DeleteFishAddress(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn

	var fishAddress FishAddress
	db.First(&fishAddress, id)
	if fishAddress.Addr == "" {
		return c.Status(500).SendString("No FishAddress Found with ID")
	}
	db.Delete(&fishAddress)
	return c.SendString("FishAddress Successfully deleted")
}
