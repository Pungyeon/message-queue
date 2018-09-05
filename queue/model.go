package queue

import "encoding/json"

// OperationType for queue operations
// Either Register or Message
type OperationType = string

// Operation determines the  action for the operation
// sent to the message queue
type Operation struct {
	Op    OperationType `json:"operation"`
	Value string        `json:"value"`
}

// NewOperationFromQueueBytes will return an operation struct parsed from a JSON byte array
func NewOperationFromQueueBytes(raw []byte) (Operation, error) {
	var operation Operation
	err := json.Unmarshal(raw[:len(raw)-1], &operation)
	if err != nil {
		return Operation{}, err
	}
	return operation, nil
}

// String returns a string value representing the Operation struct
func (operation *Operation) String() string {
	return `{` + `"operation":"` + operation.Op + `","value":"` + operation.Value + `"}`
}

// Bytes returns a byte array representing the Operation struct
func (operation *Operation) Bytes() []byte {
	return []byte(operation.String())
}

// QueueBytes returns a byte array with a nil terminator to send
// to the message queu
func (operation *Operation) QueueBytes() []byte {
	return append(operation.Bytes(), byte('\x00'))
}

const (
	Subscribe   OperationType = "subscribe"
	Message     OperationType = "message" // perhaps rename to push or send to be more clear
	Retrieve    OperationType = "retrieve"
	Notify      OperationType = "notify"
	Unsubscribe OperationType = "unsubscribe"
)
