package cleaner

import (
	"app/base/models"
	"app/base/utils"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// ClusterCleaner represents job struct for cleaning old clusters
type ClusterCleaner struct {
	logger *logrus.Logger
	conn   *gorm.DB

	ClusterRetention uint
}

// NewClusterCleaner builds ClusterCleaner
func NewClusterCleaner(clusterRetentionDays uint) (*ClusterCleaner, error) {
	logger, err := utils.CreateLogger(utils.Cfg.LoggingLevel)
	if err != nil {
		return nil, fmt.Errorf("Cannot create logger: %s", err)
	}

	if !(clusterRetentionDays > 0) {
		return nil, errors.New("CLUSTER_RETENTION_DAYS env not set")
	}

	db, err := models.GetGormConnection(utils.GetDbURL(false))
	if err != nil {
		return nil, fmt.Errorf("Cannot connect to db: %s", err)
	}

	return &ClusterCleaner{logger, db, clusterRetentionDays}, nil
}

// fetchExpiredClusters fetches expired clusters by deadline
func (c *ClusterCleaner) fetchExpiredClusters(deadline time.Time) []models.ClusterLight {
	clusters := make([]models.ClusterLight, 0)
	c.conn.Model(&models.Cluster{}).Where("last_seen < ?", deadline).Find(&clusters)
	return clusters
}

// deleteCaches deletes all caches where clusterID is used
func (c *ClusterCleaner) deleteCaches(tx *gorm.DB, clusterIDs []int64) error {
	tx.Where("cluster_id IN ?", clusterIDs).Delete(&models.ClusterCveCache{})
	return c.conn.Error
}

// deleteClusterImage deletes cluster_image table
func (c *ClusterCleaner) deleteClusterImage(tx *gorm.DB, clusterIDs []int64) error {
	tx.Where("cluster_id IN ?", clusterIDs).Delete(&models.ClusterImage{})
	return c.conn.Error
}

// deleteClusters deleted clusters table
func (c *ClusterCleaner) deleteClusters(tx *gorm.DB, clusterIDs []int64) error {
	tx.Where("ID IN ?", clusterIDs).Delete(&models.ClusterLight{})
	return c.conn.Error
}

// Clean cleans out old clusters
func (c *ClusterCleaner) Clean() error {
	c.logger.Info("Starting cluster cleaning job")

	clusterDeadline := time.Now().Truncate(time.Hour * 24).UTC().Add((-time.Hour * 24) * time.Duration(c.ClusterRetention))
	expiredClusters := c.fetchExpiredClusters(clusterDeadline)

	c.logger.Infof("Attempting to remove %d clusters", len(expiredClusters))

	expiredClustersIDs := make([]int64, len(expiredClusters))
	for _, cluster := range expiredClusters {
		expiredClustersIDs = append(expiredClustersIDs, cluster.ID)
	}

	tx := c.conn.Begin()
	defer tx.Rollback()

	err := c.deleteCaches(tx, expiredClustersIDs)
	if err != nil {
		return err
	}

	err = c.deleteClusterImage(tx, expiredClustersIDs)
	if err != nil {
		return err
	}

	err = c.deleteClusters(tx, expiredClustersIDs)
	if err != nil {
		return err
	}

	c.logger.Info("Cluster cleaning job done")
	return tx.Commit().Error
}
