package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"../calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Calculator(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	fmt.Printf("Calculator function was invoked with %v \n", req)
	calculator := req.GetCalculator()

	result := calculator.GetNumberOne() + calculator.GetNumberTwo()

	res := &calculatorpb.CalculatorResponse{
		Result: result,
	}

	return res, nil
}

func main() {
	fmt.Println("Initialize Calculator")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
