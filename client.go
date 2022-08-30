package main

import (
	"context"
	"klausapp/softaware-test-task/service"
	"log"

	"google.golang.org/grpc"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	defer conn.Close()

	c := service.NewScoringServiceClient(conn)

	startDate := service.Date{
		Date: 1563314456,
	}
	endDate := service.Date{
		Date: 1563400856,
	}

	period := service.Period{
		StartDate: &startDate,
		EndDate:   &endDate,
	}

	responseCategory, err := c.GetCategoryScoreOverPeriodOfTime(context.Background(), &period)
	if err != nil {
		log.Fatalf("Error when calling GetCategoryScoreOverPeriodOfTime: %s", err)
	}

	log.Printf("Response from GetCategoryScoreOverPeriodOfTime: %s", responseCategory)

	responseTicket, err := c.GetScoresByTicketOverPeriodOfTime(context.Background(), &period)
	if err != nil {
		log.Fatalf("Error when calling GetScoresByTicketOverPeriodOfTime: %s", err)
	}

	log.Printf("Response from GetScoresByTicketOverPeriodOfTime: %s", responseTicket)

	responseScore, err := c.GetOveralQualityScoreOverPeriodOfTime(context.Background(), &period)
	if err != nil {
		log.Fatalf("Error when calling GetOveralQualityScoreOverPeriodOfTime: %s", err)
	}

	log.Printf("Response from GetOveralQualityScoreOverPeriodOfTime: %s", responseScore)

	responseScore, err = c.GetScoreChangeOverPeriodOfTime(context.Background(), &period)
	if err != nil {
		log.Fatalf("Error when calling GetScoreChangeOverPeriodOfTime: %s", err)
	}

	log.Printf("Response from GetScoreChangeOverPeriodOfTime: %s", responseScore)

}
