package service

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/quy1003/geo-matching-api/internal/model"
	"github.com/quy1003/geo-matching-api/internal/repository"
	"github.com/quy1003/geo-matching-api/internal/repository/memory"
	"github.com/quy1003/geo-matching-api/internal/utils"
)

var (
	ErrInvalidCSVFormat  = errors.New("invalid csv format")
	ErrInvalidCoordinate = errors.New("invalid coordinate")
	ErrDatasetNotFound   = errors.New("dataset not found")
)

type DatasetService interface {
	UploadCSV(ctx context.Context, datasetName, fileName string, reader io.Reader) (model.Dataset, error)
	List(ctx context.Context) ([]model.Dataset, error)
	Get(ctx context.Context, datasetID string) (model.Dataset, []model.DatasetPoint, error)
	Delete(ctx context.Context, datasetID string) error
}

type datasetService struct {
	repo repository.DatasetRepository
}

func NewDatasetService(repo repository.DatasetRepository) DatasetService {
	return &datasetService{repo: repo}
}

func (s *datasetService) UploadCSV(ctx context.Context, datasetName, fileName string, reader io.Reader) (model.Dataset, error) {
	points, err := parseCSVPoints(reader)
	if err != nil {
		return model.Dataset{}, err
	}

	if datasetName == "" {
		datasetName = strings.TrimSuffix(fileName, ".csv")
	}

	datasetID := utils.NewID("ds")
	dataset := model.Dataset{
		ID:         datasetID,
		Name:       datasetName,
		FileName:   fileName,
		PointCount: len(points),
	}

	for i := range points {
		points[i].DatasetID = datasetID
	}

	storedDataset, _, err := s.repo.CreateDataset(ctx, dataset, points)
	return storedDataset, err
}

func (s *datasetService) List(ctx context.Context) ([]model.Dataset, error) {
	datasets, err := s.repo.ListDatasets(ctx)
	if err != nil {
		return nil, err
	}

	sort.Slice(datasets, func(i, j int) bool {
		return datasets[i].CreatedAt.After(datasets[j].CreatedAt)
	})

	return datasets, nil
}

func (s *datasetService) Get(ctx context.Context, datasetID string) (model.Dataset, []model.DatasetPoint, error) {
	dataset, err := s.repo.GetDatasetByID(ctx, datasetID)
	if err != nil {
		if errors.Is(err, memory.ErrNotFound) {
			return model.Dataset{}, nil, ErrDatasetNotFound
		}
		return model.Dataset{}, nil, err
	}

	points, err := s.repo.GetPointsByDatasetID(ctx, datasetID)
	if err != nil {
		if errors.Is(err, memory.ErrNotFound) {
			return model.Dataset{}, nil, ErrDatasetNotFound
		}
		return model.Dataset{}, nil, err
	}

	return dataset, points, nil
}

func (s *datasetService) Delete(ctx context.Context, datasetID string) error {
	err := s.repo.DeleteDataset(ctx, datasetID)
	if err != nil && errors.Is(err, memory.ErrNotFound) {
		return ErrDatasetNotFound
	}
	return err
}

func parseCSVPoints(reader io.Reader) ([]model.DatasetPoint, error) {
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidCSVFormat, err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("%w: missing data rows", ErrInvalidCSVFormat)
	}

	header := normalizeHeader(records[0])
	if len(header) < 3 || header[0] != "id" || header[1] != "latitude" || header[2] != "longitude" {
		return nil, fmt.Errorf("%w: expected headers id,latitude,longitude", ErrInvalidCSVFormat)
	}

	points := make([]model.DatasetPoint, 0, len(records)-1)
	for i, row := range records[1:] {
		if len(row) < 3 {
			return nil, fmt.Errorf("%w: row %d has fewer than 3 columns", ErrInvalidCSVFormat, i+2)
		}

		lat, err := strconv.ParseFloat(strings.TrimSpace(row[1]), 64)
		if err != nil {
			return nil, fmt.Errorf("%w: row %d latitude parse failed", ErrInvalidCoordinate, i+2)
		}
		lng, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
		if err != nil {
			return nil, fmt.Errorf("%w: row %d longitude parse failed", ErrInvalidCoordinate, i+2)
		}

		if math.Abs(lat) > 90 || math.Abs(lng) > 180 {
			return nil, fmt.Errorf("%w: row %d out of range", ErrInvalidCoordinate, i+2)
		}

		pointID := strings.TrimSpace(row[0])
		if pointID == "" {
			pointID = fmt.Sprintf("row-%d", i+1)
		}

		points = append(points, model.DatasetPoint{
			PointKey:  pointID,
			Latitude:  lat,
			Longitude: lng,
		})
	}

	return points, nil
}

func normalizeHeader(header []string) []string {
	out := make([]string, 0, len(header))
	for _, col := range header {
		out = append(out, strings.ToLower(strings.TrimSpace(col)))
	}
	return out
}
