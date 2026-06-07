package handling

import (
	"strings"

	"github.com/cRotermund/gameserver-manager/src/services/control-plane-api/lib/apimodels"
)

const ERR_UNKNOWN string = "UNKNOWN"

func ListServersApiError(err error) apimodels.APIError {
	//TODO try detect known errors on this path

	code := TryParseCommonErrors(err)

	return apimodels.APIError{
		Error: "Error listing servers",
		Code:  code,
	}
}

func ServerStateOperationError(verb string, err error) apimodels.APIError {
	//TODO try detect known errors on this path

	code := TryParseCommonErrors(err)

	var mbuilder strings.Builder
	mbuilder.WriteString("Error trying to ")
	mbuilder.WriteString(verb)
	mbuilder.WriteString(" the server")

	return apimodels.APIError{
		Error: mbuilder.String(),
		Code:  code,
	}
}

func ServerDetailError(err error) apimodels.APIError {
	//TODO try detect known errors on this path

	code := TryParseCommonErrors(err)

	return apimodels.APIError{
		Error: "Error getting server detail",
		Code:  code,
	}
}

func TryParseCommonErrors(err error) string {
	//handle any common errors

	return ERR_UNKNOWN
}
