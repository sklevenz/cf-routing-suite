package mongo

type simulator struct {
}

func (sim *simulator) ResetAll() ResultData {
	return ResultData{
		nil,
		0,
		"resetAll",
		nil,
		nil,
	}
}

func (sim *simulator) RecordRequest(requestData *RequestData) ResultData {
	return ResultData{
		nil,
		5,
		"request recorded",
		nil,
		requestData,
	}
}
