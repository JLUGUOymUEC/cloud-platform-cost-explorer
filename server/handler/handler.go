package handler

import (
	pb "cloud-cost-optimizer/gen/cost/v1"
	"cloud-cost-optimizer/server/service"
	"context"
	"fmt"
	"sync"
)

type CostOptimizerHandler struct {
	ctx             context.Context
	cloudProviders  map[string]service.CloudProvider
	responseFactory *ResponseFactory
	clients         map[string]pb.CostOptimizationService_HandleServer
	wg				sync.WaitGroup
}

func NewCostOptimizerHandler(ctx context.Context, awsService *service.AWSService, azureService *service.AzureService, alertService *service.AlertService) *CostOptimizerHandler {
	cloudProviders := map[string]service.CloudProvider{
		"AWS":   awsService,
		"Azure": azureService,
		"Alert": alertService,
	}
	responseFactory := NewResponseFactory()
	return &CostOptimizerHandler{
		ctx:             ctx,
		cloudProviders:  cloudProviders,
		responseFactory: responseFactory,
	}
}

func (handler *CostOptimizerHandler) Handle(stream pb.CostOptimizationService_HandleServer) error {	
	errChan := make(chan error, 1)
	go func(){
		select{
			case <-handler.ctx.Done():
				return 
			default:
		}
	}()
	for {

			//一直不发消息就一直阻塞所以要用goroutine
			req, err := stream.Recv()
			
			if err != nil {
				fmt.Printf("Error receiving request: %v\n", err)
				continue
			}
			switch req.RequestType.(type) {
			case *pb.HandleRequest_StreamBillingDataRequest:
				// 处理来自AWS和Azure的账单数据请求
				err = handler.handleStreamBillingDataRequest(req, stream)
			case *pb.HandleRequest_BatchGetCostTrendsRequest:
				// 处理获取成本趋势的请求
				err = handler.handleBatchGetCostTrendsRequest(req, stream)
			case *pb.HandleRequest_GetRecommendationsRequest:
				// 处理获取优化建议的请求
				err = handler.handleGetRecommendationsRequest(req, stream)
			case *pb.HandleRequest_WatchCostAlertsRequest:
				// 处理监控成本警报的请求
				err = handler.handleWatchCostAlertsRequest(req, stream)
			}
			if err != nil {
				fmt.Printf("Error handling request: %v\n", err)
				errChan <- err
			}

		}
	
			
}

func (handler *CostOptimizerHandler) handleStreamBillingDataRequest(req *pb.HandleRequest, stream pb.CostOptimizationService_HandleServer) error {
	streamBillingDatarequest := req.GetStreamBillingDataRequest()
	// 处理来自AWS和Azure的账单数据请求
	Accounts := streamBillingDatarequest.GetAccounts()
	for _, account := range Accounts {
		// 处理每个账户的账单数据
		var res *service.ProcessBillingDataRes
		if _, ok := handler.cloudProviders[account.GetProvider()]; !ok {
			res = service.NewProcessBillingDataRes("Unknown", 0.0, nil)

		} else {
			res = handler.cloudProviders[account.GetProvider()].ProcessBillingData(account)
		}

		// 将处理结果发送回客户端
		response, err := handler.responseFactory.CreateStreamBillingDataResponse(
			account.GetProvider(),
			account.GetAccountId(),
			res.GetService(),
			res.GetCost(),
			res.GetUsageDate(),
		)
		if err != nil {
			fmt.Printf("Error creating StreamBillingDataRequest: %v\n", err)
			return err
		}
		stream.Send(&pb.HandleResponse{
			ResponseType: &pb.HandleResponse_StreamBillingDataResponse{
				StreamBillingDataResponse: response,
			},
		})
	}
	return nil
}

func (handler *CostOptimizerHandler) handleBatchGetCostTrendsRequest(req *pb.HandleRequest, stream pb.CostOptimizationService_HandleServer) error {
	batchGetCostTrendsRequest := req.GetBatchGetCostTrendsRequest()
	// 处理批量获取成本趋势的请求
	Account := batchGetCostTrendsRequest.GetAccount()
	var res *service.ProcessCostTrendsRes
	if _, ok := handler.cloudProviders[Account.GetProvider()]; !ok {
		res = service.NewProcessCostTrendsRes([]*pb.CostTrend{
			{
				Service:    "Unknown",
				DailyCosts: nil,
			},
		})
	} else {
		res = handler.cloudProviders[Account.GetProvider()].ProcessCostTrends(Account)
	}
	response, err := handler.responseFactory.CreateBatchGetCostTrendsResponse(
		res.GetCostTrends(),
	)
	if err != nil {
		fmt.Printf("Error creating BatchGetCostTrendsResponse: %v\n", err)
		return err
	}
	stream.Send(&pb.HandleResponse{
		ResponseType: &pb.HandleResponse_BatchGetCostTrendsResponse{
			BatchGetCostTrendsResponse: response,
		},
	})
	return nil
}

func (handler *CostOptimizerHandler) handleGetRecommendationsRequest(req *pb.HandleRequest, stream pb.CostOptimizationService_HandleServer) error {
	getRecommendationsRequest := req.GetGetRecommendationsRequest()
	// 处理获取优化建议的请求
	Account := getRecommendationsRequest.GetAccount()
	var res *service.ProcessRecommendationRes
	if _, ok := handler.cloudProviders[Account.GetProvider()]; !ok {
		res = service.NewProcessRecommendationsRes([]*pb.Recommendation{
			{
				Title:            "Unknown",
				Description:      "Unknown",
				EstimatedSavings: 0.0,
			},
		})
	} else {
		res = handler.cloudProviders[Account.GetProvider()].ProcessRecommendation(Account)
	}
	response, err := handler.responseFactory.CreateGetRecommendationsResponse(
		res.GetRecommendations(),
	)
	if err != nil {
		fmt.Printf("Error creating GetRecommendationsResponse: %v\n", err)
		return err
	}
	stream.Send(&pb.HandleResponse{
			ResponseType: &pb.HandleResponse_GetRecommendationsResponse{
				GetRecommendationsResponse: response,
			},
		})
	return nil
}

func (handler *CostOptimizerHandler) handleWatchCostAlertsRequest(req *pb.HandleRequest, stream pb.CostOptimizationService_HandleServer) error {
	watchCostAlertsRequest := req.GetWatchCostAlertsRequest()
	// 处理监控成本警报的请求
	Accounts := watchCostAlertsRequest.GetAccounts()
	for _, account := range Accounts {
		var res *service.ProcessWatchCostAlertsRes
		if _, ok := handler.cloudProviders[account.GetProvider()]; !ok {
			res = service.NewProcessWatchCostAlertsRes("Unknown", "Unknown", 0.0, nil)
		} else {
			res = handler.cloudProviders[account.GetProvider()].ProcessWatchCostAlerts(account)
		}
		response, err := handler.responseFactory.CreateWatchCostAlertsResponse(
			res.GetTitle(),
			res.GetDescription(),
			res.GetCurrentCost(),
			res.GetAlertTime(),
		)
		if err != nil {
			fmt.Printf("Error creating WatchCostAlertsResponse: %v\n", err)
			return err
		}
		stream.Send(
			&pb.HandleResponse{
				ResponseType: &pb.HandleResponse_WatchCostAlertsResponse{
					WatchCostAlertsResponse: response,
				},
			})
	}
	return nil
}


func (handler *CostOptimizerHandler) GetClients() map[string]pb.CostOptimizationService_HandleServer {
	return handler.clients
}

func (handler *CostOptimizerHandler) AddClient(clientID string, stream pb.CostOptimizationService_HandleServer) {
	handler.clients[clientID] = stream
}

func (handler *CostOptimizerHandler) RemoveClient(clientID string) {
	delete(handler.clients, clientID)
}

func (handler *CostOptimizerHandler) StartServe(){
	for _, client :=range handler.clients{
		handler.wg.Add(1)
		go func(client pb.CostOptimizationService_HandleServer){
			defer handler.wg.Done()
			handler.Handle(client)
		}(client)
	}
	handler.wg.Wait()
}