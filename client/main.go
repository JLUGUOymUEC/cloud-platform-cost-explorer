package handler

import (
	"context"
	"fmt"
)

type Handler struct{
	ctx context.Context
}

func main() {
	context, cancel := context.WithCancel(context.Background())
	defer cancel()
	handler := &Handler{
		ctx: context,
	}
	fmt.Println("客户端结束")
}
