# Backend Source Guide

Tài liệu này mô tả cấu trúc source code backend hiện tại của Geo Matching System.

## 1. Tổng quan kiến trúc

Backend dùng kiến trúc phân lớp, tách rõ:

- `Handler` (HTTP layer)
- `Service` (business logic)
- `Repository` (data access abstraction)
- `Repository Memory` (lưu tạm in-memory cho MVP)

Mục tiêu là dễ thay thế lớp lưu trữ từ memory sang DB thật (PostgreSQL/MySQL) mà không phải sửa handler/service.

## 2. Cấu trúc thư mục chính

```text
geo-matching-api
├─ cmd/server
│  └─ main.go
├─ internal
│  ├─ app
│  │  └─ server.go
│  ├─ dto
│  │  └─ dto.go
│  ├─ handler
│  │  ├─ handler.go
│  │  ├─ dataset_handler.go
│  │  └─ matching_handler.go
│  ├─ model
│  │  └─ entities.go
│  ├─ repository
│  │  ├─ repository.go
│  │  └─ memory
│  │     └─ store.go
│  ├─ service
│  │  ├─ dataset_service.go
│  │  └─ matching_service.go
│  └─ utils
│     ├─ haversine.go
│     └─ id.go
└─ src
   ├─ README.md
   └─ dataset
      ├─ dataset-a.csv
      ├─ dataset-b.csv
      ├─ dataset-a.json
      └─ dataset-b.json
```

## 3. Chức năng từng phần

### `cmd/server/main.go`

- Entry point của ứng dụng.
- Khởi tạo server qua `app.NewServer()`.
- Chạy HTTP server tại `:8080`.

### `internal/app/server.go`

- Khởi tạo Gin engine + CORS.
- Wiring dependency injection:
  - tạo `memory.NewStore()`
  - inject vào `DatasetService`, `MatchingService`
  - inject services vào handler
- Đăng ký toàn bộ routes.

### `internal/model/entities.go`

Khai báo entity mapping theo schema dữ liệu:

- `Dataset`
- `DatasetPoint`
- `MatchingJob`
- `MatchingResult`
- `MatchingJobStatus` (`queued`, `running`, `completed`, `failed`)

### `internal/repository/repository.go`

Khai báo interface repository:

- `DatasetRepository`
- `MatchingRepository`

Đây là contract để service làm việc độc lập với loại database.

### `internal/repository/memory/store.go`

Implementation repository bằng in-memory:

- dùng `map + sync.RWMutex`
- lưu dữ liệu dataset, points, jobs, results
- hỗ trợ CRUD dataset và lifecycle matching job

### `internal/service/dataset_service.go`

Business logic dataset:

- parse/validate CSV upload
- kiểm tra header `id,latitude,longitude`
- validate tọa độ hợp lệ
- tạo dataset và lưu points
- list/get/delete dataset

### `internal/service/matching_service.go`

Business logic matching:

- validate input (`radius > 0`)
- lấy điểm từ dataset A/B
- chạy thuật toán nested loop + Haversine
- lưu `matching_results`
- cập nhật trạng thái `matching_jobs`

### `internal/handler/*.go`

HTTP handlers:

- `GET /health`
- `POST /datasets/upload` (CSV)
- `GET /datasets`
- `GET /datasets/:id`
- `DELETE /datasets/:id`
- `POST /geo-matching/run`

### `internal/dto/dto.go`

Định nghĩa request/response object cho API:

- `RunMatchingRequest`
- `RunMatchingResponse`
- `DatasetDetailResponse`

### `internal/utils/*.go`

Utility dùng chung:

- `haversine.go`: tính khoảng cách mét giữa 2 tọa độ
- `id.go`: phát sinh ID dạng prefix cho dữ liệu runtime

## 4. Luồng xử lý chính

### Upload dataset

1. Client gọi `POST /datasets/upload`
2. Handler nhận file CSV
3. Service parse + validate dữ liệu
4. Service gọi repository để lưu dataset + points
5. Trả metadata dataset cho client

### Run matching

1. Client gọi `POST /geo-matching/run`
2. Service tạo `matching_job` trạng thái `running`
3. Service duyệt điểm A x B, tính Haversine
4. Lọc các cặp có `distance <= radius`
5. Lưu `matching_results`
6. Cập nhật `matching_job` thành `completed`
7. Trả `job + matches`

## 5. Dữ liệu mẫu

Thư mục `src/dataset` chứa dữ liệu mẫu để test upload/matching:

- `dataset-a.csv`, `dataset-b.csv`
- `dataset-a.json`, `dataset-b.json`

Lưu ý: API hiện chỉ nhận upload CSV. JSON ở đây dùng làm dữ liệu tham khảo.

## 6. Hướng mở rộng sang DB thật

Khi tích hợp DB:

1. Tạo package repository mới (ví dụ `internal/repository/postgres`)
2. Implement đầy đủ interface trong `repository.go`
3. Đổi wiring tại `internal/app/server.go` từ `memory.NewStore()` sang `postgres.NewStore(...)`
4. Giữ nguyên handler/service

