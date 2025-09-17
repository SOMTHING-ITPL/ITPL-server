package update_client

import "google.golang.org/grpc"

type UpdateClient struct {
	Conn   *grpc.ClientConn
	Client ConcertUpdaterClient
}
