package main

import (
	"context"
	"fmt"
	"log"
	users "users/proto"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Client initialized...")

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

	// // Testing get user profile endpoint test
	// req := &users.GetUserProfileRequest{
	// 	Id: "1",
	// }
	// res, err := c.GetUserProfile(context.Background(), req)
	// if err != nil {
	// 	log.Fatalf("error while performing create %v", err)
	// }
	// fmt.Println(res)

	// // Testing get Delete request
	// req := &users.DeleteUserProfileRequest{
	// 	Id: "1",
	// }
	// res, err := c.DeleteUserProfile(context.Background(), req)
	// if err != nil {
	// 	log.Fatalf("error while performing create %v", err)
	// }
	// fmt.Println(res)

	// // Testing get update request
	// req := &users.UpdateUserProfileRequest{
	// 	UserProfile: &users.UserProfile{
	// 		Id:        "1",
	// 		FirstName: "manish",
	// 		LastName:  "chouksey",
	// 		Email:     "manish@appointy.com",
	// 	},
	// }
	// res, err := c.UpdateUserProfile(context.Background(), req)
	// if err != nil {
	// 	log.Fatalf("error while performing update %v", err)
	// }
	// fmt.Println(res)

	// // Testing the listing endpoint
	// req := &users.ListUsersProfilesRequest{
	// 	Query: "ani",
	// }
	// res, err := c.ListUsersProfiles(context.Background(), req)
	// if err != nil {
	// 	log.Fatalf("error in the list method %v", err)
	// }
	// fmt.Println(res)
}
