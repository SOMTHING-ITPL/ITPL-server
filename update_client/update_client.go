package update_client

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/SOMTHING-ITPL/ITPL-server/performance"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(socket string) (*UpdateClient, error) {
	conn, err := grpc.NewClient(socket, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
		return &UpdateClient{}, err
	}

	client := NewConcertUpdaterClient(conn)

	return &UpdateClient{Conn: conn, Client: client}, nil
}

func (c *UpdateClient) UpdateConcert(concert *[]performance.Performance) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //5초 timeout?
	defer cancel()

	resp, err := c.Client.UpdateConcerts(ctx, req)
	if err != nil {
		log.Fatalf("UpdateConcerts RPC failed: %v", err)
	}

	fmt.Printf("✅ 서버 응답: success=%v, message=%s\n", resp.GetSuccess(), resp.GetMessage())

}

func main() {
	// 서버 연결
	conn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := NewConcertUpdaterClient(conn)

	// 요청 데이터 준비
	req := &UpdateRequest{
		Concerts: []*Concert{
			{
				Id:        "123",
				Title:     "뮤지컬 라이온킹",
				Genre:     1,
				Cast:      []string{"홍길동", "김철수"},
				Keyword:   []string{"뮤지컬", "가족"},
				UpdatedAt: time.Now().Format(time.RFC3339),
			},
			{
				Id:        "456",
				Title:     "재즈 콘서트",
				Genre:     2,
				Cast:      []string{"존 콜트레인"},
				Keyword:   []string{"재즈", "연주"},
				UpdatedAt: time.Now().Format(time.RFC3339),
			},
		},
	}

	// RPC 호출
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.UpdateConcerts(ctx, req)
	if err != nil {
		log.Fatalf("UpdateConcerts RPC failed: %v", err)
	}

	fmt.Printf("✅ 서버 응답: success=%v, message=%s\n", resp.GetSuccess(), resp.GetMessage())
}
