package mocks

import (
	"github.com/golang/mock/gomock"
	"github.com/mauhlik/go-index/internal/go-index/providers"
)

type MockProvider struct {
	providers.Provider
	mock *gomock.Controller
}

func NewMockProvider(ctrl *gomock.Controller) *MockProvider {
	return &MockProvider{
		Provider: nil,
		mock:     ctrl,
	}
}

func (m *MockProvider) GetVersions(_, _ string) ([]string, error) {
	return []string{"0.0.0", "0.0.1", "1.0.0", "2.0.0"}, nil
}
