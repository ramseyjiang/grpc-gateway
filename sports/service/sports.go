package service

import (
	"github.com/ramseyjiang/grpc-gateway/sports/db"
	"github.com/ramseyjiang/grpc-gateway/sports/proto/sports"
	"golang.org/x/net/context"
)

type Sports interface {
	// ListSports will return a collection of sports.
	ListSports(ctx context.Context, in *sports.ListSportsRequest) (*sports.ListSportsResponse, error)
}

// sportsService implements the sports interface.
type sportsService struct {
	sportsRepo db.SportsRepo
}

// NewSportsService instantiates and returns a new sportsService.
func NewSportsService(sportsRepo db.SportsRepo) Sports {
	return &sportsService{sportsRepo}
}

func (s *sportsService) ListSports(ctx context.Context, in *sports.ListSportsRequest) (*sports.ListSportsResponse, error) {
	sportsResult, err := s.sportsRepo.List(in.Filter)
	if err != nil {
		return nil, err
	}

	return &sports.ListSportsResponse{Sports: sportsResult}, nil
}
