package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"../calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Client initialized")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	defer cc.Close()
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	c := calculatorpb.NewCalculatorServiceClient(cc)

	//doUnary(c)
	//doServerStreaming(c)
	doClientStreaming(c)

}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.CalculatorRequest{
		Calculator: &calculatorpb.Calculator{
			NumberOne: 23,
			NumberTwo: 57,
		},
	}
	res, err := c.Calculator(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while Calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")

	req := &calculatorpb.CalculateManyTimesRequest{
		Calculator: &calculatorpb.Calculator{
			NumberOne: 533,
		},
	}

	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			//we've reached the end of the stream
			break
		}

		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	requests := []*calculatorpb.CalculateAverageRequest{
		&calculatorpb.CalculateAverageRequest{
			Calculator: &calculatorpb.Calculator{
				NumberOne: 10,
			},
		},
		&calculatorpb.CalculateAverageRequest{
			Calculator: &calculatorpb.Calculator{
				NumberOne: 50,
			},
		},
		&calculatorpb.CalculateAverageRequest{
			Calculator: &calculatorpb.Calculator{
				NumberOne: 30,
			},
		},
		&calculatorpb.CalculateAverageRequest{
			Calculator: &calculatorpb.Calculator{
				NumberOne: 20,
			},
		},
		&calculatorpb.CalculateAverageRequest{
			Calculator: &calculatorpb.Calculator{
				NumberOne: 70,
			},
		},
	}

	stream, err := c.ComputeAverage(context.Background())

	if err != nil {
		log.Fatalf("error while calling ComputeAverage: %v", err)
	}

	// we iterate over our slice and send each message individually
	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("error while receiving respose from ComputeAverage: %v", err)
	}

	fmt.Printf("LongGreet response: %v\n", res)

}
