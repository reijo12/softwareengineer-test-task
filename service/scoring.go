package service

import (
	"context"
	"database/sql"
	"klausapp/softaware-test-task/database"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type Server struct {
	UnimplementedScoringServiceServer
	Database *sql.DB
}

func calculateScoreWeightedScore(scoreSum int64, scoreCount int64, weight float64) int64 {

	averageScore := float64(scoreSum / scoreCount)
	weightedAverageScore := averageScore * weight
	return int64(weightedAverageScore * 20)
}

func (s *Server) GetCategoryScoreOverPeriodOfTime(ctx context.Context, period *Period) (*ScoresByCategories, error) {
	valid, err := isValidPeriod(period)
	if !valid {
		return nil, err
	}

	response, err := database.GetScoresByCategoriesOverPeriodOfTime(ctx, s.Database, period.StartDate.Date, period.EndDate.Date)
	if err != nil {
		return nil, err
	}

	scoreByCategoriesList := &ScoresByCategories{}

	scoreByCategoryMap := make(map[string]*ScoreByCategory)
	totalRatingCountMap := make(map[string]int64)
	totalAverageWeightedScoreByCategoryMap := make(map[string]int64)

	for _, categoryRating := range response {

		calculatedAverageWeightedScore := calculateScoreWeightedScore(int64(categoryRating.Rating), int64(categoryRating.Count), categoryRating.Weight)

		totalRatingCountMap[categoryRating.CategoryName] += int64(categoryRating.Count)
		totalAverageWeightedScoreByCategoryMap[categoryRating.CategoryName] += calculatedAverageWeightedScore

		scoreMapRefrence, ok := scoreByCategoryMap[categoryRating.CategoryName]
		if !ok {
			scoreMapRefrence = &ScoreByCategory{
				CategoryScore: &CategoryScore{
					CategoryName: categoryRating.CategoryName,
					Score: calculateScoreWeightedScore(
						totalAverageWeightedScoreByCategoryMap[categoryRating.CategoryName],
						totalRatingCountMap[categoryRating.CategoryName],
						1.0,
					),
				},
				Ratings: totalRatingCountMap[categoryRating.CategoryName],
				DateScores: []*DateScore{
					{
						Date:  categoryRating.Date,
						Score: calculatedAverageWeightedScore,
					},
				},
			}
			scoreByCategoryMap[categoryRating.CategoryName] = scoreMapRefrence
		} else {
			scoreMapRefrence.DateScores = append(scoreMapRefrence.DateScores, &DateScore{
				Date:  categoryRating.Date,
				Score: calculatedAverageWeightedScore,
			})
			scoreMapRefrence.CategoryScore.Score = calculateScoreWeightedScore(
				totalAverageWeightedScoreByCategoryMap[categoryRating.CategoryName],
				totalRatingCountMap[categoryRating.CategoryName],
				1.0,
			)
		}
	}

	for _, item := range scoreByCategoryMap {
		scoreByCategoriesList.ScoresByCategories = append(scoreByCategoriesList.ScoresByCategories, item)
	}

	return scoreByCategoriesList, nil
}

func (s *Server) GetScoresByTicketOverPeriodOfTime(ctx context.Context, period *Period) (*ScoresByTickets, error) {
	valid, err := isValidPeriod(period)
	if !valid {
		return nil, err
	}

	response, err := database.GetScoresByTicketOverPeriodOfTime(ctx, s.Database, period.StartDate.Date, period.EndDate.Date)
	if err != nil {
		return nil, err
	}

	scoresByTickets := &ScoresByTickets{}
	scoreByTicketsMap := make(map[int64]*ScoreByTicket)

	for _, ticketRating := range response {

		ticketEntry, ok := scoreByTicketsMap[ticketRating.TicketID]
		if !ok {
			ticketEntry = &ScoreByTicket{
				TicketId: ticketRating.TicketID,
				CategoryScores: []*CategoryScore{
					{
						CategoryName: ticketRating.CategoryName,
						Score:        ticketRating.Score,
					},
				},
			}
			scoreByTicketsMap[ticketRating.TicketID] = ticketEntry
		} else {
			ticketEntry.CategoryScores = append(ticketEntry.CategoryScores, &CategoryScore{
				CategoryName: ticketRating.CategoryName,
				Score:        ticketRating.Score,
			})
		}
	}

	for _, item := range scoreByTicketsMap {
		scoresByTickets.ScoresByTickets = append(scoresByTickets.ScoresByTickets, item)
	}

	return scoresByTickets, nil
}

func (s *Server) GetOveralQualityScoreOverPeriodOfTime(ctx context.Context, period *Period) (*Score, error) {
	valid, err := isValidPeriod(period)
	if !valid {
		return nil, err
	}

	response, err := database.GetOveralQualityScoreOverPeriodOfTime(ctx, s.Database, period.StartDate.Date, period.EndDate.Date)
	if err != nil {
		return nil, err
	}

	score := Score{
		Score: response.Score,
	}

	return &score, nil
}

func (s *Server) GetScoreChangeOverPeriodOfTime(ctx context.Context, period *Period) (*Score, error) {
	valid, err := isValidPeriod(period)
	if !valid {
		return nil, err
	}

	periodLength := period.EndDate.Date - period.StartDate.Date

	initialPeriodScore, err := database.GetOveralQualityScoreOverPeriodOfTime(ctx, s.Database, period.StartDate.Date, period.EndDate.Date)

	if err != nil {
		return nil, err
	}

	newStartDate := period.StartDate.Date - periodLength
	previousPeriodScore, err := database.GetOveralQualityScoreOverPeriodOfTime(ctx, s.Database, newStartDate, period.StartDate.Date)

	if err != nil {
		return nil, err
	}

	scoreChange := initialPeriodScore.Score - previousPeriodScore.Score

	score := Score{
		Score: scoreChange,
	}

	return &score, nil
}

func isValidPeriod(period *Period) (bool, error) {
	if period.StartDate == nil {
		return false, status.Error(codes.InvalidArgument, "Missing startdate")
	}

	if period.EndDate == nil {
		return false, status.Error(codes.InvalidArgument, "Missing endDate")
	}

	if period.EndDate.Date < period.StartDate.Date {
		return false, status.Error(codes.InvalidArgument, "End date cant be bigger than start date")
	}

	return true, nil
}
