package internalgrpc

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/app"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func convertGet(e app.Event) (*GetEventResponse, error) {
	date, err := ptypes.TimestampProto(e.Date)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to convert event")
	}
	item := &GetEventResponse{
		Id:      string(e.ID),
		Title:   e.Title,
		Date:    date,
		Latency: int64(e.Latency),
		Note:    e.Note,
		UserId:  e.UserID,
		Notify:  int64(e.Notify),
	}
	return item, nil
}
func convertUpdate(e app.Event) (*UpdateEventResponse, error) {
	date, err := ptypes.TimestampProto(e.Date)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to convert event")
	}
	item := &UpdateEventResponse{
		Id:      string(e.ID),
		Title:   e.Title,
		Date:    date,
		Latency: int64(e.Latency),
		Note:    e.Note,
		UserId:  e.UserID,
		Notify:  int64(e.Notify),
	}
	return item, nil
}

func convertList(events []app.Event) (*ListEventResponse, error) {
	items := make([]*ListEventItem, len(events))
	var err error
	var i int
	for _, e := range events {
		date, err := ptypes.TimestampProto(e.Date)
		if err != nil {
			return nil, status.Error(codes.Internal, "unable to convert event")
		}
		evt := &ListEventItem{
			Id:      string(e.ID),
			Title:   e.Title,
			Date:    date,
			Latency: int64(e.Latency),
			Note:    e.Note,
			UserId:  e.UserID,
			Notify:  int64(e.Notify),
		}
		items[i] = evt
		i++
	}
	return &ListEventResponse{
		Results: items,
	}, err
}

func convertCreate(e app.EventID) (*CreateEventResponse, error) {
	return &CreateEventResponse{
		Id: string(e),
	}, nil
}
