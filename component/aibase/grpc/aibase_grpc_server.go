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
	"net"

	grpc "google.golang.org/grpc"
)

type AiBaseServer struct {
	UnimplementedAIBaseServiceServer
}

func (s *AiBaseServer) Call(ctx context.Context, req *CallRequest) (resp *CallResponse, err error) {
	return &CallResponse{}, nil
}

// 流式请求
func (s *AiBaseServer) Stream(s1 AIBaseService_StreamServer) error {
	for {
		r1, err2 := s1.Recv()
		if err2 != nil {
			return err2
		}
		log.Println(r1.Data)
	}
	return nil
}
func StartTestServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterAIBaseServiceServer(s, &AiBaseServer{})
	log.Println("Server started listening on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
