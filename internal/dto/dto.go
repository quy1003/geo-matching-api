package dto

import "github.com/quy1003/geo-matching-api/internal/model"

type DatasetDetailResponse struct {
	Dataset model.Dataset        `json:"dataset"`
	Points  []model.DatasetPoint `json:"points"`
}

type RunMatchingRequest struct {
	DatasetA string `json:"datasetA" binding:"required"`
	DatasetB string `json:"datasetB" binding:"required"`
	Radius   int    `json:"radius" binding:"required"`
}

type MatchPair struct {
	PointA         model.DatasetPoint `json:"pointA"`
	PointB         model.DatasetPoint `json:"pointB"`
	DistanceMeters float64            `json:"distanceMeters"`
}

type RunMatchingResponse struct {
	Job     model.MatchingJob `json:"job"`
	Matches []MatchPair       `json:"matches"`
}
