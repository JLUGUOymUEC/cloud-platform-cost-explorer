package service

import (
	pb "cloud-cost-optimizer/gen/cost/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type CloudProvider interface {
	//按日获取aws azure alert账单
	ProcessBillingData(account *pb.Account, timerange *pb.TimeRange) ([]*ProcessBillingDataRes, error)
	//获取aws azure alert账单趋势
	ProcessCostTrends(account *pb.Account, service string) (*ProcessCostTrendsRes, error)
	//获取aws azure alert计划推荐
	ProcessRecommendation(account *pb.Account) (*ProcessRecommendationRes, error)
	//获取aws azure alert监控
	ProcessWatchCostAlerts(account *pb.Account, costThreshold float64) (*ProcessWatchCostAlertsRes, error)
}

// 返回值结构体类型
type ProcessBillingDataRes struct {
	tag       string
	cost      float64
	usageDate *timestamppb.Timestamp
}

type ProcessCostTrendsRes struct {
	trend pb.CostTrend
}

type ProcessRecommendationRes struct {
	recommendations []*pb.Recommendation
}

type ProcessWatchCostAlertsRes struct {
	title        string
	description  string
	current_cost float64
	alertTime    *timestamppb.Timestamp
}

// 返回值结构体获取方法
func (res *ProcessBillingDataRes) GetTag() string {
	return res.tag
}

func (res *ProcessBillingDataRes) GetCost() float64 {
	return res.cost
}

func (res *ProcessBillingDataRes) GetUsageDate() *timestamppb.Timestamp {
	return res.usageDate
}

func (res *ProcessCostTrendsRes) GetCostTrend() *pb.CostTrend {
	return &res.trend
}

func (res *ProcessRecommendationRes) GetRecommendations() []*pb.Recommendation {
	return res.recommendations
}

func (res *ProcessWatchCostAlertsRes) GetTitle() string {
	return res.title
}

func (res *ProcessWatchCostAlertsRes) GetDescription() string {
	return res.description
}

func (res *ProcessWatchCostAlertsRes) GetCurrentCost() float64 {
	return res.current_cost
}

func (res *ProcessWatchCostAlertsRes) GetAlertTime() *timestamppb.Timestamp {
	return res.alertTime
}

// 结构体构造函数
func NewProcessBillingDataRes(tag string, cost float64, usageDate *timestamppb.Timestamp) *ProcessBillingDataRes {
	return &ProcessBillingDataRes{
		tag:       tag,
		cost:      cost,
		usageDate: usageDate,
	}
}

func NewProcessCostTrendsRes(costTrend *pb.CostTrend) *ProcessCostTrendsRes {
	return &ProcessCostTrendsRes{
		trend: *costTrend,
	}
}

func NewProcessRecommendationsRes(recommendations []*pb.Recommendation) *ProcessRecommendationRes {
	return &ProcessRecommendationRes{
		recommendations: recommendations,
	}
}

func NewProcessWatchCostAlertsRes(title string, description string, current_cost float64, alertTime *timestamppb.Timestamp) *ProcessWatchCostAlertsRes {
	return &ProcessWatchCostAlertsRes{
		title:        title,
		description:  description,
		current_cost: current_cost,
		alertTime:    alertTime,
	}
}
