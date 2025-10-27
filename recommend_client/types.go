package recommend_client

import grpc "google.golang.org/grpc"

type RecommenderRPCClient struct {
	Conn   *grpc.ClientConn
	Client RecommenderClient
}
