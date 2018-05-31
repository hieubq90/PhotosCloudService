# Photos Cloud Service
**Photos Cloud Service** đơn giản là một service handle upload hình ảnh viết bằng Go lang.

## Các chức năng hiện có **(v0.1.1)**:
1. **Handle Upload image file:**
* Xác định chính xác file_type bằng cách đọc 512 bytes đầu tiên của file.
* Chỉ cho phép upload file_type  `image/png và image/jpeg`
* Giới hạn kích thước file ảnh thông qua file cấu hình _photos_cloud_service.yaml_
* Thay đổi nơi lưu trữ file thông qua file cấu hình _photos_cloud_service.yaml_
* Tự sinh file name bằng UUID.V4
* Tự tạo thư mục theo ngày
* Sử dụng NGINX để host thư mục ảnh
2. **Resize uploaded image**
* Sử dụng thư viện `bimg.v1 (Go) và libvips (C++)` cho tốc độ xử lý ảnh nhanh và tiết kiệm tài nguyên
* Có thể thêm bớt các kích thước resize ảnh trong file cấu hình _photos_cloud_service.yaml_ (mặc định đang có 320 & 720)

## Tính năng mới dự kiến:
* Có thẻ thêm chức năng serve_static thư mục chứa file đã upload lên. Nhưng nên sử dụng NGINX để serve static để đạt hiệu năng tốt hơn (phiên bản hiện tại đang sử dụng nginx để serve static)

## Build:
**Yêu cầu đã setup môi trường Go & cài đặt gb** ([https://getgb.io/](https://getgb.io/))
```
git clone [https://github.com/hieubq90/PhotosCloudService](https://github.com/hieubq90/PhotosCloudService) PhotosCloudService
cd PhotosCloudService/photos_cloud_service
gb build all
```

## Changes log:
**V0.1.0**:
1. **Handle Upload image file:**
* Xác định chính xác file_type bằng cách đọc 512 bytes đầu tiên của file.
* Chỉ cho phép upload file_type  `image/png và image/jpeg`
* Giới hạn kích thước file ảnh thông qua file cấu hình _photos_cloud_service.yaml_
* Thay đổi nơi lưu trữ file thông qua file cấu hình _photos_cloud_service.yaml_
* Tự sinh file name bằng UUID.V4
* Tự tạo thư mục theo ngày
 
