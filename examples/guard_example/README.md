# Goblin Guard Example

Ví dụ này minh họa cách sử dụng Guard System trong Goblin Framework để bảo vệ các route và xác thực người dùng.

## Tính năng

1. **JWT Authentication**
   - Đăng nhập và nhận JWT token
   - Tự động xác thực token trong mọi request
   - Bỏ qua xác thực cho các route công khai

2. **Role-based Authorization**
   - AdminGuard để bảo vệ các route chỉ dành cho admin
   - Kiểm tra role của người dùng từ JWT token

3. **Middleware Integration**
   - LoggerMiddleware để ghi log cho mọi request
   - JWTGuard middleware để xác thực token

4. **Controller Organization**
   - AuthController: Xử lý đăng nhập
   - AdminController: Quản lý người dùng (chỉ admin)
   - UserController: Quản lý profile người dùng
   - PublicController: Cung cấp thông tin công khai

## Cấu trúc

```
guard_example/
├── adapter.go      # Adapter để tích hợp guard package với core package
├── main.go         # File chính chứa logic của ứng dụng
└── README.md       # Tài liệu này
```

## API Endpoints

1. **Public Routes**
   - `GET /ping`: Kiểm tra server
   - `GET /public`: Lấy thông tin công khai

2. **Authentication**
   - `POST /auth/login`: Đăng nhập và nhận JWT token

3. **Protected Routes**
   - `GET /users/profile`: Lấy thông tin profile người dùng (yêu cầu token)
   - `GET /admin/users`: Lấy danh sách users (yêu cầu token và quyền admin)

## Cách sử dụng

1. **Khởi động server**
   ```bash
   go run main.go
   ```

2. **Đăng nhập để lấy token**
   ```bash
   curl -X POST http://localhost:8080/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username": "admin", "password": "password"}'
   ```

3. **Truy cập các route được bảo vệ**
   ```bash
   # Lấy profile người dùng
   curl http://localhost:8080/users/profile \
     -H "Authorization: Bearer YOUR_TOKEN"

   # Lấy danh sách users (chỉ admin)
   curl http://localhost:8080/admin/users \
     -H "Authorization: Bearer YOUR_TOKEN"
   ```

## Lưu ý

1. Đây là ví dụ đơn giản, trong thực tế cần:
   - Mã hóa password
   - Lưu trữ users trong database thật
   - Xử lý refresh token
   - Thêm các biện pháp bảo mật khác

2. Guard System có thể mở rộng với:
   - Custom guards
   - Role-based access control (RBAC)
   - Permission-based access control (PBAC)
   - Rate limiting
   - IP filtering
   - và nhiều tính năng khác 