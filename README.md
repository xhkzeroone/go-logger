# go-logger

Một thư viện logging mạnh mẽ và linh hoạt cho Go, được xây dựng trên nền tảng logrus với các tính năng nâng cao như định dạng tùy chỉnh, bảo mật thông tin nhạy cảm và hỗ trợ context.

## 🌟 Tính năng chính

### 1. **Định dạng log linh hoạt**
- Hỗ trợ pattern tùy chỉnh với các placeholder động
- Cấu hình qua file XML
- Hiển thị timestamp, level, file, line, function name
- Hỗ trợ request ID từ context

### 2. **Bảo mật thông tin nhạy cảm**
- Tự động che giấu email, số điện thoại, thẻ tín dụng
- Cấu hình regex patterns tùy chỉnh
- Hỗ trợ nhiều định dạng số điện thoại (quốc tế, nội địa)

### 3. **Tích hợp context**
- Tự động lấy request ID từ context
- Hỗ trợ distributed tracing
- Log với context-aware

### 4. **Cấu hình dễ dàng**
- File cấu hình XML đơn giản
- Fallback configuration khi không tìm thấy file
- Hỗ trợ multiple log levels

## 📦 Cài đặt

```bash
go get github.com/xhkzeroone/go-logger
```

## 🚀 Cách sử dụng

### 1. **Khởi tạo cơ bản**

```go
package main

import "github.com/xhkzeroone/go-logger/logger"

func main() {
    // Khởi tạo logger với cấu hình mặc định
    err := logger.Init()
    if err != nil {
        panic(err)
    }
    
    // Sử dụng logger
    logger.Log.Info("Ứng dụng đã khởi động")
    logger.Log.WithField("user", "admin").Info("User đăng nhập")
}
```

### 2. **Sử dụng với context**

```go
package main

import (
    "context"
    "github.com/xhkzeroone/go-logger/logger"
)

func main() {
    logger.Init()
    
    // Tạo context với request ID
    ctx := context.WithValue(context.Background(), "requestId", "req-123")
    
    // Log với context
    log := logger.WithContext(ctx)
    log.Info("Request được xử lý")
    log.WithField("status", "success").Info("Hoàn thành request")
}
```

### 3. **Kích hoạt bảo mật thông tin nhạy cảm**

```go
package main

import "github.com/xhkzeroone/go-logger/logger"

func main() {
    // Đăng ký formatter bảo mật
    logger.RegisterSensitiveMessageFormater()
    logger.Init()
    
    // Thông tin nhạy cảm sẽ được che giấu tự động
    logger.Log.Info("Email: user@example.com") // Output: Email: **@**.***
    logger.Log.Info("Phone: +84123456789")     // Output: Phone: +84*********
}
```

### 4. **Tùy chỉnh formatter**

```go
package main

import "github.com/xhkzeroone/go-logger/logger"

type CustomMessageFormatter struct{}

func (c *CustomMessageFormatter) Format(message string) string {
    return "[CUSTOM] " + message
}

func main() {
    // Đăng ký custom formatter
    logger.RegisterMessageFormater(&CustomMessageFormatter{})
    logger.Init()
    
    logger.Log.Info("Test message") // Output: [CUSTOM] Test message
}
```

### 5. **Log ra JSON và mask dữ liệu nhạy cảm**

Bạn có thể log ra JSON đồng thời vẫn che giấu (mask) các thông tin nhạy cảm như email, số điện thoại bằng cách sử dụng `JSONFormatter`.

#### Ví dụ:

```go
package main

import "github.com/xhkzeroone/go-logger/logger"

func main() {
    // Đăng ký formatter bảo mật để mask message
    logger.RegisterSensitiveMessageFormater()
    logger.Init()

    // Sử dụng formatter JSON có mask
    logger.Log.SetFormatter(&logger.JSONFormatter{
        MsgFormatter: logger.GetMessageFormater(),
        FunctionNameFormatter: logger.GetFunctionNameFormatter(),
    })

    logger.Log.WithField("logger", "main").Info("Email: user@example.com | Phone: +84123456789")
}
```

#### Kết quả log:
```json
{
  "timestamp": "2024-01-15 10:30:45",
  "level": "info",
  "message": "Email: **@**.*** | Phone: +84*********",
  "logger": "main"
}
```

## ⚙️ Cấu hình

### File cấu hình log (`log-config.xml`)

```xml
<logConfig>
    <timestampFormat>2006-01-02 15:04:05</timestampFormat>
    <pattern>%timestamp% | %level% | %requestId% | %file%:%line% | %function% | %message%</pattern>
    <level>info</level>
</logConfig>
```

**Các placeholder hỗ trợ:**
- `%timestamp%`: Thời gian log
- `%level%`: Mức độ log (INFO, WARN, ERROR, etc.)
- `%requestId%`: ID của request (từ context)
- `%file%`: Tên file gọi log
- `%line%`: Số dòng gọi log
- `%function%`: Tên function gọi log
- `%message%`: Nội dung message
- `%fieldName%`: Giá trị của field tùy chỉnh

### File cấu hình bảo mật (`sensitive-patterns.xml`)

```xml
<patterns>
    <pattern>
        <type>email</type>
        <regex>[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}</regex>
        <replacement>**@**.***</replacement>
    </pattern>
    <pattern>
        <type>phone_international</type>
        <regex>\+84\d{9,10}</regex>
        <replacement>+84*********</replacement>
    </pattern>
</patterns>
```

## 📋 Ví dụ output

```
2024-01-15 10:30:45 | INFO | req-123 | main.go:15 | main | Ứng dụng đã khởi động
2024-01-15 10:30:45 | INFO | req-123 | main.go:16 | main | Email: **@**.***
2024-01-15 10:30:45 | WARN | null | main.go:25 | someFunc | Gọi hàm someFunc
```

## ✅ Ưu điểm

1. **Dễ sử dụng**: API đơn giản, tương thích với logrus
2. **Linh hoạt**: Hỗ trợ nhiều pattern và formatter tùy chỉnh
3. **Bảo mật**: Tự động che giấu thông tin nhạy cảm
4. **Hiệu suất cao**: Sử dụng logrus - một trong những logger nhanh nhất cho Go
5. **Cấu hình linh hoạt**: Hỗ trợ file XML và fallback configuration
6. **Context-aware**: Tích hợp tốt với distributed systems
7. **Extensible**: Dễ dàng mở rộng với custom formatters

## ❌ Nhược điểm

1. **Phụ thuộc logrus**: Không độc lập hoàn toàn, phụ thuộc vào thư viện bên ngoài
2. **Cấu hình XML**: Có thể không quen thuộc với một số developer
3. **Learning curve**: Cần hiểu về patterns và formatters để tùy chỉnh nâng cao
4. **File-based config**: Không hỗ trợ cấu hình từ environment variables
5. **Limited output formats**: Chỉ hỗ trợ text format, không có JSON/structured logging
6. **No async logging**: Không hỗ trợ logging bất đồng bộ

## 🔧 Yêu cầu hệ thống

- Go 1.23.1 trở lên
- Dependencies: `github.com/sirupsen/logrus`

## 📄 License

MIT License

## 🤝 Đóng góp

Mọi đóng góp đều được chào đón! Vui lòng tạo issue hoặc pull request.

## 📞 Hỗ trợ

Nếu có vấn đề hoặc câu hỏi, vui lòng tạo issue trên GitHub repository.
