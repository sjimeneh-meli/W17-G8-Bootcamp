package tests

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/stretchr/testify/mock"
)

func GetSectionValidationMock() *SectionValidationMock {
	return &SectionValidationMock{}
}

type SectionValidationMock struct {
	mock.Mock
}

func (m *SectionValidationMock) ValidateSectionRequestStruct(r requests.SectionRequest) error {
	args := m.Called(r)
	return args.Error(0)
}
