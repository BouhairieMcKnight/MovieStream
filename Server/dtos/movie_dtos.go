package dtos

import "github.com/BouhairieMcKnight/MovieStream/Server/domain"

type AdminReview struct {
	AdminReview string `json:"admin_review"`
}

type AdminReviewResponse struct {
	RankingName string `json:"ranking_name"`
	AdminReview string `json:"admin_review"`
}

type MovieSearchDto struct {
	Movies       []domain.Movie `json:"results"`
	TotalPages   int            `json:"total_pages"`
	TotalResults int            `json:"total_results"`
	Term         string         `json:"term"`
	GenreId      int            `json:"genre_id"`
	GenreName    string         `json:"genre_name"`
}