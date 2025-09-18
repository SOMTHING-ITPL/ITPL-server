package update_client

import (
	"context"
	"fmt"
	"log"
	"strconv"
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

func (c *UpdateClient) UpdateConcert(concerts []performance.Performance) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //5초 timeout?
	defer cancel()

	req := []*Concert{}
	for _, ccn := range concerts {
		parsedCast, err := performance.ParsingCast(*ccn.Cast)
		if err != nil {
			return fmt.Errorf("fail to parse cast: %w", err)
		}
		parsedKeyword, err := performance.ParsingKeyword(ccn.Keyword)
		if err != nil {
			return fmt.Errorf("fail to parse keyword: %w", err)
		}

		req = append(req, &Concert{
			Id:        strconv.FormatUint(uint64(ccn.ID), 10),
			Title:     ccn.Title,
			Genre:     int32(ccn.Genre),
			Cast:      parsedCast,
			Keyword:   parsedKeyword,
			UpdatedAt: ccn.UpdatedAt.Format(time.RFC3339),
		})
	}

	resp, err := c.Client.UpdateConcerts(ctx, &UpdateRequest{
		Concerts: req,
	})

	if err != nil {
		log.Fatalf("UpdateConcerts RPC failed: %v", err)
		return err
	}

	fmt.Printf("서버 응답: success=%v, message=%s\n", resp.GetSuccess(), resp.GetMessage())
	return nil
}
