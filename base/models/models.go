package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

// Account table
type Account struct {
	ID    int64  `gorm:"type:bigint;primaryKey;autoIncrement"`
	OrgID string `gorm:"type:text"`
}

func (a Account) TableName() string {
	return "account"
}

// Arch table
type Arch struct {
	ID   int64  `gorm:"type:bigint;primaryKey;autoIncrement"`
	Name string `gorm:"type:text;not null"`
}

func (Arch) TableName() string {
	return "arch"
}

// Cluster table
type Cluster struct {
	ID                int64        `gorm:"type:bigint;primaryKey;autoIncrement"`
	UUID              uuid.UUID    `gorm:"type:uuid;unique"`
	DisplayName       string       `gorm:"type:text"`
	Status            string       `gorm:"type:text"`
	Type              string       `gorm:"type:text"`
	Version           string       `gorm:"type:text"`
	Provider          string       `gorm:"type:text"`
	AccountID         int64        `gorm:"type:bigint;not null"`
	CveCacheCritical  int32        `gorm:"type:int;not null;default:0"`
	CveCacheImportant int32        `gorm:"type:int;not null;default:0"`
	CveCacheModerate  int32        `gorm:"type:int;not null;default:0"`
	CveCacheLow       int32        `gorm:"type:int;not null;default:0"`
	LastSeen          time.Time    `gorm:"type:timestamp with time zone not null"`
	Workload          pgtype.JSONB `gorm:"type:jsonb"`
}

func (Cluster) TableName() string {
	return "cluster"
}

// ClusterLight table (Cluster without Workload JSONB)
type ClusterLight struct {
	ID                int64     `gorm:"type:bigint;primaryKey;autoIncrement"`
	UUID              uuid.UUID `gorm:"type:uuid;unique"`
	DisplayName       string    `gorm:"type:text"`
	Status            string    `gorm:"type:text"`
	Type              string    `gorm:"type:text"`
	Version           string    `gorm:"type:text"`
	Provider          string    `gorm:"type:text"`
	AccountID         int64     `gorm:"type:bigint;not null"`
	CveCacheCritical  int32     `gorm:"type:int;not null;default:0"`
	CveCacheImportant int32     `gorm:"type:int;not null;default:0"`
	CveCacheModerate  int32     `gorm:"type:int;not null;default:0"`
	CveCacheLow       int32     `gorm:"type:int;not null;default:0"`
	LastSeen          time.Time `gorm:"type:timestamp with time zone not null"`
}

func (ClusterLight) TableName() string {
	return "cluster"
}

// Repository table
type Repository struct {
	ID           int64     `gorm:"type:bigint;primaryKey;autoIncrement"`
	PyxisID      string    `gorm:"type:text;not null;unique"`
	ModifiedDate time.Time `gorm:"type:timestamp with time zone not null"`
	Registry     string    `gorm:"type:text;not null"`
	Repository   string    `gorm:"type:text;not null"`
}

func (i Repository) TableName() string {
	return "repository"
}

// Image table
type Image struct {
	ID                    int64     `gorm:"type:bigint;primaryKey;autoIncrement"`
	PyxisID               string    `gorm:"type:text;not null;unique"`
	ModifiedDate          time.Time `gorm:"type:timestamp with time zone not null"`
	ManifestListDigest    *string   `gorm:"type:text;index:image_manifest_list_digest_idx"`
	ManifestSchema2Digest *string   `gorm:"type:text;index:image_manifest_schema2_digest_idx"`
	DockerImageDigest     *string   `gorm:"type:text;index:image_docker_image_digest_idx"`
	ArchID                int64     `gorm:"type:bigint"`
}

func (i Image) TableName() string {
	return "image"
}

// RepositoryImage table
type RepositoryImage struct {
	RepositoryID int64 `gorm:"type:bigint;index:repository_image_repository_id_image_id_key"`
	ImageID      int64 `gorm:"type:bigint;index:repository_image_repository_id_image_id_key"`
}

func (ic RepositoryImage) TableName() string {
	return "repository_image"
}

// Cve table
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
	ExploitData  []byte     `gorm:"type:jsonb,default:null"`
}

func (c Cve) TableName() string {
	return "cve"
}

// ImageCve table
type ImageCve struct {
	ImageID int64 `gorm:"type:bigint;index:image_cve_image_id_cve_id_key"`
	CveID   int64 `gorm:"type:bigint;index:image_cve_image_id_cve_id_key"`
}

func (ic ImageCve) TableName() string {
	return "image_cve"
}

// ClusterImage table
type ClusterImage struct {
	ClusterID int64 `gorm:"type:bigint;index:cluster_image_cluster_id_image_id_key"`
	ImageID   int64 `gorm:"type:bigint;index:cluster_image_cluster_id_image_id_key"`
}

func (ci ClusterImage) TableName() string {
	return "cluster_image"
}

// ClusterCveCache table
type ClusterCveCache struct {
	ClusterID  int64 `gorm:"type:bigint;index:cluster_cve_cache_cluster_id_cve_id_key"`
	CveID      int64 `gorm:"type:bigint;index:cluster_cve_cache_cluster_id_cve_id_key"`
	ImageCount int32 `gorm:"type:int;not null;default:0"`
}

func (ccc ClusterCveCache) TableName() string {
	return "cluster_cve_cache"
}

// AccountCveCache table
type AccountCveCache struct {
	AccountID    int64 `gorm:"type:bigint;index:account_cve_cache_account_id_cve_id_key"`
	CveID        int64 `gorm:"type:bigint;index:account_cve_cache_account_id_cve_id_key"`
	ClusterCount int32 `gorm:"type:int;not null;default:0"`
	ImageCount   int32 `gorm:"type:int;not null;default:0"`
}

func (acc AccountCveCache) TableName() string {
	return "account_cve_cache"
}
