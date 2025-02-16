package info

type InfoService struct {
}

func NewInfoService() *InfoService {
	return &InfoService{}
}

func (service InfoService) GetInfoForUser(username string) string {
	return "This is a demo service"
}
