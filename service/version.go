package service

import (
	"github.com/SladkyCitron/vlna/version"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type VersionService struct{}

func NewVersionService() *VersionService {
	return &VersionService{}
}

func (s *VersionService) GetVersion() string {
	return version.Version
}

func (s *VersionService) OpenGitHub() error {
	return application.Get().Browser.OpenURL("https://github.com/SladkyCitron/vlna")
}
