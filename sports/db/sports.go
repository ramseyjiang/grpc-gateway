package db

import (
	"database/sql"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ramseyjiang/grpc-gateway/sports/proto/sports"
)

// SportsRepo provides repository access to sports.
type SportsRepo interface {
	// Init will initialise our sports repository.
	Init() error

	// List will return a list of sports.
	List(filter *sports.ListSportsRequestFilter) ([]*sports.Sport, error)
}

type sportsRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewSportsRepo creates a new sports repository.
func NewSportsRepo(db *sql.DB) SportsRepo {
	return &sportsRepo{db: db}
}

// Init prepares the sport repository dummy data.
func (s *sportsRepo) Init() error {
	var err error

	s.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy sports.
		err = s.seed()
	})

	return err
}

func (s *sportsRepo) List(filter *sports.ListSportsRequestFilter) ([]*sports.Sport, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getSportQueries()[sportsList]

	query, args = s.applyFilter(query, filter)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return s.scanSports(rows)
}

func (s *sportsRepo) applyFilter(query string, filter *sports.ListSportsRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if filter == nil {
		return query, args
	}

	if filter.Visible == true {
		clauses = append(clauses, "visible = 1")
	}

	if len(clauses) != 0 && len(filter.Column) == 0 && len(filter.OrderBy) == 0 {
		query += " WHERE " + strings.Join(clauses, " ")
	}

	if len(filter.Column) > 0 && len(filter.OrderBy) > 0 {
		clauses = append(clauses, "ORDER BY "+filter.Column+" "+filter.OrderBy)
		query += " WHERE " + strings.Join(clauses, " ")
	}

	// check sql correct or not
	// log.Println(filter, query)
	return query, args
}

func (s *sportsRepo) scanSports(rows *sql.Rows) ([]*sports.Sport, error) {
	var allSports []*sports.Sport

	for rows.Next() {
		var sport sports.Sport
		var advertisedStart time.Time
		var sportStart time.Time
		var sportEnd time.Time

		if err := rows.Scan(&sport.Id, &sport.Name, &sport.Visible, &sport.Status, &sportStart, &sportEnd, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		startTime, err := ptypes.TimestampProto(sportStart)
		if err != nil {
			return nil, err
		}
		sport.StartTime = startTime

		endTime, err := ptypes.TimestampProto(sportEnd)
		if err != nil {
			return nil, err
		}
		sport.EndTime = endTime

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}
		sport.AdvertisedStartTime = ts

		if time.Now().Unix() > ts.AsTime().Unix() {
			sport.Status = "CLOSED"
		} else {
			sport.Status = "OPEN"
		}

		allSports = append(allSports, &sport)
	}

	return allSports, nil
}
