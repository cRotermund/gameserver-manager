package mapping

import (
	"github.com/cRotermund/gameserver-manager/src/services/control-plane-api/lib/apimodels"
	"github.com/cRotermund/gameserver-manager/src/services/control-plane-api/lib/awsutils"
)

func InstanceFlyweightToSummary(
	sif awsutils.ServerInstanceFlyweight,
	nameTagKey string,
	gameTagKey string) apimodels.ServerSummary {
	return apimodels.ServerSummary{
		ServerID: sif.Id,
		Name:     getTagValue(nameTagKey, sif.Tags),
		Game:     getTagValue(gameTagKey, sif.Tags),
		Specs:    apimodels.HardwareSpecs{}, //TODO
		Status:   apimodels.ServerStatus(sif.State),
	}
}

func InstanceDetailsToDetail(
	sid awsutils.ServerInstanceDetails,
	nameTagKey string,
	gameTagKey string) apimodels.ServerDetail {
	return apimodels.ServerDetail{
		ServerSummary:    InstanceFlyweightToSummary(sid.ServerInstanceFlyweight, nameTagKey, gameTagKey),
		IPAddress:        nil, //TODO
		ConnectedClients: nil, //TODO
		ResourceUsage:    nil, //TODO
	}
}

func getTagValue(tagkey string, tags map[string]string) string {
	value, ok := tags[tagkey]

	if ok {
		return value
	} else {
		return ""
	}
}
