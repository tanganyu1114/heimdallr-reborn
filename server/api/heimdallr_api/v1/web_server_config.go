package v1

import (
	metav1 "github.com/tanganyu1114/heimdallr-reborn/internal/pkg/meta/v1"
)

type BifrostGroupMeta struct {
	metav1.UintObjectMeta `json:",inline"`
	Children              []*BifrostMeta `json:"children"`
}

type BifrostMeta struct {
	metav1.UintObjectMeta `json:",inline"`
	Children              []*WebSrvMeta `json:"children"`
}

type WebSrvMeta struct {
	metav1.StringObjectMeta `json:",inline"`
}
