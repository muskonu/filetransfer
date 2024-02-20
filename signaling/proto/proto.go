package proto

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

var validate = validator.New()

const offerJson = `{"source":"","message":""}`

const (
	CmdRegister = iota + 1
	CmdAnswer
	CmdOffer
	CmdCandidate
)

const (
	BadPayloadResponse = iota + 101
	CloseResponse
	OfferResponse
	AnswerResponse
	CandidateResponse
)

type Request struct { //with CmdOffer, CmdAnswer, CmdCandidate
	Cmd    int64  `json:"command"`
	Source string `json:"source" validate:"required,gte=8,lte=40"`
	Target string `json:"target" validate:"required,gte=8,lte=40"`
	Body   string `json:"-"` // carry offer, answer, candidate
}

type Response struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

type RegisterRequest struct { // with CmdRegister
	Cmd      int64  `json:"command"`
	SourceID string `json:"source" validate:"required,gte=8,lte=40"`
}

func NewResponse(code int64, message string) []byte {
	response, _ := json.Marshal(Response{code, message})
	return response
}

func PayloadRequest(payload []byte) (*Request, error) {
	request := &Request{}
	err := json.Unmarshal(payload, request)
	if err != nil {
		return nil, err
	}
	r := gjson.GetBytes(payload, "body")
	request.Body = r.String()
	err = validate.Struct(request)
	if err != nil {
		logrus.Errorf("parse payloadRequest error: %v requset:%s", err, string(payload))
		return nil, err
	}
	return request, nil
}

func PayloadRegisterRequest(payload []byte) (*RegisterRequest, error) {
	initRequest := &RegisterRequest{}
	err := json.Unmarshal(payload, initRequest)
	if err != nil {
		return nil, err
	}
	err = validate.Struct(initRequest)
	if err != nil {
		return nil, err
	}
	return initRequest, nil
}

func (req *Request) ToOfferString() string {
	v, _ := sjson.Set(offerJson, "source", req.Source)
	s, _ := sjson.Set(v, "message", req.Body)
	return s
}

func (req *Request) ToString() string {
	return req.Body
}
