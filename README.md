# Permission Service

Service **Permission** là một phần của hệ thống Huyshop (Microservices architecture), được phát triển bằng Go (Golang) và giao tiếp chủ yếu qua gRPC. Service này đóng vai trò quan trọng trong việc quản lý phân quyền (RBAC - Role-Based Access Control), quản lý quyền truy cập của người dùng vào các trang và nhóm chức năng.

## 🚀 Tính năng chính

- **Quản lý Vai trò (Roles)**: Tạo mới, cập nhật, xóa, danh sách và kiểm tra vai trò (Role).
- **Quản lý Trang (Pages)**: Quản lý các trang hệ thống mà người dùng có thể truy cập (Thêm, Sửa, Xóa, Liệt kê).
- **Quản lý Nhóm quyền (Groups)**: Gom nhóm các quyền hạn để dễ dàng cấp phát và quản lý.
- **Giao tiếp gRPC**: Cung cấp các API gRPC tốc độ cao cho các service khác (như User/API gateway) để kiểm tra quyền hạn.

## 🛠 Yêu cầu hệ thống

- **Go**: 1.18+ (khuyến nghị)
- **Database**: MySQL (lưu trữ dữ liệu vai trò, trang, và nhóm)
- **Docker & Docker Compose** (tùy chọn cho việc deploy)

## ⚙️ Cấu hình môi trường

Tạo file `.env` ở thư mục gốc của project (có thể tham khảo file `.env.example`) với các cấu hình sau:

```env
# Cổng chạy gRPC server (Mặc định: 7000)
GRPC_PORT=7000

# Cấu hình kết nối MySQL
DB_PATH=user:password@tcp(host:port)

# Tên Database
DB_NAME=permission
```

## 📦 Hướng dẫn chạy và phát triển (Development)

Dự án sử dụng `Makefile` và `urfave/cli/v2` để quản lý các lệnh chạy.

### 1. Khởi tạo Database (Tạo bảng tự động)

Trước khi chạy ứng dụng lần đầu tiên, bạn cần khởi tạo các bảng trong Database:

```bash
make cdb
# hoặc chạy lệnh trực tiếp: go build && ./permission createDb
```

### 2. Khởi chạy Service

Chạy service ở chế độ bình thường:

```bash
make start
# hoặc chạy lệnh trực tiếp: go build && ./permission start
```

### 3. Build file thực thi (Binary)

Build ứng dụng thành file thực thi cho môi trường Linux (sử dụng cho Docker/Production):

```bash
make build
# Lệnh thực thi: CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o permission .
```

## 🐳 Khởi chạy với Docker

Dự án đã có sẵn `Dockerfile` (sử dụng base image `alpine:3.14`) tối ưu dung lượng cho môi trường production.

1. **Build image**:

```bash
make build # Build file binary (permission) trước
docker build -t huyshop/permission:latest .
```

2. **Chạy container**:

```bash
docker run -d \
  -p 7000:7000 \
  --name permission_service \
  --env-file .env \
  huyshop/permission:latest
```
