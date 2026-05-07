package dtos

type AdminReview struct {
	AdminReview string `json:"admin_review"`
}

type AdminReviewResponse struct {
	RankingName string `json:"ranking_name"`
	AdminReview string `json:"admin_review"`
}