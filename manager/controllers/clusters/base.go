package clusters

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Controller struct {
	Conn *gorm.DB
}

var logger *logrus.Logger
