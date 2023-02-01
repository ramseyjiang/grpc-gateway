package db

import (
	"database/sql"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ramseyjiang/grpc-gateway/racing/proto/racing"
)

// RacesRepo provides repository access to races.
type RacesRepo interface {
	// Init will initialise our races repository.
	Init() error

	// List will return a list of races.
	List(filter *racing.ListRacesRequestFilter) ([]*racing.Race, error)

	// GetRace will return a race collection by id
	GetRace(req *racing.GetRaceRequest) (*racing.Race, error)
}

type racesRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewRacesRepo creates a new races repository.
func NewRacesRepo(db *sql.DB) RacesRepo {
	return &racesRepo{db: db}
}

// Init prepares the race repository dummy data.
func (r *racesRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy races.
		err = r.seed()
	})

	return err
}

func (r *racesRepo) GetRace(req *racing.GetRaceRequest) (*racing.Race, error) {
	var (
		err   error
		query string
	)

	query = getRaceQueries()[racesList]

	if req.Id > 0 {
		id := strconv.Itoa(int(req.Id))
		query += " WHERE id=" + id
	}

	row, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	return r.scanRace(row)
}

func (r *racesRepo) List(filter *racing.ListRacesRequestFilter) ([]*racing.Race, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getRaceQueries()[racesList]

	query, args = r.applyFilter(query, filter)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanRaces(rows)
}

func (r *racesRepo) scanRace(row *sql.Rows) (*racing.Race, error) {
	var advertisedStart time.Time
	race := &racing.Race{}

	for row.Next() {
		if err := row.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		race.AdvertisedStartTime = ts
		if time.Now().Unix() > ts.AsTime().Unix() {
			race.Status = "CLOSED"
		} else {
			race.Status = "OPEN"
		}
	}

	return race, nil
}

func (r *racesRepo) applyFilter(query string, filter *racing.ListRacesRequestFilter) (string, []interface{}) {
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

	if len(filter.MeetingIds) > 0 {
		clauses = append(clauses, "meeting_id IN ("+strings.Repeat("?,", len(filter.MeetingIds)-1)+"?)")

		for _, meetingID := range filter.MeetingIds {
			args = append(args, meetingID)
		}
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

func (r *racesRepo) scanRaces(
	rows *sql.Rows,
) ([]*racing.Race, error) {
	var races []*racing.Race

	for rows.Next() {
		var race racing.Race
		var advertisedStart time.Time

		if err := rows.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		race.AdvertisedStartTime = ts
		if time.Now().Unix() > ts.AsTime().Unix() {
			race.Status = "CLOSED"
		} else {
			race.Status = "OPEN"
		}

		races = append(races, &race)
	}

	return races, nil
}
