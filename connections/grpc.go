package connections

import (
	"fmt"
	"google.golang.org/grpc"
	"io"
)

func GRPCStream[T any](stream grpc.ClientStream, input string) Streamer[*T] {
	var generator Streamer[*T] = func() (*T, error) {
		m := new(T)
		err := stream.RecvMsg(m)
		if err == io.EOF {
			return nil, fmt.Errorf("stream for input %s ended successfully", input)
		} else if err != nil {
			return nil, err
		}
		return m, nil
	}

	return generator
}
