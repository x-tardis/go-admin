package models

import "gorm.io/gorm/schema"

type ActiveRecord interface {
	schema.Tabler
	SetCreator(createBy uint)
	SetUpdator(updateBy uint)
	Generate() ActiveRecord
	GetId() uint
}
