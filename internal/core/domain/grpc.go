package domain

type ProtobufAny struct {
	Type  string                 `json:"@type"`
	Extra map[string]interface{} `json:"-"`
}

type RpcStatus struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Details []ProtobufAny `json:"details"`
}
