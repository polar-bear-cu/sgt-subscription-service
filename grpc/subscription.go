package grpc

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	subscriptionv1 "github.com/polar-bear-cu/sgt-proto/gen/go/subscription/v1"
	"github.com/polar-bear-cu/sgt-subscription-service/models"
	"github.com/polar-bear-cu/sgt-subscription-service/usecases"
)

type SubscriptionServer struct {
	subscriptionv1.UnimplementedSubscriptionServiceServer
	uc *usecases.SubscriptionUsecase
}

func NewSubscriptionServer(uc *usecases.SubscriptionUsecase) *SubscriptionServer {
	return &SubscriptionServer{uc: uc}
}

func (s *SubscriptionServer) GetSubscriptionsForReport(
	ctx context.Context,
	req *subscriptionv1.GetSubscriptionsForReportRequest,
) (*subscriptionv1.GetSubscriptionsForReportResponse, error) {
	subs, err := s.uc.ListByUser(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	return &subscriptionv1.GetSubscriptionsForReportResponse{Subscription: toProtoList(subs)}, nil
}

func (s *SubscriptionServer) GetUpcomingForBilling(
	ctx context.Context,
	req *subscriptionv1.GetUpcomingForBillingRequest,
) (*subscriptionv1.GetUpcomingForBillingResponse, error) {
	within := time.Duration(req.GetWithinHours()) * time.Hour

	subs, err := s.uc.GetUpcomingForBilling(ctx, within)
	if err != nil {
		return nil, err
	}

	return &subscriptionv1.GetUpcomingForBillingResponse{Subscription: toProtoList(subs)}, nil
}

func toProtoList(subs []models.Subscription) []*subscriptionv1.Subscription {
	out := make([]*subscriptionv1.Subscription, 0, len(subs))
	for _, sub := range subs {
		out = append(out, &subscriptionv1.Subscription{
			Id:          sub.ID,
			UserId:      sub.UserID,
			Name:        sub.Name,
			BillingDate: timestamppb.New(sub.BillingDate),
		})
	}
	return out
}
