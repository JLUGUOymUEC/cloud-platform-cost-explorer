package service

import (
	pb "cloud-cost-optimizer/gen/cost/v1"
)

type AWSService struct {
}

func NewAWSService() *AWSService{
	return &AWSService{}
}

func (service *AWSService) ProcessBillingData(account *pb.Account) *ProcessBillingDataRes {
	return nil
}

func (service *AWSService) ProcessCostTrends(account *pb.Account) *ProcessCostTrendsRes{
	return nil
}

func (service *AWSService) ProcessRecommendation(account *pb.Account) *ProcessRecommendationRes{
	return nil
}

func (service *AWSService) ProcessWatchCostAlerts(account *pb.Account) *ProcessWatchCostAlertsRes{
	return nil
}