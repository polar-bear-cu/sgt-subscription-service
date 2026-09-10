package grpc

import (
	"context"

	subscriptionv1 "github.com/polar-bear-cu/sgt-proto/gen/go/subscription/v1"
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

	out := make([]*subscriptionv1.Subscription, 0, len(subs))
	for _, sub := range subs {
		out = append(out, &subscriptionv1.Subscription{
			Id:     sub.ID,
			UserId: sub.UserID,
		})
	}
	return &subscriptionv1.GetSubscriptionsForReportResponse{Subscription: out}, nil
}
