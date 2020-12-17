package internalgrpc

import (
	"context"
	"time"

	"github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/app"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AdapterService struct {
	UnimplementedEventsServer
	app app.Application
}

func (a AdapterService) GetEvent(ctx context.Context, request *GetEventRequest) (*GetEventResponse, error) {
	e, err := a.app.GetEvent(ctx, app.EventID(request.Id))
	if err != nil {
		return nil, err
	}
	return convertGet(e)
}

func (a AdapterService) GetEventList(ctx context.Context, _ *emptypb.Empty) (*ListEventResponse, error) {
	e, err := a.app.GetEventList(ctx)
	if err != nil {
		return nil, err
	}
	return convertList(e)
}

func (a AdapterService) CreateEvent(ctx context.Context, request *CreateEventRequest) (*CreateEventResponse, error) {
	e, err := a.app.CreateEvent(ctx, request.Title, request.Date.AsTime(), time.Duration(request.Latency), request.Note, request.UserId, time.Duration(request.Notify))
	if err != nil {
		return nil, err
	}
	return convertCreate(e)
}

func (a AdapterService) UpdateEvent(ctx context.Context, request *UpdateEventRequest) (*UpdateEventResponse, error) {
	e, err := a.app.UpdateEvent(ctx, app.EventID(request.Id), request.Title, request.Date.AsTime(), time.Duration(request.Latency), request.Note, request.UserId, time.Duration(request.Notify))
	if err != nil {
		return nil, err
	}
	return convertUpdate(e)
}

func (a AdapterService) DeleteEvent(ctx context.Context, request *DeleteEventRequest) (*emptypb.Empty, error) {
	id := request.Id
	err := a.app.DeleteEvent(ctx, app.EventID(id))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func NewInternalAdapter(app app.Application) AdapterService {
	return AdapterService{
		app: app,
	}
}
