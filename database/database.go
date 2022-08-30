package database

import (
	"context"
	"database/sql"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CategoryRating struct {
	CategoryName string
	Date         string
	Rating       int32
	Count        int32
	Weight       float64
}

type TicketScore struct {
	TicketID     int64
	CategoryName string
	Score        int64
}

type Score struct {
	Score int64
}

func GetScoresByCategoriesOverPeriodOfTime(ctx context.Context, db *sql.DB, startDate int64, endDate int64) ([]CategoryRating, error) {

	preparedStmt, err := db.PrepareContext(ctx, getScoresByCategoriesOverPeriodOfTimeQuery())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed To Fetch")
	}

	rows, err := preparedStmt.QueryContext(ctx, startDate, endDate)

	if err != nil {
		return nil, status.Error(codes.Internal, "Could not get rows")
	}

	categoryRatings := []CategoryRating{}

	for rows.Next() {
		categoryRating := CategoryRating{}
		err = rows.Scan(
			&categoryRating.CategoryName,
			&categoryRating.Weight,
			&categoryRating.Rating,
			&categoryRating.Count,
			&categoryRating.Date,
		)
		if err != nil {
			return nil, status.Error(codes.Internal, "Issues creating readable entity")
		}
		categoryRatings = append(categoryRatings, categoryRating)
	}

	return categoryRatings, rows.Err()
}

func GetScoresByTicketOverPeriodOfTime(ctx context.Context, db *sql.DB, startDate int64, endDate int64) ([]TicketScore, error) {

	preparedStmt, err := db.PrepareContext(ctx, getScoresByTicketOverPeriodOfTimeQuery())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed To Fetch")
	}

	rows, err := preparedStmt.QueryContext(ctx, startDate, endDate)

	if err != nil {
		return nil, status.Error(codes.Internal, "Could not get rows")
	}

	ticketScores := []TicketScore{}

	for rows.Next() {
		ticketScore := TicketScore{}
		err = rows.Scan(
			&ticketScore.TicketID,
			&ticketScore.CategoryName,
			&ticketScore.Score,
		)
		if err != nil {
			return nil, status.Error(codes.Internal, "Issues creating readable entity")
		}
		ticketScores = append(ticketScores, ticketScore)
	}

	return ticketScores, rows.Err()
}

func GetOveralQualityScoreOverPeriodOfTime(ctx context.Context, db *sql.DB, startDate int64, endDate int64) (*Score, error) {
	preparedStmt, err := db.PrepareContext(ctx, getOveralQualityScoreQuery())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed To Fetch")
	}

	rows, err := preparedStmt.QueryContext(ctx, startDate, endDate)

	if err != nil {
		return nil, status.Error(codes.Internal, "Could not get rows")
	}

	overalScore := Score{}

	for rows.Next() {

		err = rows.Scan(
			&overalScore.Score,
		)
		if err != nil {
			log.Println(err)
			return nil, status.Error(codes.Internal, "Issues creating readable entity")
		}
	}

	return &overalScore, nil
}

func getOveralQualityScoreQuery() string {
	return `
		SELECT 
			COALESCE(ROUND(SUM(r.rating) / COUNT(r.rating) * rc.weight * 20), 0) as score
		FROM 
			ratings AS r
		INNER JOIN
			rating_categories as rc ON r.rating_category_id = rc.id
		WHERE
			r.created_at BETWEEN datetime(?, 'unixepoch') AND datetime(?, 'unixepoch')
		LIMIT 1
	`
}

func getScoresByTicketOverPeriodOfTimeQuery() string {
	return `
		SELECT 
			r.ticket_id,
			rc.name,
			COALESCE(ROUND(SUM(r.rating) / COUNT(r.rating) * rc.weight * 20), 0) as score
		FROM 
			ratings AS r
		INNER JOIN
			rating_categories as rc ON r.rating_category_id = rc.id
		WHERE
			r.created_at BETWEEN datetime(?, 'unixepoch') AND datetime(?, 'unixepoch')
		GROUP BY r.ticket_id,r.rating_category_id;
	`
}

func getScoresByCategoriesOverPeriodOfTimeQuery() string {

	return `
		SELECT 
			rc.name,
			rc.weight,
			SUM(r.rating) as rating_sum,
			COUNT(r.rating) as rating_count,
			r.created_at as created_at
		FROM 
			ratings AS r
		INNER JOIN
			rating_categories as rc ON r.rating_category_id = rc.id
		WHERE
			r.created_at BETWEEN datetime(?, 'unixepoch') AND datetime(?, 'unixepoch')
		GROUP BY r.rating_category_id, DATE(r.created_at)
		;
	`
}
