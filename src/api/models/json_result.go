package models

// JSONSuccessResult represents the model of a synchronous call result
type JSONSuccessResult struct {
	Code          int         `json:"code" example:"200"`
	Message       string      `json:"message,omitempty" example:"Success"`
	Data          interface{} `json:"data,omitempty"`
	CorrelationId string      `json:"correlation_id,omitempty" example:"705e4dcb-3ecd-24f3-3a35-3e926e4bded5"`
}

// JSONAcceptedResult represents the model of an async call result. Requires implementation of state machine
type JSONAcceptedResult struct {
	Code          int         `json:"code" example:"202"`
	Id            string      `json:"id" example:"123-456-789-abc-def"`
	Message       string      `json:"message,omitempty" example:"Accepted"`
	Data          interface{} `json:"data,omitempty"`
	CorrelationId string      `json:"correlation_id,omitempty" example:"705e4dcb-3ecd-24f3-3a35-3e926e4bded5"`
}

// JSONFailureResult represents the model of a call result for a request which was deemed inappropriate by the server
type JSONFailureResult struct {
	Code          int         `json:"code" example:"400"`
	Data          interface{} `json:"data,omitempty"`
	Error         string      `json:"error,omitempty" example:"There was an error processing the request"`
	Stack         string      `json:"stacktrace,omitempty"`
	CorrelationId string      `json:"correlation_id,omitempty" example:"705e4dcb-3ecd-24f3-3a35-3e926e4bded5"`
}

// JSONNotFoundResult represents the model of a call result for a request which references a missing object
type JSONNotFoundResult struct {
	Code          int         `json:"code" example:"404"`
	Data          interface{} `json:"data,omitempty"`
	CorrelationId string      `json:"correlation_id,omitempty" example:"705e4dcb-3ecd-24f3-3a35-3e926e4bded5"`
}
