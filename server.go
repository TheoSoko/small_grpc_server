package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"small_grpc_server/proto"

	"google.golang.org/grpc"
)

var maxNumber = proto.Max{
	MessageID: 0,
	Max:       0,
}

type server struct {
	proto.UnimplementedMaxNumberServer
}

func (s *server) RegisterNumber(stream proto.MaxNumber_RegisterNumberServer) error {
	for {
		message, err := stream.Recv()
		// end of file = stream terminé
		if err == io.EOF {
			fmt.Print("End of File \n")
			return nil
		}
		if err != nil {
			fmt.Print("error after receving message \n", err)
			return err
		}

		if message.Num > maxNumber.Max {
			maxNumber.Max = message.Num
			maxNumber.MessageID = message.ID
		}

		if err := stream.Send(&maxNumber); err != nil {
			fmt.Print("error after sending message \n", err)
			return err
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:443")
	if err != nil {
		log.Fatal("failed to listen", err)
	}

	s := grpc.NewServer()
	proto.RegisterMaxNumberServer(s, &server{})

	log.Print("server listening at", listener.Addr())

	if err := s.Serve(listener); err != nil {
		log.Fatal("failed to serve", err)
	}
}
