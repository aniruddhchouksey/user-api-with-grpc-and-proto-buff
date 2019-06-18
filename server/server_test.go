package main

import (
	"context"
	"log"
	"testing"
	users "users/proto"

	"google.golang.org/grpc"
)

func TestCreateUserProfile(t *testing.T) {
	conn, err := grpc.Dial("localhost:50055", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect, %v ", err)
	}

	// Making sure that the connection is cloased
	defer conn.Close()

	c := users.NewUserProfilesClient(conn)

	req := &users.CreateUserProfileRequest{
		UserProfile: &users.UserProfile{
			Id:        "1",
			FirstName: "aniruddhisawesome",
			LastName:  "chouksey",
			Email:     "aniruddh@appointy.com",
		},
	}

	res, err := c.CreateUserProfile(context.Background(), req)
	if err != nil {
		log.Fatalf("error while performing create %v", err)
	}
	log.Printf("response from the client : %v", res.Email)

}
