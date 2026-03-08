package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/quy1003/geo-matching-api/internal/dto"
	"github.com/quy1003/geo-matching-api/internal/model"
	"github.com/quy1003/geo-matching-api/internal/repository"
	"github.com/quy1003/geo-matching-api/internal/repository/memory"
	"github.com/quy1003/geo-matching-api/internal/utils"
)

var ErrInvalidRadius = errors.New("radius must be greater than 0")

type MatchingService interface {
	Run(ctx context.Context, req dto.RunMatchingRequest) (dto.RunMatchingResponse, error)
}

type matchingService struct {
	datasetRepo  repository.DatasetRepository
	matchingRepo repository.MatchingRepository
}

func NewMatchingService(datasetRepo repository.DatasetRepository, matchingRepo repository.MatchingRepository) MatchingService {
	return &matchingService{
		datasetRepo:  datasetRepo,
		matchingRepo: matchingRepo,
	}
}

func (s *matchingService) Run(ctx context.Context, req dto.RunMatchingRequest) (dto.RunMatchingResponse, error) {
	if req.Radius <= 0 {
		return dto.RunMatchingResponse{}, ErrInvalidRadius
	}

	_, err := s.datasetRepo.GetDatasetByID(ctx, req.DatasetA)
	if err != nil {
		if errors.Is(err, memory.ErrNotFound) {
			return dto.RunMatchingResponse{}, fmt.Errorf("datasetA %w", ErrDatasetNotFound)
		}
		return dto.RunMatchingResponse{}, err
	}
	_, err = s.datasetRepo.GetDatasetByID(ctx, req.DatasetB)
	if err != nil {
		if errors.Is(err, memory.ErrNotFound) {
			return dto.RunMatchingResponse{}, fmt.Errorf("datasetB %w", ErrDatasetNotFound)
		}
		return dto.RunMatchingResponse{}, err
	}

	pointsA, err := s.datasetRepo.GetPointsByDatasetID(ctx, req.DatasetA)
	if err != nil {
		return dto.RunMatchingResponse{}, err
	}
	pointsB, err := s.datasetRepo.GetPointsByDatasetID(ctx, req.DatasetB)
	if err != nil {
		return dto.RunMatchingResponse{}, err
	}

	job, err := s.matchingRepo.CreateJob(ctx, model.MatchingJob{
		ID:           utils.NewID("job"),
		DatasetAID:   req.DatasetA,
		DatasetBID:   req.DatasetB,
		RadiusMeters: req.Radius,
		Status:       model.MatchingJobStatusRunning,
	})
	if err != nil {
		return dto.RunMatchingResponse{}, err
	}

	start := time.Now()
	matches := make([]dto.MatchPair, 0)
	results := make([]model.MatchingResult, 0)

	for _, pointA := range pointsA {
		for _, pointB := range pointsB {
			distance := utils.HaversineMeters(pointA.Latitude, pointA.Longitude, pointB.Latitude, pointB.Longitude)
			if distance <= float64(req.Radius) {
				matches = append(matches, dto.MatchPair{
					PointA:         pointA,
					PointB:         pointB,
					DistanceMeters: distance,
				})
				results = append(results, model.MatchingResult{
					JobID:          job.ID,
					PointAID:       pointA.ID,
					PointBID:       pointB.ID,
					DistanceMeters: distance,
				})
			}
		}
	}

	if _, err = s.matchingRepo.SaveResults(ctx, job.ID, results); err != nil {
		_, _ = s.matchingRepo.FailJob(ctx, job.ID, err.Error())
		return dto.RunMatchingResponse{}, err
	}

	durationMS := int(time.Since(start).Milliseconds())
	job, err = s.matchingRepo.CompleteJob(ctx, job.ID, len(matches), durationMS)
	if err != nil {
		return dto.RunMatchingResponse{}, err
	}

	return dto.RunMatchingResponse{
		Job:     job,
		Matches: matches,
	}, nil
}
