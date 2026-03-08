package model

import "time"

// tags json -> map field khi serialize/ deserialize JSON (API request/response)
// tags db -> map field khi serialize/ deserialize DB (DB query result -> struct)
// *Dataset: đại diện cho một tập dữ liệu địa lý, chứa thông tin về tên, số lượng điểm, thời gian tạo và xóa (nếu có).
// *DatasetPoint: đại diện cho một điểm dữ liệu trong một tập dữ liệu, chứa thông tin về tọa độ và khóa điểm.
// *MatchingJob: đại diện cho một công việc so khớp giữa hai tập dữ liệu, chứa thông tin về trạng thái, số lượng kết quả, thời gian thực hiện và lỗi (nếu có).
// *MatchingResult: đại diện cho một kết quả so khớp giữa hai điểm từ hai tập dữ liệu khác nhau, chứa thông tin về khoảng cách giữa hai điểm.

type Dataset struct {
	ID         string    `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	FileName   string    `json:"fileName" db:"file_name"`
	PointCount int       `json:"pointCount" db:"point_count"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	// DeletedAt có thể null nếu dataset chưa bị xóa, nên dùng con trỏ để phân biệt giữa giá trị null và giá trị thời gian thực tế.
	DeletedAt *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
}

type DatasetPoint struct {
	ID        int64     `json:"id" db:"id"`
	DatasetID string    `json:"datasetId" db:"dataset_id"`
	PointKey  string    `json:"pointKey" db:"point_key"`
	Latitude  float64   `json:"latitude" db:"latitude"`
	Longitude float64   `json:"longitude" db:"longitude"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type MatchingJobStatus string

const (
	MatchingJobStatusQueued    MatchingJobStatus = "queued"
	MatchingJobStatusRunning   MatchingJobStatus = "running"
	MatchingJobStatusCompleted MatchingJobStatus = "completed"
	MatchingJobStatusFailed    MatchingJobStatus = "failed"
)

type MatchingJob struct {
	ID           string            `json:"id" db:"id"`
	DatasetAID   string            `json:"datasetAId" db:"dataset_a_id"`
	DatasetBID   string            `json:"datasetBId" db:"dataset_b_id"`
	RadiusMeters int               `json:"radiusMeters" db:"radius_meters"`
	Status       MatchingJobStatus `json:"status" db:"status"`
	ResultCount  int               `json:"resultCount" db:"result_count"`
	// DurationMS có thể null nếu công việc chưa hoàn thành, nên dùng con trỏ để phân biệt giữa giá trị null và giá trị thực tế.
	DurationMS *int `json:"durationMs,omitempty" db:"duration_ms"`
	// ErrorMessage có thể null nếu công việc không gặp lỗi, nên dùng con trỏ để phân biệt giữa giá trị null và giá trị thực tế.
	ErrorMessage *string   `json:"errorMessage,omitempty" db:"error_message"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	// CompletedAt có thể null nếu công việc chưa hoàn thành, nên dùng con trỏ để phân biệt giữa giá trị null và giá trị thời gian thực tế.
	CompletedAt *time.Time `json:"completedAt,omitempty" db:"completed_at"`
}

type MatchingResult struct {
	ID             int64     `json:"id" db:"id"`
	JobID          string    `json:"jobId" db:"job_id"`
	PointAID       int64     `json:"pointAId" db:"point_a_id"`
	PointBID       int64     `json:"pointBId" db:"point_b_id"`
	DistanceMeters float64   `json:"distanceMeters" db:"distance_meters"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
}
