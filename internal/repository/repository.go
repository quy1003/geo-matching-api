package repository

import (
	"context"

	"github.com/quy1003/geo-matching-api/internal/model"
)

type DatasetRepository interface {
	CreateDataset(ctx context.Context, dataset model.Dataset, points []model.DatasetPoint) (model.Dataset, []model.DatasetPoint, error)
	ListDatasets(ctx context.Context) ([]model.Dataset, error)
	GetDatasetByID(ctx context.Context, datasetID string) (model.Dataset, error)
	GetPointsByDatasetID(ctx context.Context, datasetID string) ([]model.DatasetPoint, error)
	DeleteDataset(ctx context.Context, datasetID string) error
}

type MatchingRepository interface {
	CreateJob(ctx context.Context, job model.MatchingJob) (model.MatchingJob, error)
	CompleteJob(ctx context.Context, jobID string, resultCount int, durationMS int) (model.MatchingJob, error)
	FailJob(ctx context.Context, jobID string, message string) (model.MatchingJob, error)
	SaveResults(ctx context.Context, jobID string, results []model.MatchingResult) ([]model.MatchingResult, error)
}
