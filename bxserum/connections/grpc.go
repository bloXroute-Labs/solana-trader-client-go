package connections

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"io"
)

func GRPCStream[T any](stream grpc.ClientStream, input string, responseChan chan *T) error {
	response, err := recvGRPCResponse[T](stream, input)
	if err != nil {
		return err
	}

	go func(response *T, stream grpc.ClientStream, input string) {
		responseChan <- response

		for {
			response, err = recvGRPCResponse[T](stream, input)
			if err != nil {
				logrus.Errorf(err.Error())
				return
			} else {
				responseChan <- response
			}
		}
	}(response, stream, input)

	return nil
}

func recvGRPCResponse[T any](stream grpc.ClientStream, input string) (*T, error) {
	m := new(T)
	err := stream.RecvMsg(m)
	if err == io.EOF {
		return nil, fmt.Errorf("stream for input %s ended successfully", input)
	} else if err != nil {
		return nil, err
	}

	return m, nil
}
