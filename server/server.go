package main

import (
	"database/sql"
	"klausapp/softaware-test-task/service"
	"log"
	"net"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen to port 9000: %v", err)
		panic(err)
	}

	database, dberr := sql.Open("sqlite3", "./database.db")

	if dberr != nil {
		log.Println("Can't connect to database!")
		panic(dberr)
	}

	if database != nil {
		log.Println("Database has been loaded")
	}

	server := service.Server{
		Database: database,
	}

	grpcServer := grpc.NewServer()

	service.RegisterScoringServiceServer(grpcServer, &server)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9000: %v", err)
	}
}
