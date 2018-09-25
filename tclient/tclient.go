package tclient

import (
	"encoding/json"
	"log"

	"github.com/rocymp/atx-server/model"
	"github.com/rocymp/atx-server/proto"
	"github.com/rocymp/zero"
	"os/exec"
	"fmt"
)

type TClient struct {
	client *zero.SocketClient
	db     *model.RdbUtils
}

func NewTClient(addr string, db *model.RdbUtils) *TClient {
	c := zero.NewSocketClient(addr, 3)
	if c == nil {
		log.Printf("Connect failed\n")
		return nil
	}

	c.Online()

	return &TClient{
		client: c,
		db:     db,
	}
}

func (tc *TClient) HandleMessage(msg *zero.Message) {
	rm := new(proto.RoomMessage)
	if msg.GetCMD() == int32(proto.StartRoomMessage) || msg.GetCMD() == int32(proto.StopRoomMessage) {
		json.Unmarshal(msg.GetData(),rm)
		exec.Command("python",fmt.Sprintf("d:\room.py -c %s -r %d",rm.Command,rm.Rid))
	}
}

func (tc *TClient) ReportDevices() {
	devices, err := tc.db.DeviceList()
	if err != nil {
		log.Printf("Get Device err %#v\n", err)
		return
	}

	b, err := json.Marshal(devices)
	tc.client.SendMessage(int32(proto.UpReportMessage), b)
}
