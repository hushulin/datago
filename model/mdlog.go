package model

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"main/database"
	"main/utils"
	"net/url"
	"strconv"
)

// MdLog mdLog mdLogs 日志
type MdLog struct {
	gorm.Model
	Title   string `json:"title"`
	Content string `json:"content"`
}

func GetMdLogs(c *fiber.Ctx) error {
	db := database.DBConn
	var mdLogs []MdLog
	db.Find(&mdLogs)
	return c.JSON(mdLogs)
}

func GetMdLogs5(c *fiber.Ctx) error {
	db := database.DBConn
	var mdLogs []MdLog
	var mdClicks []MdClick
	if err := db.Model(&MdClick{}).Select("markdown_id, SUM(score) AS count").Group("markdown_id").Order("count desc").Limit(5).Find(&mdClicks).Error; err != nil {
		return c.Status(503).SendString(err.Error())
	}
	var markdownIDs []uint
	for _, click := range mdClicks {
		markdownIDs = append(markdownIDs, click.MarkdownId)
	}
	db.Where("id IN (?)", markdownIDs).Find(&mdLogs)
	return c.JSON(mdLogs)
}

func GetMdLogByPage(p *utils.Pagination) (mdLogs []MdLog, err error) {
	db := database.DBConn
	err = db.Model(&MdLog{}).Scopes(p.GormPaginate()).Order("id desc").Find(&mdLogs).Error
	if err != nil {
		return nil, err
	}
	var total int64
	db.Model(&MdLog{}).Count(&total)
	p.Total = cast.ToInt(total)
	return
}

func GetMdLog(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var mdLog MdLog
	db.Find(&mdLog, id)
	return c.JSON(mdLog)
}

func GetMdLogByKeyword(c *fiber.Ctx) error {
	keyword := c.Params("keyword")
	unescape, err := url.QueryUnescape(keyword)
	if err != nil {
		return err
	}
	db := database.DBConn
	var mdLogs []MdLog
	db.Where("`title` LIKE ? OR `content` LIKE ?", "%"+unescape+"%", "%"+unescape+"%").Find(&mdLogs)
	return c.JSON(mdLogs)
}

func GetMdLogLatest100(c *fiber.Ctx) error {
	db := database.DBConn
	var mdLogs []MdLog
	db.Select("ID,title").Order("id desc").Limit(100).Find(&mdLogs)
	return c.JSON(mdLogs)
}

func GetMdLog2(c *fiber.Ctx) MdLog {
	id := c.Params("id")
	db := database.DBConn
	var mdLog MdLog
	db.Find(&mdLog, id)
	return mdLog
}

func NewMdLog(c *fiber.Ctx) error {
	db := database.DBConn
	mdLog := new(MdLog)
	if err := c.BodyParser(mdLog); err != nil {
		return c.Status(503).SendString(err.Error())
	}
	db.Create(&mdLog)
	//"/markdown-page/" + strconv.Itoa(int(mdLog.ID))
	return c.Redirect(fmt.Sprintf("/markdown-page/%d", mdLog.ID))
}

func DeleteMdLog(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn

	var mdLog MdLog
	db.First(&mdLog, id)
	db.Delete(&mdLog)
	return c.SendString("MdLog Successfully deleted")
}

func UpdateMdLog(c *fiber.Ctx) error {
	mdLog := new(MdLog)
	if err := c.BodyParser(mdLog); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	id, _ := strconv.Atoi(c.Params("id"))
	db := database.DBConn
	db.Model(&MdLog{}).Where("id = ?", id).Updates(map[string]interface{}{
		"title":   mdLog.Title,
		"content": mdLog.Content,
	})
	//return c.SendString("MdLog Successfully updated")
	return c.RedirectBack("更新成功")
}
