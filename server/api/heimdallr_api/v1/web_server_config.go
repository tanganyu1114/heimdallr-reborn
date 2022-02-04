package v1

import metav1 "gin-vue-admin/internal/pkg/meta/v1"

type WebSrvGroupMeta struct {
	metav1.LabelMeta `json:",inline"`
	Children         []*WebSrvHostMeta `json:"children"`
}

type WebSrvHostMeta struct {
	metav1.LabelMeta `json:",inline"`
	Children         []*WebSrvMeta `json:"children"`
}

type WebSrvMeta struct {
	metav1.LabelMeta `json:",inline"`
}
