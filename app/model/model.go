package model

import (
	"database/sql"
	"time"
)

// Admin undefined
type Admin struct {
	ID          int64     `json:"id" gorm:"id"`
	Name        string    `json:"name" gorm:"name"`
	Password    string    `json:"password" gorm:"password"`
	CreatedTime time.Time `json:"created_time" gorm:"created_time"`
	UpdatedTime time.Time `json:"updated_time" gorm:"updated_time"`
}

// TableName 表名称
func (*Admin) TableName() string {
	return "admin"
}

type BookInfo struct {
	Id                 uint32 `gorm:"column:id;primary_key;NOT NULL;comment:'书的id'"`
	Uid                int64  `gorm:"column:uid;default:NULL;comment:'书的uid'"`
	BookName           string `gorm:"column:book_name;default:NULL;comment:'书名'"`
	Author             string `gorm:"column:author;default:NULL;comment:'作者'"`
	PublishingHouse    string `gorm:"column:publishing_house;default:NULL;comment:'出版社'"`
	Translator         string `gorm:"column:translator;default:NULL;comment:'译者'"`
	Num                int32  `gorm:"column:num;default:NULL;comment:'书的数量'"`
	PublishDate        string `gorm:"column:publish_date;default:NULL;comment:'出版时间'"`
	Pages              int32  `gorm:"column:pages;default:100;comment:'页数'"`
	ISBN               string `gorm:"column:ISBN;default:NULL;comment:'ISBN号码'"`
	Price              string `gorm:"column:price;default:1;comment:'价格'"`
	BriefIntroduction  string `gorm:"column:brief_introduction;default:;comment:'内容简介'"`
	AuthorIntroduction string `gorm:"column:author_introduction;default:;comment:'作者简介'"`
	ImgUrl             string `gorm:"column:img_url;default:NULL;comment:'封面地址'"`
	DelFlg             int32  `gorm:"column:del_flg;default:0;comment:'删除标识'"`
}

func (b *BookInfo) TableName() string {
	return "book_info"
}

// User undefined
type User struct {
	ID          int64     `json:"id" gorm:"id"`
	Name        string    `json:"name" gorm:"name"`
	Password    string    `json:"password" gorm:"password"`
	CreatedTime time.Time `json:"created_time" gorm:"created_time"`
	UpdatedTime time.Time `json:"updated_time" gorm:"updated_time"`
}

// TableName 表名称
func (*User) TableName() string {
	return "user"
}

// BookUser undefined
type BookUser struct {
	ID          int64     `json:"id" gorm:"id"`
	UserId      int64     `json:"user_id" gorm:"user_id"`
	BookId      int64     `json:"book_id" gorm:"book_id"`
	Status      int64     `json:"status" gorm:"status"`
	Time        int64     `json:"time" gorm:"time"`
	CreatedTime time.Time `json:"created_time" gorm:"created_time"`
	UpdatedTime time.Time `json:"updated_time" gorm:"updated_time"`
}

// TableName 表名称
func (*BookUser) TableName() string {
	return "book_user"
}

type MysqlConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
}

type Roles struct {
	RoleId      int64          `gorm:"column:role_id;primary_key;AUTO_INCREMENT;NOT NULL"`
	RoleName    string         `gorm:"column:role_name;NOT NULL"`
	Description sql.NullString `gorm:"column:description"`
}

func (r *Roles) TableName() string {
	return "roles"
}

type Permissions struct {
	PermissionId   int64          `gorm:"column:permission_id;primary_key;AUTO_INCREMENT;NOT NULL"`
	PermissionName string         `gorm:"column:permission_name;NOT NULL"`
	Description    sql.NullString `gorm:"column:description"`
}

func (p *Permissions) TableName() string {
	return "permissions"
}

type RolePermissions struct {
	RoleId       int64 `gorm:"column:role_id;primary_key;NOT NULL"`
	PermissionId int64 `gorm:"column:permission_id;NOT NULL"`
}

func (r *RolePermissions) TableName() string {
	return "role_permissions"
}

type UserRoles struct {
	UserId int64 `gorm:"column:user_id;primary_key;NOT NULL"`
	RoleId int64 `gorm:"column:role_id;NOT NULL"`
}

func (u *UserRoles) TableName() string {
	return "user_roles"
}

type Orders struct {
	Id           int32     `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	BookId       string    `gorm:"column:book_id;NOT NULL"`
	UserId       string    `gorm:"column:user_id;NOT NULL"`
	Quantity     int32     `gorm:"column:quantity;NOT NULL"`
	TotalAmount  string    `gorm:"column:total_amount;NOT NULL"`
	PurchaseTime time.Time `gorm:"column:purchase_time;NOT NULL"`
	Status       string    `gorm:"column:status;NOT NULL"`
}

func (o *Orders) TableName() string {
	return "orders"
}
