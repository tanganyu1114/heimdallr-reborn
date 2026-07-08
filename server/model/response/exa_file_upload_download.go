package response

import "github.com/tanganyu1114/heimdallr-reborn/model"

type ExaFileResponse struct {
	File model.ExaFileUploadAndDownload `json:"file"`
}
