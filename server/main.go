package main

import (
	"cloud-cost-optimizer/server/handler"
	"cloud-cost-optimizer/server/service"
	"context"
	"fmt"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	azureService := service.NewAzureService()
	awsService := service.NewAWSService()
	alertService := service.NewAlertService()
	costOptimizerHandler := handler.NewCostOptimizerHandler(ctx, awsService, azureService, alertService)
	costOptimizerHandler.StartServe()
	fmt.Println("gRPC服务器开始服务")
}
