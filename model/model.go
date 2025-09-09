package model

import (
	"gorm.io/gorm"
	"time"
)

// 交易类型枚举
type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

// 交易状态枚举
type TransactionStatus string

const (
	Valid   TransactionStatus = "valid"
	Deleted TransactionStatus = "deleted"
	Pending TransactionStatus = "pending"
)

// 分类类型枚举
type CategoryType string

const (
	CategoryIncome  CategoryType = "income"
	CategoryExpense CategoryType = "expense"
)

// 成员角色枚举
type MemberRole string

const (
	RoleAdmin  MemberRole = "admin"
	RoleMember MemberRole = "member"
	RoleViewer MemberRole = "viewer"
)

// 家庭表
type Family struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Members   []Member  `gorm:"foreignkey:FamilyID" json:"members,omitempty"`
}

// 成员表
type Member struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	FamilyID  uint       `json:"family_id" gorm:"index"`
	Family    Family     `json:"family,omitempty" gorm:"foreignKey:FamilyID"`
	Name      string     `gorm:"size:50;not null" json:"name"`
	Role      MemberRole `gorm:"type:ENUM('admin', 'member', 'viewer');default:'member'" json:"role"`
	Phone     string     `gorm:"size:20" json:"phone"`
	Email     string     `gorm:"size:100" json:"email"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	Status    int8       `gorm:"default:1" json:"status"` // 1=正常，0=已移除
}

// 分类表
type Category struct {
	ID        uint         `gorm:"primaryKey" json:"id"`
	Name      string       `gorm:"size:100;not null" json:"name"`
	Type      CategoryType `gorm:"type:ENUM('income', 'expense');not null" json:"type"`
	ParentID  *uint        `json:"parent_id"`
	Parent    *Category    `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Path      string       `gorm:"size:500;index" json:"path"`
	Level     int8         `json:"level"`
	SortOrder int          `json:"sort_order"`
	IsDeleted bool         `gorm:"default:false" json:"is_deleted"`
	CreatedAt time.Time    `json:"created_at" gorm:"autoCreateTime"`
	Children  []Category   `gorm:"foreignkey:ParentID" json:"children,omitempty"`
}

// 标签表
type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	FamilyID  uint      `json:"family_id" gorm:"index"`
	Family    Family    `json:"family,omitempty" gorm:"foreignKey:FamilyID"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Type      string    `gorm:"size:50" json:"type"` // merchant, occasion, area等
	Color     string    `gorm:"size:7" json:"color"` // #FF5733
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// 收支流水表
type Transaction struct {
	ID              uint              `gorm:"primaryKey" json:"id"`
	CreatedAt       time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time         `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt    `gorm:"index" json:"-"`
	FamilyID        uint              `json:"family_id" gorm:"index"`
	Family          Family            `json:"family,omitempty" gorm:"foreignKey:FamilyID"`
	MemberID        uint              `json:"member_id" gorm:"index"`
	Member          Member            `json:"member,omitempty" gorm:"foreignKey:MemberID"`
	Amount          float64           `gorm:"type:DECIMAL(12,2);not null" json:"amount"`
	Type            TransactionType   `gorm:"type:ENUM('income', 'expense');not null" json:"type"`
	CategoryID      uint              `json:"category_id" gorm:"index"`
	Category        Category          `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	TransactionTime time.Time         `gorm:"not null" json:"transaction_time"`
	Note            string            `gorm:"type:TEXT" json:"note"`
	ImageURL        string            `gorm:"size:500" json:"image_url"`
	Status          TransactionStatus `gorm:"type:ENUM('valid', 'deleted', 'pending');default:'valid'" json:"status"`
	PaymentMethod   string            `gorm:"size:50" json:"payment_method"` // 支付方式：现金、银行卡、支付宝、微信等
	Labels          []Tag             `gorm:"many2many:transaction_tags;" json:"labels"`
}

// 流水-标签关联表
type TransactionTag struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	TransactionID uint        `json:"transaction_id" gorm:"index"`
	Transaction   Transaction `json:"transaction,omitempty" gorm:"foreignKey:TransactionID"`
	TagID         uint        `json:"tag_id" gorm:"index"`
	Tag           Tag         `json:"tag,omitempty" gorm:"foreignKey:TagID"`
	CreatedAt     time.Time   `json:"created_at" gorm:"autoCreateTime"`
}
