package service

import (
	pb "cloud-cost-optimizer/gen/cost/v1"
)

type AlertService struct {
}

func NewAlertService() *AlertService {
	return &AlertService{}
}

func (service *AlertService) ProcessBillingData(account *pb.Account) *ProcessBillingDataRes {
	return nil
}

func (service *AlertService) ProcessCostTrends(account *pb.Account) *ProcessCostTrendsRes{
	return nil
}

func (service *AlertService) ProcessRecommendation(account *pb.Account) *ProcessRecommendationRes{
	return nil
}

func (service *AlertService) ProcessWatchCostAlerts(account *pb.Account) *ProcessWatchCostAlertsRes{
	return nil
}