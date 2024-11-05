# BookMall
精简的书城平台

# BookMall API Documentation

BookMall 是一个基于 Gin 框架的 RESTful API 项目，用于管理图书馆的书籍借阅系统。该项目使用了多种中间件和工具来处理路由、验证和安全性。

## 目录结构

- `/app`：包含应用程序的主要逻辑和模型。
  - `/logic`：包含处理业务逻辑的函数。
  - `/model`：定义数据库模型和数据结构。
  - `/tools`：包含辅助工具函数，例如生成验证码等。

- `/docs`：包含 Swagger API 文档的配置和说明。

- `/app/view`：包含 HTML 模板文件。

- `/app/images`：包含静态图片资源。

- `/app/middleware`：包含自定义中间件，例如用户验证和管理员验证。

## API 路由

### 登录界面

- `GET /index`：显示登录页面。
- `GET /wxack`：显示微信登录页面。
- `GET /BorrowingRecord`：显示借阅记录页面。
- `GET /RBACPermissionManagement`：显示基于角色的访问控制权限管理页面。

### 用户相关

- `POST /user/login`：用户登录。
- `POST /user/wxlogin`：微信登录。
- `POST /user/update-login-status`：更新登录状态。
- `GET /user/GetToken`：获取用户 Token。
- `POST /user/SendEmailCaptcha`：发送邮箱验证码。
- `POST /user/VerifyEmailCaptcha`：验证邮箱验证码。
- `POST /user/SendSMSCaptcha`：发送手机验证码。
- `POST /user/VerifySMSCaptcha`：验证手机验证码。
- `GET /user/wechat`：微信签名验证。
- `GET /user/wechat/login`：微信登录跳转。
- `GET /user/wechat/Callback`：微信回调处理。
- `GET /user/wechat/check_login`：检查微信登录状态。
- `GET /user/UserList`：管理员界面，显示用户列表。
- `POST /user/book_info/borrow`：借阅书籍。
- `POST /user/book_info/return`：归还书籍。
- `GET /user/UserLogin`：用户登录页面。
- `POST /user/PicUpload`：用户上传图片。

### 管理员相关

- `POST /admin/login`：管理员登录。
- `GET /admin/logout`：管理员登出。
- `GET /admin/AdminLogin`：管理员登录页面。
- `GET /admin/AdminList`：显示管理员列表。

### 游客相关

- `GET /visitor/VisitorLogout`：游客登出。
- `GET /visitor/VisitorLogin`：游客登录页面。

### 书籍相关

- `GET /book_info/GetBook`：获取书籍信息。
- `GET /book_info/list`：获取书籍列表。
- `GET /book_info/BooksBorrowingRecord`：获取书籍借阅记录。
- `POST /book_info/AddBook`：添加书籍。
- `PUT /book_info/SaveBook`：保存书籍信息。
- `DELETE /book_info/DelBook`：删除书籍。
- `GET /book_info/GetPaginatedBooks`：获取分页书籍列表。

### 借阅模块

- `GET /book_user/GetBookUserList`：获取借阅用户列表。
- `POST /book_user/UpdateBookUser`：更新借阅用户信息。

### 权限模块

- `GET /rights/GetRoles`：获取角色列表。
- `GET /rights/GetPermissions`：获取权限列表。
- `GET /rights/GetRole_Pre`：获取角色权限。
- `GET /rights/GetUer_Roles`：获取用户角色。
- `PUT /rights/UpdateRole_Pre`：更新角色权限。
- `POST /rights/UpdateUer_Roles`：更新用户角色权限。

### 支付模块

- `POST /alipay/pay`：支付宝支付。
- `GET /alipay/callback`：支付宝回调处理。
- `POST /alipay/notify`：支付宝支付通知。

### 验证码相关

- `GET /captcha`：生成并显示验证码。
- `POST /captcha/verify`：验证用户输入的验证码。

## 安装和运行

1. 克隆项目到本地：
    ```
    git clone "git@github.com:lbyyh/Bookmall.git"
    ```
2. 确保您已经安装了 Go 和 Gin 框架。
3. 运行程序：
    ```
    go run main.go
    ```
4. 程序默认监听在 `:8087` 端口。

## 访问 API

API 文档可以通过访问 `http://localhost:8087/swagger/index.html` 来查看。

## 注意事项

- 请确保您的数据库连接配置正确。
- 在生产环境中，您可能需要配置 HTTPS 以保证安全性。
- 请确保所有的依赖项都已经安装并配置正确。

