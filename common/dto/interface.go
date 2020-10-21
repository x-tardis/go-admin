package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/thinkgos/sharp/core/paginator"
	"gorm.io/gorm/schema"
)

type Index interface {
	Generate() Index
	Bind(ctx *gin.Context) error
	GetPaginatorParam() paginator.Param
	GetNeedSearch() interface{}
}

type Control interface {
	Generate() Control
	Bind(ctx *gin.Context) error
	GenerateM() (ActiveRecord, error)
	GetId() interface{}
}

type ActiveRecord interface {
	schema.Tabler
	SetCreator(createBy uint)
	SetUpdator(updateBy uint)
	Generate() ActiveRecord
	GetId() uint
}
