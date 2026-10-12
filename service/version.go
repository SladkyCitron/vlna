package service

import "github.com/SladkyCitron/vlna/version"

type VersionService struct{}

func NewVersionService() *VersionService {
	return &VersionService{}
}

func (s *VersionService) GetVersion() string {
	return version.Version
}
