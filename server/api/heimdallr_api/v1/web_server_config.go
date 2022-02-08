package v1

import metav1 "gin-vue-admin/internal/pkg/meta/v1"

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
