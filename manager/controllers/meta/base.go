package meta

import "gorm.io/gorm"

type Controller struct {
	Conn *gorm.DB
}
