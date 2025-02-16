package info

type Handler struct {
	InfoService *InfoService
}

func NewHandler() *Handler {
	return &Handler{}
}

func (handler Handler) GetInfoForUser(username string) string {
	return handler.InfoService.GetInfoForUser(username)
}
