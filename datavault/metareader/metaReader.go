package metareader

import "github.com/guinso/gxdoc/datavault/definition"

//MetaReader interface to read metadata of datavault from database
type MetaReader interface {
	GetHubDefinition(hubName string) (*definition.HubDefinition, error)
	GetLinkDefinition(linkName string) (*definition.LinkDefinition, error)
	GetSateliteDefinition(satName string) (*definition.SateliteDefinition, error)
}
