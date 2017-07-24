package definition

import (
	"fmt"
	"gxdoc/util"
)

//HubReference is schema used by Link and Satelink to describe reference to hub
type HubReference struct {
	HubName  string
	Revision int
}

// GetDbTableName is to get equivalence database table name
func (hubRef *HubReference) GetDbTableName() string {
	return fmt.Sprintf("hub_%s_rev%d", util.ToSnakeCase(hubRef.HubName), hubRef.Revision)
}

// GetHashKey is to get equivalence database hash key table column name
func (hubRef *HubReference) GetHashKey() string {
	return fmt.Sprintf("%s_hash_key", util.ToSnakeCase(hubRef.HubName))
}
