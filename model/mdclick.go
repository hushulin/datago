package model

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"main/database"
	"main/utils"
	"net/url"
	"strconv"
)

// MdClick md 点击记录 用于分析常用
// Ref 格式：type-random 由前端生成 首页下拉-date+random
type MdClick struct {
	gorm.Model
	MarkdownId uint   `json:"markdown_id"`
	Ref        string `json:"ref"`
	Score      uint   `json:"score"`
	Random     string `json:"random"`
}

func GetMdClicks(c *fiber.Ctx) error {
	db := database.DBConn
	var mdClicks []MdClick
	db.Find(&mdClicks)
	return c.JSON(mdClicks)
}

func GetMdClickByRandom(random string) bool {
	db := database.DBConn
	var count int64
	db.Model(&MdClick{}).Where("random = ?", random).Count(&count)
	if count > 0 {
		return true
	}
	return false
}

func GetMdClicks5(c *fiber.Ctx) error {
	db := database.DBConn
	var mdClicks []MdClick
	db.Order("id desc").Limit(5).Find(&mdClicks)
	return c.JSON(mdClicks)
}

func GetMdClickByPage(p *utils.Pagination) (mdClicks []MdClick, err error) {
	db := database.DBConn
	err = db.Model(&MdClick{}).Scopes(p.GormPaginate()).Order("id desc").Find(&mdClicks).Error
	if err != nil {
		return nil, err
	}
	var total int64
	db.Model(&MdClick{}).Count(&total)
	p.Total = cast.ToInt(total)
	return
}

func GetMdClick(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var mdClick MdClick
	db.Find(&mdClick, id)
	return c.JSON(mdClick)
}

func GetMdClickByKeyword(c *fiber.Ctx) error {
	keyword := c.Params("keyword")
	unescape, err := url.QueryUnescape(keyword)
	if err != nil {
		return err
	}
	db := database.DBConn
	var mdClicks []MdClick
	db.Where("`title` LIKE ? OR `content` LIKE ?", "%"+unescape+"%", "%"+unescape+"%").Find(&mdClicks)
	return c.JSON(mdClicks)
}

func GetMdClick2(c *fiber.Ctx) MdClick {
	id := c.Params("id")
	db := database.DBConn
	var mdClick MdClick
	db.Find(&mdClick, id)
	return mdClick
}

// NewMdClick 第一个要使用的接口，点击就产生一条
func NewMdClick(c *fiber.Ctx) error {
	db := database.DBConn
	mdClick := new(MdClick)
	if err := c.BodyParser(mdClick); err != nil {
		return c.Status(503).SendString(err.Error())
	}
	db.Create(&mdClick)
	return c.JSON(mdClick)
}

// NewMdClickForInner 第一个要使用的接口 内部使用
func NewMdClickForInner(markdownId uint, ref string, score uint, random string) *MdClick {
	db := database.DBConn
	mdClick := new(MdClick)
	mdClick.Ref = ref
	mdClick.Score = score
	mdClick.MarkdownId = markdownId
	mdClick.Random = random
	db.Create(&mdClick)
	return mdClick
}

func DeleteMdClick(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn

	var mdClick MdClick
	db.First(&mdClick, id)
	db.Delete(&mdClick)
	return c.SendString("MdClick Successfully deleted")
}

func UpdateMdClick(c *fiber.Ctx) error {
	mdClick := new(MdClick)
	if err := c.BodyParser(mdClick); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	id, _ := strconv.Atoi(c.Params("id"))
	db := database.DBConn
	db.Model(&MdClick{}).Where("id = ?", id).Updates(map[string]interface{}{
		"ref":   mdClick.Ref,
		"score": mdClick.Score,
	})
	return c.SendString("MdClick Successfully updated")
}
