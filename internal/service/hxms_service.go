package service

import "go-scheduler/config"

type HXMSService interface {
	GetUser()
}

type hxmsService struct {
	config *config.Config
}

func NewHXMSService(config *config.Config) HXMSService {
	return &hxmsService{config: config}
}

func (hs *hxmsService) GetUser() {
	// do something
}
