package models

type ControlBy struct {
	Creator uint `gorm:"index;comment:'创建者'"`
	Updator uint `gorm:"index;comment:'更新者'"`
}

func (e *ControlBy) SetCreator(Creator uint) {
	e.Creator = Creator
}

func (e *ControlBy) SetUpdator(updator uint) {
	e.Updator = updator
}
