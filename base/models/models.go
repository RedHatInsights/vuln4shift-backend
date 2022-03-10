package models

import (
	"time"

	"github.com/google/uuid"
)

// account	 table
type Account struct {
	ID   int64  `gorm:"type:bigint;primaryKey;autoIncrement"`
	Name string `gorm:"type:text"`
}

func (a Account) TableName() string {
	return "account"
}

// cluster table
type Cluster struct {
	ID                int64     `gorm:"type:bigint;primaryKey;autoIncrement"`
	UUID              uuid.UUID `gorm:"type:uuid;unique"`
	Status            string    `gorm:"type:text;not null"`
	Version           string    `gorm:"type:text;not null"`
	Provider          string    `gorm:"type:text"`
	AccountID         int64     `gorm:"type:bigint;not null"`
	CveCacheCritical  int32     `gorm:"type:int;not null;default:0"`
	CveCacheImportant int32     `gorm:"type:int;not null;default:0"`
	CveCacheModerate  int32     `gorm:"type:int;not null;default:0"`
	CveCacheLow       int32     `gorm:"type:int;not null;default:0"`
}

func (Cluster) TableName() string {
	return "cluster"
}

// image table
type Image struct {
	ID          int64  `gorm:"type:bigint;primaryKey;autoIncrement"`
	Digest      string `gorm:"type:text;not null;unique"`
	HealthIndex string `gorm:"type:char"`
}

func (i Image) TableName() string {
	return "image"
}

// cve table
type Cve struct {
	ID           int64      `gorm:"type:bigint;primaryKey;autoIncrement"`
	Name         string     `gorm:"type:text;not null;unique"`
	Description  string     `gorm:"type:text;not null"`
	Severity     Severity   `gorm:"type:severity;not null"`
	Cvss3Score   float32    `gorm:"type:numeric(5,3)"`
	Cvss3Metrics string     `gorm:"type:text"`
	Cvss2Score   float32    `gorm:"type:numeric(5,3)"`
	Cvss2Metrics string     `gorm:"type:text"`
	PublicDate   *time.Time `gorm:"type:timestamp with time zone null"`
	ModifiedDate *time.Time `gorm:"type:timestamp with time zone null"`
	RedhatURL    string     `gorm:"type:text"`
	SecondaryURL string     `gorm:"type:text"`
}

func (c Cve) TableName() string {
	return "cve"
}

// image_cve table
type ImageCve struct {
	ImageID int64 `gorm:"type:bigint;index:image_cve_image_id_cve_id_key"`
	CveID   int64 `gorm:"type:bigint;index:image_cve_image_id_cve_id_key"`
}

func (ic ImageCve) TableName() string {
	return "image_cve"
}

// cluster_image table
type ClusterImage struct {
	ClusterID int64 `gorm:"type:bigint;index:cluster_image_cluster_id_image_id_key"`
	ImageID   int64 `gorm:"type:bigint;index:cluster_image_cluster_id_image_id_key"`
}

func (ci ClusterImage) TableName() string {
	return "cluster_image"
}

// cluster_cve_cache table
type ClusterCveCache struct {
	ClusterID  int64 `gorm:"type:bigint;index:cluster_cve_cache_cluster_id_cve_id_key"`
	CveID      int64 `gorm:"type:bigint;index:cluster_cve_cache_cluster_id_cve_id_key"`
	ImageCount int32 `gorm:"type:int;not null;default:0"`
}

func (ccc ClusterCveCache) TableName() string {
	return "cluster_cve_cache"
}

// account_cve_cache
type AccountCveCache struct {
	AccountID    int64 `gorm:"type:bigint;index:account_cve_cache_account_id_cve_id_key"`
	CveID        int64 `gorm:"type:bigint;index:account_cve_cache_account_id_cve_id_key"`
	ClusterCount int32 `gorm:"type:int;not null;default:0"`
	ImageCount   int32 `gorm:"type:int;not null;default:0"`
}

func (acc AccountCveCache) TableName() string {
	return "account_cve_cache"
}
