package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

// 第三个info参数表示当前是对应的那个gRPC方法
// 第四个handler参数对应当前的gRPC方法函数
// 首先是日志输出info参数，然后调用handler对应的gRPC方法函数

func filter(
	ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	log.Println("filter:", info)

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	return handler(ctx, req)
}
