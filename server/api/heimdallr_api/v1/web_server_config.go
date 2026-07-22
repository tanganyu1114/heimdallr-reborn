package v1

type BifrostGroupMeta struct {
	UintObjectMeta `json:",inline"`
	Children       []*BifrostMeta `json:"children"`
}

type BifrostMeta struct {
	UintObjectMeta `json:",inline"`
	Children       []*WebSrvMeta `json:"children"`
}

type WebSrvMeta struct {
	StringObjectMeta `json:",inline"`
}
