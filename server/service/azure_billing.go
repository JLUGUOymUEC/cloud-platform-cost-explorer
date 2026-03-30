package service

import (
	pb "cloud-cost-optimizer/gen/cost/v1"
)

type AzureService struct{

}

func NewAzureService() *AzureService {
	return &AzureService{}
}

func (service *AzureService) ProcessBillingData(account *pb.Account) *ProcessBillingDataRes {
	return nil
}

func (service *AzureService) ProcessCostTrends(account *pb.Account) *ProcessCostTrendsRes{
	return nil
}

func (service *AzureService) ProcessRecommendation(account *pb.Account) *ProcessRecommendationRes{
	return nil
}

func (service *AzureService) ProcessWatchCostAlerts(account *pb.Account) *ProcessWatchCostAlertsRes{
	return nil
}