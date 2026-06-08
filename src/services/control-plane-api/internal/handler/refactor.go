package handler

//////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////
//
//  TODO - move this somewhere else.  It shouldn't be here.  Error handling needs some love
//
//////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"strings"

	"github.com/cRotermund/gameserver-manager/src/services/control-plane-api/internal/models"
)

const ERR_UNKNOWN string = "UNKNOWN"

func ListServersApiError(err error) models.APIError {
	//TODO try detect known errors on this path

	code := TryParseCommonErrors(err)

	return models.APIError{
		Error: "Error listing servers",
		Code:  code,
	}
}

func ServerStateOperationError(verb string, err error) models.APIError {
	//TODO try detect known errors on this path

	code := TryParseCommonErrors(err)

	var mbuilder strings.Builder
	mbuilder.WriteString("Error trying to ")
	mbuilder.WriteString(verb)
	mbuilder.WriteString(" the server")

	return models.APIError{
		Error: mbuilder.String(),
		Code:  code,
	}
}

func ServerDetailError(err error) models.APIError {
	//TODO try detect known errors on this path

	code := TryParseCommonErrors(err)

	return models.APIError{
		Error: "Error getting server detail",
		Code:  code,
	}
}

func TryParseCommonErrors(err error) string {
	//handle any common errors

	return ERR_UNKNOWN
}
