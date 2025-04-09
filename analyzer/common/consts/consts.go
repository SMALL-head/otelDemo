package consts

const (
	TempoTraceIDAPIWithDateTemplate = "/api/traces/%s?start=%s&end=%s" // 特定traceId查询
	TempoTraceIDAPITemplate         = "/api/traces/%s"                 // 特定traceId查询
	TempoSearchAPI                  = "/api/search"                    // hiding的traceQL进行全量的查询
)
