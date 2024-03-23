// Copyright (C) 2024 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package aibase

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	// 连接到 gRPC 服务器
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()
	c := NewAIBaseServiceClient(conn)

	// 调用普通请求 Call
	callReq := &CallRequest{Data: []byte("Hello")}
	callRes, err := c.Call(context.Background(), callReq)
	if err != nil {
		log.Fatalf("Call failed: %v", err)
	}
	log.Printf("Call response: %s", string(callRes.Result))

	// 调用流式请求 Stream
	// streamReq := &StreamRequest{Data: []byte("Hello")}
	stream, err := c.Stream(context.Background())
	if err != nil {
		log.Fatalf("Stream failed: %v", err)
	}
	for {
		// stream.Recv()
		err := stream.Send(&StreamRequest{})
		if err != nil {
			log.Fatalf("Stream receive failed: %v", err)
		}
	}
}
