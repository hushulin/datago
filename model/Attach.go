package model

import (
	"gorm.io/gorm"
	"main/database"
)

type Attachment struct {
	gorm.Model
	FilePath  string `json:"file_path"`
	FileName  string `json:"file_name"`
	Extension string `json:"extension"`
	MimeType  string `json:"mime_type"`
	Size      int    `json:"size"`
}

// NewAttachmentForInner 第一个要使用的接口 内部使用
func NewAttachmentForInner(filePath string) *Attachment {
	db := database.DBConn
	attachment := new(Attachment)
	attachment.FilePath = filePath
	db.Create(&attachment)
	return attachment
}
