package postgres

import (
	"gorm.io/gorm"
)

const TenantIDKey = "tenant_id"

// RegisterTenantScope registra un scope global para filtrar por tenant_id
func RegisterTenantScope(db *gorm.DB) {
	db.Callback().Query().Before("gorm:query").Register("tenant_scope", tenantScopeCallback)
	db.Callback().Create().Before("gorm:create").Register("tenant_scope_create", tenantScopeCreateCallback)
	db.Callback().Update().Before("gorm:update").Register("tenant_scope_update", tenantScopeUpdateCallback)
	db.Callback().Delete().Before("gorm:delete").Register("tenant_scope_delete", tenantScopeDeleteCallback)
}

func tenantScopeCallback(db *gorm.DB) {
	if tenantID, ok := db.Get(TenantIDKey); ok && tenantID != "" {
		db.Where("tenant_id = ?", tenantID)
	}
}

func tenantScopeCreateCallback(db *gorm.DB) {
	if tenantID, ok := db.Get(TenantIDKey); ok && tenantID != "" {
		db.Statement.SetColumn("tenant_id", tenantID)
	}
}

func tenantScopeUpdateCallback(db *gorm.DB) {
	if tenantID, ok := db.Get(TenantIDKey); ok && tenantID != "" {
		db.Where("tenant_id = ?", tenantID)
	}
}

func tenantScopeDeleteCallback(db *gorm.DB) {
	if tenantID, ok := db.Get(TenantIDKey); ok && tenantID != "" {
		db.Where("tenant_id = ?", tenantID)
	}
}
