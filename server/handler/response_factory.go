package handler

import (
	pb "cloud-cost-optimizer/gen/cost/v1"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type ResponseFactory struct{}

func NewResponseFactory() *ResponseFactory {
	return &ResponseFactory{}
}

func (responseFactory *ResponseFactory) CreateStreamBillingDataResponse(provider string, accountId string, tag string, cost float64, usageDate *timestamppb.Timestamp) (*pb.StreamBillingDataResponse, error) {
	if provider == "" {
		return nil, fmt.Errorf("invalid input: provider cannot be empty")
	}
	if accountId == "" || tag == "" {
		return nil, fmt.Errorf("Provider: %s, invalid input", provider)
	}
	if cost < 0 {
		return nil, fmt.Errorf("Provider: %s,invalid input: cost cannot be negative", provider)
	}
	if usageDate == nil {
		return nil, fmt.Errorf("Provider: %s, invalid input: usageDate cannot be nil", provider)
	}
	return &pb.StreamBillingDataResponse{
		Provider:  provider,
		AccountId: accountId,
		Tag:       tag,
		Cost:      cost,
		UsageDate: usageDate,
	}, nil
}

func (responseFactory *ResponseFactory) CreateBatchGetCostTrendsResponse(trend *pb.CostTrend) (*pb.BatchGetCostTrendsResponse, error) {
	if trend == nil {
		return nil, fmt.Errorf("trends cannot be nil")
	}
	return &pb.BatchGetCostTrendsResponse{
		Trend: trend,
	}, nil
}

func (responseFactory *ResponseFactory) CreateGetRecommendationsResponse(recommendations []*pb.Recommendation) (*pb.GetRecommendationsResponse, error) {
	if recommendations == nil {
		return nil, fmt.Errorf("recommendations cannot be nil")
	}
	return &pb.GetRecommendationsResponse{
		Recommendations: recommendations,
	}, nil
}

func (responseFactory *ResponseFactory) CreateWatchCostAlertsResponse(title string, description string, current_cost float64, alert_time *timestamppb.Timestamp) (*pb.WatchCostAlertsResponse, error) {
	if title == "" {
		return nil, fmt.Errorf("title cannot be empty")
	}
	if description == "" {
		return nil, fmt.Errorf("description cannot be empty")
	}
	if current_cost < 0 {
		return nil, fmt.Errorf("current_cost cannot be negative")
	}
	if alert_time == nil {
		return nil, fmt.Errorf("alert_time cannot be nil")
	}
	return &pb.WatchCostAlertsResponse{
		Title:       title,
		Description: description,
		CurrentCost: current_cost,
		AlertTime:   alert_time,
	}, nil
}
