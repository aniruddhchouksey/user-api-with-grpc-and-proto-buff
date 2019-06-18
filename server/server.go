package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	users "users/proto"

	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
)

const (
	psqlInfo = "host=localhost port=5432 user=postgres password=redhat dbname=postgres sslmode=disable"
)

type server struct {
	db *sql.DB
}

func main() {
	listener, err := net.Listen("tcp", "localhost:50055")
	if err != nil {
		log.Fatalf("cannot intialize the listener, %v", err)
	}
	srv := grpc.NewServer()
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println(psqlInfo)
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	users.RegisterUserProfilesServer(srv, &server{db})
	if err := srv.Serve(listener); err != nil {
		log.Fatalf("failed to listen %v", err)
	}

}

func (s *server) CreateUserProfile(ctx context.Context, req *users.CreateUserProfileRequest) (*users.UserProfile, error) {
	fmt.Printf("create user profile invoked %v", req)
	// Generating a unique id for the database
	id, err := uuid.NewV4()
	if err != nil {
		errors.Wrap(err, "error occured in create user profile method")
	}
	req.UserProfile.Id = id.String()
	sqlStatement := `INSERT INTO users (id, firstname, lastname, email) VALUES ($1, $2, $3, $4)`
	var result sql.Result
	if result, err = s.db.Exec(sqlStatement, req.UserProfile.Id, req.UserProfile.FirstName, req.UserProfile.LastName, req.UserProfile.Email); err != nil {
		errors.Wrap(err, "error occured while inserting into the database")
	}
	fmt.Println((result))
	return req.UserProfile, nil
}

func (s *server) GetUserProfile(ctx context.Context, req *users.GetUserProfileRequest) (*users.UserProfile, error) {
	sqlStatement := `SELECT * FROM users WHERE id = $1`
	result, err := s.db.Query(sqlStatement, req.Id)
	if err != nil {
		errors.Wrap(err, "error occured in get user profile method")
	}
	var userProfiles []users.UserProfile
	for result.Next() {
		var tmpuser users.UserProfile
		err := result.Scan(&tmpuser.Id, &tmpuser.FirstName, &tmpuser.LastName, &tmpuser.Email)
		if err != nil {
			errors.Wrap(err, "error occured in getuserprofile method")
		}
		userProfiles = append(userProfiles, tmpuser)
	}
	return &userProfiles[0], nil
}

func (s *server) DeleteUserProfile(ctx context.Context, req *users.DeleteUserProfileRequest) (*google_protobuf.Empty, error) {

	sqlStatement := `DELETE FROM users WHERE id = $1`

	_, err := s.db.Query(sqlStatement, req.Id)

	if err != nil {
		errors.Wrap(err, "error occured in deleteuserprofile method")
	}

	return &google_protobuf.Empty{}, nil
}

func (s *server) UpdateUserProfile(ctx context.Context, req *users.UpdateUserProfileRequest) (*users.UserProfile, error) {

	stmt, err := s.db.Prepare(`UPDATE users SET firstname=$1, lastname=$2, email=$3 WHERE id = $4`)
	if err != nil {
		errors.Wrap(err, "error occured in update user profile method")
	}
	_, err = stmt.Exec(req.UserProfile.FirstName, req.UserProfile.LastName, req.UserProfile.Email, req.UserProfile.Id)
	if err != nil {
		errors.Wrap(err, "error occured in updateuserprofile method")
	}

	return req.UserProfile, nil
}

func (s *server) ListUsersProfiles(ctx context.Context, req *users.ListUsersProfilesRequest) (*users.ListUsersProfilesResponse, error) {
	sqlStatement := `SELECT * FROM users WHERE firstname LIKE $1`

	result, err := s.db.Query(sqlStatement, "%"+req.Query+"%")
	if err != nil {
		errors.Wrap(err, "error occured in list user profile method")
	}
	// createing a temporary place to hold the outpu from the database
	var userProfiles users.ListUsersProfilesResponse
	var tmpuser users.UserProfile
	for result.Next() {
		err := result.Scan(&tmpuser.Id, &tmpuser.FirstName, &tmpuser.LastName, &tmpuser.Email)
		if err != nil {
			errors.Wrap(err, "error occured in list user profile method while scanning")
		}
		userProfiles.Profiles = append(userProfiles.Profiles, &tmpuser)
	}
	return &userProfiles, nil
}
