package services

func GetSectionService() SectionServiceI {
	return &SectionService{}
}

type SectionServiceI interface{}

type SectionService struct{}
