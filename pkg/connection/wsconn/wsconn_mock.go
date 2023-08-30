package wsconn

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/exasol/exasol-driver-go/pkg/types"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/mock"
)

type WebsocketConnectionMock struct {
	mock.Mock
}

func CreateWebsocketConnectionMock() *WebsocketConnectionMock {
	return &WebsocketConnectionMock{}
}

func (mock *WebsocketConnectionMock) SimulateSQLQueriesResponse(request interface{}, results []types.SqlQueryResponseResultSet) {
	marshalledResults := []json.RawMessage{}
	for _, r := range results {
		marshalledResults = append(marshalledResults, JsonMarshall(r))
	}
	mock.SimulateResponse(request, baseOKResponse(types.SqlQueriesResponse{NumResults: len(results), Results: marshalledResults}))
}
func (mock *WebsocketConnectionMock) SimulateOKResponse(request interface{}, response interface{}) {
	mock.SimulateResponse(request, baseOKResponse(response))
}

func baseOKResponse(payload interface{}) types.BaseResponse {
	return types.BaseResponse{Status: "ok", ResponseData: JsonMarshall(payload)}
}

func JsonMarshall(payload interface{}) json.RawMessage {
	data, err := json.Marshal(payload)
	if err != nil {
		panic(fmt.Errorf("failed to marshal data %v: %w", payload, err))
	}
	return data
}

func (mock *WebsocketConnectionMock) SimulateResponse(request interface{}, response interface{}) {
	requestMessage := JsonMarshall(request)
	mock.OnWriteTextMessage(requestMessage, nil)
	var responseMessage []byte
	if response != nil {
		responseMessage = JsonMarshall(response)
		mock.OnReadTextMessage(responseMessage, nil)
	}
	log.Printf("Simulate request: %s -> response: %s", string(requestMessage), string(responseMessage))
}

func (mock *WebsocketConnectionMock) OnWriteTextMessage(data []byte, returnedError error) {
	mock.On("WriteMessage", websocket.TextMessage, data).Return(returnedError).Once()
}

func (mock *WebsocketConnectionMock) OnReadTextMessage(data []byte, returnedError error) {
	mock.On("ReadMessage").Return(websocket.TextMessage, data, returnedError).Once()
}

func (mock *WebsocketConnectionMock) OnClose(returnedError error) {
	mock.On("Close").Return(returnedError)
}

func (mock *WebsocketConnectionMock) WriteMessage(messageType int, data []byte) error {
	log.Printf("ws.WriteMessage(%d, `%s`)", messageType, string(data))
	mockArgs := mock.Called(messageType, data)
	return mockArgs.Error(0)
}

func (mock *WebsocketConnectionMock) ReadMessage() (messageType int, response []byte, err error) {
	mockArgs := mock.Called()
	responseData := mockArgs.Get(1).([]byte)
	log.Printf("ws.ReadMessage() -> `%s`", string(responseData))
	return mockArgs.Int(0), responseData, mockArgs.Error(2)
}

func (mock *WebsocketConnectionMock) Close() error {
	mockArgs := mock.Called()
	return mockArgs.Error(0)
}
