package service

import (
	pb "cloud-cost-optimizer/gen/cost/v1"
	"context"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AWSService struct {
	ceClient *costexplorer.Client
	ctx      context.Context
}

func NewAWSServiceWithKeys(accessKey, secretKey, region string) (*costexplorer.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKey, secretKey, "",
		)))

	if err != nil {
		print(err)
		return nil, err
	}
	return costexplorer.NewFromConfig(cfg), err
}

func NewAWSService(ctx context.Context) *AWSService {
	accessKey, secretKey, region := "", "", "us-east-1"
	ceClient, err := NewAWSServiceWithKeys(accessKey, secretKey, region)
	if err != nil {
		panic(err)
	}
	return &AWSService{
		ceClient: ceClient,
		ctx:      ctx,
	}
}

// tag cost usagedate 根据tag来划分
func (service *AWSService) ProcessBillingData(account *pb.Account, timerange *pb.TimeRange) ([]*ProcessBillingDataRes, error) {
	var tags []string
	var groupDefinitionLists []types.GroupDefinition
	var tag string
	var records []*ProcessBillingDataRes

	timePeriod := &types.DateInterval{
		Start: aws.String(timerange.StartTime.AsTime().Format("2006-01-02")), //需要先指定格式转成字符串
		End:   aws.String(timerange.EndTime.AsTime().Format("2006-01-02")),
	}
	tagsResult, err := service.ceClient.GetTags(service.ctx, &costexplorer.GetTagsInput{
		TimePeriod: timePeriod,
	})
	if err != nil {
		return nil, err
	}
	tags = tagsResult.Tags

	for _, tag := range tags {
		groupDefinitionLists = append(groupDefinitionLists, types.GroupDefinition{
			Type: types.GroupDefinitionTypeTag,
			Key:  aws.String(tag),
		})
	}
	costAndUsageData, err := service.ceClient.GetCostAndUsage(service.ctx, &costexplorer.GetCostAndUsageInput{
		Granularity: types.GranularityDaily,
		TimePeriod:  timePeriod,
		Metrics:     []string{"UnblendedCost"}, //UnblendedCost 原始账单费用 AmortizedCost分摊成本
		GroupBy:     groupDefinitionLists,
	})
	if err != nil {
		return nil, err
	}

	if len(costAndUsageData.ResultsByTime) == 0 {
		return records, nil
	}
	for _, result := range costAndUsageData.ResultsByTime {
		for _, metric := range result.Groups {

			//防止长度不够2
			if len(metric.Keys) < 2 {
				tag = metric.Keys[0]
			} else {
				tag = metric.Keys[1]
			}
			cost, _ := strconv.ParseFloat(*metric.Metrics["UnblendedCost"].Amount, 64)
			usageDate, err := time.Parse("2006-01-02", *result.TimePeriod.Start)
			if err != nil {
				return nil, err
			}
			usageDateTimestamp := timestamppb.New(usageDate)
			records = append(records, NewProcessBillingDataRes(tag, cost, usageDateTimestamp))
		}

	}
	return records, nil
}

func (service *AWSService) ProcessCostTrends(account *pb.Account, tag string) (*ProcessCostTrendsRes, error) {
	var records *ProcessCostTrendsRes
	today := timestamppb.New(time.Now())
	endDate := timestamppb.New(today.AsTime().AddDate(0, 30, 0))
	timePeriod := &types.DateInterval{
		Start: aws.String(today.AsTime().Format("2006-01-02")), //需要先指定格式转成字符串
		End:   aws.String(endDate.AsTime().Format("2006-01-02")),
	}
	costForecastData, err := service.ceClient.GetCostForecast(service.ctx, &costexplorer.GetCostForecastInput{
		Granularity: types.GranularityDaily,
		TimePeriod:  timePeriod,
		Metric:      types.MetricUnblendedCost,
		Filter: &types.Expression{
			Tags: &types.TagValues{
				Key: aws.String(tag),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	for _, record := range costForecastData.ForecastResultsByTime {
		cost, _ := strconv.ParseFloat(*record.Total.Amount, 64)
		records = NewProcessCostTrendsRes(cost)
	}
	return records, nil
}

func (service *AWSService) ProcessRecommendation(account *pb.Account) (*ProcessRecommendationRes, error) {
	return nil
}

func (service *AWSService) ProcessWatchCostAlerts(account *pb.Account, costThreshold float64) (*ProcessWatchCostAlertsRes, error) {
	return nil
}
