package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/quy1003/geo-matching-api/internal/model"
)

var ErrNotFound = errors.New("resource not found")

type Store struct {
	mu            sync.RWMutex
	datasets      map[string]model.Dataset
	datasetPoints map[string][]model.DatasetPoint
	jobs          map[string]model.MatchingJob
	results       map[string][]model.MatchingResult
	nextPointID   int64
	nextResultID  int64
}

func NewStore() *Store {
	return &Store{
		datasets:      make(map[string]model.Dataset),
		datasetPoints: make(map[string][]model.DatasetPoint),
		jobs:          make(map[string]model.MatchingJob),
		results:       make(map[string][]model.MatchingResult),
	}
}

func (s *Store) CreateDataset(_ context.Context, dataset model.Dataset, points []model.DatasetPoint) (model.Dataset, []model.DatasetPoint, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	dataset.CreatedAt = time.Now().UTC()
	s.datasets[dataset.ID] = dataset

	storedPoints := make([]model.DatasetPoint, 0, len(points))
	for _, p := range points {
		s.nextPointID++
		p.ID = s.nextPointID
		p.CreatedAt = time.Now().UTC()
		storedPoints = append(storedPoints, p)
	}
	s.datasetPoints[dataset.ID] = storedPoints

	return dataset, append([]model.DatasetPoint(nil), storedPoints...), nil
}

func (s *Store) ListDatasets(_ context.Context) ([]model.Dataset, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]model.Dataset, 0, len(s.datasets))
	for _, dataset := range s.datasets {
		out = append(out, dataset)
	}
	return out, nil
}

func (s *Store) GetDatasetByID(_ context.Context, datasetID string) (model.Dataset, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	dataset, ok := s.datasets[datasetID]
	if !ok {
		return model.Dataset{}, ErrNotFound
	}
	return dataset, nil
}

func (s *Store) GetPointsByDatasetID(_ context.Context, datasetID string) ([]model.DatasetPoint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	points, ok := s.datasetPoints[datasetID]
	if !ok {
		if _, datasetExists := s.datasets[datasetID]; !datasetExists {
			return nil, ErrNotFound
		}
		return []model.DatasetPoint{}, nil
	}
	return append([]model.DatasetPoint(nil), points...), nil
}

func (s *Store) DeleteDataset(_ context.Context, datasetID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.datasets[datasetID]; !ok {
		return ErrNotFound
	}

	delete(s.datasets, datasetID)
	delete(s.datasetPoints, datasetID)
	return nil
}

func (s *Store) CreateJob(_ context.Context, job model.MatchingJob) (model.MatchingJob, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	job.CreatedAt = time.Now().UTC()
	s.jobs[job.ID] = job
	return job, nil
}

func (s *Store) CompleteJob(_ context.Context, jobID string, resultCount int, durationMS int) (model.MatchingJob, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, ok := s.jobs[jobID]
	if !ok {
		return model.MatchingJob{}, ErrNotFound
	}

	job.Status = model.MatchingJobStatusCompleted
	job.ResultCount = resultCount
	job.DurationMS = &durationMS
	completedAt := time.Now().UTC()
	job.CompletedAt = &completedAt
	s.jobs[jobID] = job

	return job, nil
}

func (s *Store) FailJob(_ context.Context, jobID string, message string) (model.MatchingJob, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, ok := s.jobs[jobID]
	if !ok {
		return model.MatchingJob{}, ErrNotFound
	}

	job.Status = model.MatchingJobStatusFailed
	job.ErrorMessage = &message
	completedAt := time.Now().UTC()
	job.CompletedAt = &completedAt
	s.jobs[jobID] = job

	return job, nil
}

func (s *Store) SaveResults(_ context.Context, jobID string, results []model.MatchingResult) ([]model.MatchingResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.jobs[jobID]; !ok {
		return nil, ErrNotFound
	}

	storedResults := make([]model.MatchingResult, 0, len(results))
	for _, result := range results {
		s.nextResultID++
		result.ID = s.nextResultID
		result.CreatedAt = time.Now().UTC()
		storedResults = append(storedResults, result)
	}
	s.results[jobID] = storedResults

	return append([]model.MatchingResult(nil), storedResults...), nil
}
