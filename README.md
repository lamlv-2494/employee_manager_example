# Employee Manager Example

A RESTful API server for managing employees and departments, built with Go and MySQL (no third-party frameworks).

## Tech Stack

- **Language:** Go 1.26
- **Database:** MySQL
- **Driver:** `github.com/go-sql-driver/mysql`
- **HTTP:** `net/http` (standard library)

## Project Structure

```
cmd/app/          # Entry point (main.go)
internal/
  handlers/       # HTTP handlers (employee, department)
  middleware/     # Logging & Recover middleware
  models/         # Data models
  repositories/   # Database access layer
  services/       # Business logic layer
  utils/          # Exporter (CSV, JSON),utils
  texts/          # Error/message constants
exports/          # Exported data files (CSV, JSON)
schema.sql        # Database schema
```

## Setup

### 1. Clone & cài dependencies

```bash
git clone <repo-url>
cd employee_manager_example
go mod tidy
```

### 2. Tạo database

```bash
mysql -u root -p < schema.sql
```

### 3. Tạo file `.env`

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=employee_management
```

### 4. Chạy server

```bash
go run cmd/app/main.go
```

Server khởi động tại `http://localhost:8080`

---

## API Endpoints

### Employees

| Method | Endpoint             | Mô tả                        |
|--------|----------------------|------------------------------|
| GET    | /employees           | Lấy danh sách nhân viên      |
| POST   | /employees           | Tạo nhân viên mới            |
| GET    | /employees/{id}      | Lấy nhân viên theo ID        |
| PUT    | /employees/{id}      | Cập nhật nhân viên           |
| DELETE | /employees/{id}      | Xoá nhân viên                |
| GET    | /employees/search    | Tìm kiếm nhân viên           |

### Departments

| Method | Endpoint                        | Mô tả                                   |
|--------|---------------------------------|-----------------------------------------|
| GET    | /departments                    | Lấy danh sách phòng ban                 |
| POST   | /departments                    | Tạo phòng ban mới                       |
| GET    | /departments/{id}/employees     | Lấy danh sách nhân viên theo phòng ban  |

---

## Middleware

- **LoggingMiddleware:** Ghi log method, URL, thời gian xử lý mỗi request ra console.

Thứ tự áp dụng:

```
Request → LoggingMiddleware → Handler
```

---

## Data Models

### Employee

```json
{
  "id": 1,
  "name": "Nguyen Van A",
  "age": 25,
  "position": "Engineer",
  "salary": 1500.00,
  "department_id": 1,
  "department_name": "Engineering"
}
```

### Department

```json
{
  "id": 1,
  "name": "Engineering"
}
```

---

## Export

Khi khởi động, server tự động export dữ liệu nhân viên ra:
- `exports/employees.csv`
- `exports/employees.json`
