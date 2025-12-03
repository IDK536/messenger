package controller

import (
	"messanger/internal/config"
)

type Controller interface {
	Start(cfg *config.Config)
}
