package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

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

func (*server) PrimeNumberDecomposition(req *calculatorpb.CalculateManyTimesRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("PrimeNumberDecomposition was invoked with %v\n", req)
	numberOne := req.GetCalculator().GetNumberOne()

	divisor := int32(2)
	number := numberOne

	for number > 1 {
		if number%divisor == 0 {
			res := &calculatorpb.CalculateManyTimesResponse{
				Result: divisor,
			}
			stream.Send(res)
			time.Sleep(1000 * time.Millisecond)
			number = number / divisor

		} else {
			divisor = divisor + 1
		}
	}
	return nil

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
