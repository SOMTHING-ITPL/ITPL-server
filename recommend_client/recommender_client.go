package recommend_client

import (
	context "context"
	"log"
	"strconv"
	"time"

	"github.com/SOMTHING-ITPL/ITPL-server/artist"
	"github.com/SOMTHING-ITPL/ITPL-server/performance"
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(socket string) (*RecommenderRPCClient, error) {
	conn, err := grpc.NewClient(socket, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
		return &RecommenderRPCClient{}, err
	}

	client := NewRecommenderClient(conn)

	return &RecommenderRPCClient{Conn: conn, Client: client}, nil
}

func (c *RecommenderRPCClient) PerformanceRecommendation(userRepo *user.Repository, performanceRepo *performance.Repository, artistRepo *artist.Repository, userID uint, topK int32) ([]performance.Performance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// user가 없을 경우 err
	_, err := userRepo.GetById(userID)
	if err != nil {
		log.Printf("failed to get user")
		return nil, err
	}

	// load genres (genre id)
	genres, err := userRepo.GetUserGenres(userID)
	if err != nil {
		log.Printf("failed to get user genres")
		return nil, err
	}
	var reqgenres []int32
	for _, genre := range genres {
		reqgenres = append(reqgenres, int32(genre.ID))
	}

	// load artists (artist name)
	artists, err := artistRepo.GetUserArtists(userID)
	if err != nil {
		log.Printf("failed to get user artists")
		return nil, err
	}
	var reqartists []string
	for _, artist := range artists {
		reqartists = append(reqartists, artist.Name)
	}

	// load user like (performance title)
	favs, err := performanceRepo.GetUserLike(userID)
	if err != nil {
		log.Printf("failed to get user likes")
		return nil, err
	}
	var reqfavs []string
	for _, fav := range favs {
		reqfavs = append(reqfavs, fav.Title)
	}

	userRequest := &UserRequest{
		UserId:  int32(userID),
		Genres:  reqgenres,
		Artists: reqartists,
		FavIds:  reqfavs,
		Topk:    topK,
	}

	/*request to server stub*/
	resp, err := c.Client.Recommend(ctx, userRequest)
	if err != nil {
		return nil, err
	}

	var result []performance.Performance
	for _, item := range resp.Concerts {
		if item == nil {
			log.Printf("WARNING: Received a nil Concert item from gRPC server. Skipping.")
			continue // nil 항목은 건너뛰고 다음 루프로 이동
		}
		ID, err := strconv.ParseInt(item.Id, 10, 64)
		if err != nil {
			log.Printf("failed to convert rpc response(concert.Id) to int")
			return nil, err
		}
		ap, err := performanceRepo.GetPerformanceById(uint(ID))
		if err != nil {
			log.Printf("failed to find performance by ID, after GRPC CALL - recommendation")
			return nil, err
		}
		result = append(result, *ap)

	}

	return result, nil
}
