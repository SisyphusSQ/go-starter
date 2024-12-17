package task

import "github.com/qiniu/qmgo/field"

type TaskList struct {
	field.DefaultField `bson:",inline"`

	ServiceName string `json:"serviceName" bson:"serviceName"`
	TicketID    int    `json:"ticketId" bson:"ticketId"`
	ElectionID  string `json:"electionId" bson:"electionId"`
	State       string `json:"state" bson:"state"`
	ConfigFile  string `json:"configFile" bson:"configFile"`
	SrcMongo    string `json:"srcMongo" bson:"srcMongo"`
	DestMongo   string `json:"destMongo" bson:"destMongo"`
	MetaData    string `json:"metaData" bson:"metaData"`
}
