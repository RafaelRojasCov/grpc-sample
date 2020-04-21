package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	greetpb "../greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello Client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	defer cc.Close()
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	c := greetpb.NewGreetServiceClient(cc)

	// doUnary(c)

	// doServerStreaming(c)

	// doClientStreaming(c)

	doBidiStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Rafael",
			LastName:  "Rojas",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while Calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Rafael",
			LastName:  "Rojas",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)

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

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Rafael",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Cookie",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tania",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Rayitas",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())

	if err != nil {
		log.Fatalf("error while calling LongGreet: %v", err)
	}

	// we iterate over our slice and send each message individually
	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("error while receiving respose from LongGreet: %v", err)
	}

	fmt.Printf("LongGreet response: %v\n", res)

}

func doBidiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Bidireccional Streaming RPC...")

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Rafael",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Cookie",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tania",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Rayitas",
			},
		},
	}

	// we create a stream by invoking the client

	stream, err := c.GreetEveryone(context.Background())

	if err != nil {
		log.Fatalf("Error while creating a stream: %v", err)
		return
	}

	waitc := make(chan struct{})

	go func() { // go Routine
		// function to send a bunch of messages
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// we send a bunch of messages to the server (go Routines)

	go func() { // go Routine
		// function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving a stream: %v", err)
				break
			}

			// we receive a bunch of messages from the server (go Routines)
			fmt.Printf("Received: %v", res.GetResult())
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc
}
