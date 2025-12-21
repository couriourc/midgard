package collection

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/midgard/gateway/internal/database"
	"github.com/midgard/gateway/internal/openapi"
	"gorm.io/gorm"
)

// CollectionManager manages collections
type CollectionManager struct {
	db *gorm.DB
}

// NewCollectionManager creates a new collection manager
func NewCollectionManager(db *gorm.DB) *CollectionManager {
	return &CollectionManager{db: db}
}

// GetAllCollections gets all collections
func (cm *CollectionManager) GetAllCollections() ([]database.Collection, error) {
	var collections []database.Collection
	if err := cm.db.Preload("Endpoints").Find(&collections).Error; err != nil {
		return nil, err
	}
	return collections, nil
}

// GetCollection gets a collection by ID
func (cm *CollectionManager) GetCollection(id string) (*database.Collection, error) {
	var collection database.Collection
	if err := cm.db.Preload("Endpoints").First(&collection, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &collection, nil
}

// CreateCollection creates a new collection
func (cm *CollectionManager) CreateCollection(coll *database.Collection) error {
	// Check if prefix already exists
	var existing database.Collection
	if err := cm.db.Where("prefix = ?", coll.Prefix).First(&existing).Error; err == nil {
		return fmt.Errorf("prefix '%s' already exists", coll.Prefix)
	}
	
	if coll.ID == "" {
		coll.ID = uuid.New().String()
	}
	coll.CreatedAt = time.Now()
	coll.UpdatedAt = time.Now()
	return cm.db.Create(coll).Error
}

// CheckPrefixExists checks if a prefix already exists
func (cm *CollectionManager) CheckPrefixExists(prefix string, excludeID string) (bool, error) {
	var count int64
	query := cm.db.Model(&database.Collection{}).Where("prefix = ?", prefix)
	if excludeID != "" {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// UpdateCollection updates a collection
func (cm *CollectionManager) UpdateCollection(id string, coll *database.Collection) error {
	coll.UpdatedAt = time.Now()
	return cm.db.Model(&database.Collection{}).Where("id = ?", id).Updates(coll).Error
}

// DeleteCollection deletes a collection
func (cm *CollectionManager) DeleteCollection(id string) error {
	return cm.db.Delete(&database.Collection{}, "id = ?", id).Error
}

// ToggleCollection toggles the active state of a collection
func (cm *CollectionManager) ToggleCollection(id string) error {
	var coll database.Collection
	if err := cm.db.First(&coll, "id = ?", id).Error; err != nil {
		return err
	}
	coll.Active = !coll.Active
	coll.UpdatedAt = time.Now()
	return cm.db.Save(&coll).Error
}

// ImportOpenAPI imports OpenAPI spec and creates endpoints for a collection
func (cm *CollectionManager) ImportOpenAPI(collectionID string, openAPIURL string, openAPIJSON []byte) error {
	// Get collection
	coll, err := cm.GetCollection(collectionID)
	if err != nil {
		return fmt.Errorf("collection not found: %w", err)
	}

	// Import OpenAPI spec
	endpoints, err := openapi.ImportOpenAPI(openAPIURL, openAPIJSON, coll.BaseURL)
	if err != nil {
		return fmt.Errorf("failed to import OpenAPI: %w", err)
	}

	// Delete existing endpoints
	if err := cm.db.Where("collection_id = ?", collectionID).Delete(&database.Endpoint{}).Error; err != nil {
		return err
	}

	// Create new endpoints
	for i := range endpoints {
		endpoints[i].CollectionID = collectionID
		endpoints[i].CreatedAt = time.Now()
		endpoints[i].UpdatedAt = time.Now()
	}

	if len(endpoints) > 0 {
		if err := cm.db.Create(&endpoints).Error; err != nil {
			return err
		}
	}

	// Update collection's OpenAPI URL if provided
	if openAPIURL != "" {
		coll.OpenAPIURL = openAPIURL
		coll.UpdatedAt = time.Now()
		if err := cm.db.Save(coll).Error; err != nil {
			return err
		}
	}

	return nil
}

// GetCollectionByPrefix gets a collection by its prefix
func (cm *CollectionManager) GetCollectionByPrefix(prefix string) (*database.Collection, error) {
	var collection database.Collection
	if err := cm.db.Preload("Endpoints").First(&collection, "prefix = ? AND active = ?", prefix, true).Error; err != nil {
		return nil, err
	}
	return &collection, nil
}
