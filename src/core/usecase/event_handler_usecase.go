package usecase

type IEventHandlerUsecase interface {
}
type EventHandlerUsecase struct {
}

func NewEventHandlerUsecase() IEventHandlerUsecase {
	return &EventHandlerUsecase{}
}
