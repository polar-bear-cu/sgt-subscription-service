package grpc

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *SubscriptionServer) ListDueReminders(
	ctx context.Context,
	req *subscriptionv1.ListDueRemindersRequest,
) (*subscriptionv1.ListDueRemindersResponse, error) {
	if err := validateDate(req.GetDate()); err != nil {
		return nil, err
	}

	due, err := s.uc.ListDueReminders(ctx, req.GetDate())
	if err != nil {
		return nil, err
	}

	out := make([]*subscriptionv1.DueReminder, 0, len(due))
	for _, d := range due {
		out = append(out, &subscriptionv1.DueReminder{
			Subscription: toProto(d.Subscription),
			BillingDue:   d.BillingDue,
			TrialEndDue:  d.TrialEndDue,
		})
	}
	return &subscriptionv1.ListDueRemindersResponse{Reminders: out}, nil
}

func (s *SubscriptionServer) AdvanceBillingDates(
	ctx context.Context,
	req *subscriptionv1.AdvanceBillingDatesRequest,
) (*subscriptionv1.AdvanceBillingDatesResponse, error) {
	if err := validateDate(req.GetDate()); err != nil {
		return nil, err
	}

	advanced, converted, err := s.uc.AdvanceBillingDates(ctx, req.GetDate())
	if err != nil {
		log.Printf("advance billing dates %s: advanced=%d converted=%d err=%v", req.GetDate(), advanced, converted, err)
		return nil, err
	}
	return &subscriptionv1.AdvanceBillingDatesResponse{
		AdvancedCount:        advanced,
		TrialsConvertedCount: converted,
	}, nil
}

func (s *SubscriptionServer) DeleteSubscriptionsByUser(
	ctx context.Context,
	req *subscriptionv1.DeleteSubscriptionsByUserRequest,
) (*subscriptionv1.DeleteSubscriptionsByUserResponse, error) {
	deleted, err := s.uc.DeleteByUser(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	return &subscriptionv1.DeleteSubscriptionsByUserResponse{DeletedCount: deleted}, nil
}

func validateDate(date string) error {
	if _, err := time.Parse(time.DateOnly, date); err != nil {
		return status.Error(codes.InvalidArgument, "date must be YYYY-MM-DD")
	}
	return nil
}

func toProto(sub models.Subscription) *subscriptionv1.Subscription {
	p := &subscriptionv1.Subscription{
		Id:                     sub.ID,
		UserId:                 sub.UserID,
		Name:                   sub.Name,
		BillingDate:            timestamppb.New(sub.NextBillingDate),
		Cost:                   sub.Cost,
		Type:                   sub.Type,
		Category:               sub.Category,
		Status:                 sub.Status,
		ReminderTimeInAdvanced: sub.ReminderTimeInAdvanced,
	}
	if sub.FtEndDate != nil {
		p.FtEndDate = timestamppb.New(*sub.FtEndDate)
	}
	return p
}

func toProtoList(subs []models.Subscription) []*subscriptionv1.Subscription {
	out := make([]*subscriptionv1.Subscription, 0, len(subs))
	for _, sub := range subs {
		out = append(out, toProto(sub))
	}
	return out
}
