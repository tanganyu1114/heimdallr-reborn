package response

import "github.com/tanganyu1114/heimdallr-reborn/model"

type FilePathResponse struct {
	FilePath string `json:"filePath"`
}

type FileResponse struct {
	File model.ExaFile `json:"file"`
}
