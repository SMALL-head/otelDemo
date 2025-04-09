package otelmodel

type TraceData struct {
	Batches []Batch `json:"batches"`
}

type Batch struct {
	Resource   Resource     `json:"resource"`
	ScopeSpans []ScopeSpans `json:"scopeSpans"`
}

type Resource struct {
	Attributes []Attribute `json:"attributes"`
}

type ScopeSpans struct {
	Scope InstrumentationScope `json:"scope"`
	Spans []Span               `json:"spans"`
}

type InstrumentationScope struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"` // 有些scope没有version
}

type Span struct {
	TraceID           string      `json:"traceId"`
	SpanID            string      `json:"spanId"`
	ParentSpanID      string      `json:"parentSpanId,omitempty"`
	Name              string      `json:"name"`
	Kind              string      `json:"kind"`
	StartTimeUnixNano string      `json:"startTimeUnixNano"`
	EndTimeUnixNano   string      `json:"endTimeUnixNano"`
	Attributes        []Attribute `json:"attributes"`
	Status            Status      `json:"status"`
}

type Attribute struct {
	Key   string `json:"key"`
	Value Value  `json:"value"`
}

type Value struct {
	StringValue string `json:"stringValue,omitempty"`
	IntValue    string `json:"intValue,omitempty"`
}

type Status struct {
	// 可根据需要补充 status 字段，例如 code 和 message 等
}
