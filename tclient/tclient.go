package tclient

import (
	"encoding/json"
	"log"
	"time"

	"fmt"
	"os/exec"

	"github.com/rocymp/atx-server/model"
	"github.com/rocymp/atx-server/proto"
	"github.com/rocymp/atx-server/util"
	"github.com/rocymp/zero"
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

	tc := &TClient{
		client: c,
		db:     db,
	}

	//注册处理函数
	tc.client.RegMessageHandler(tc.HandleMessage)

	return tc
}

func (tc *TClient) HandleMessage(msg *zero.Message) {
	rm := new(proto.RoomMessage)
	if msg.GetCMD() == int32(proto.StartRoomMessage) || msg.GetCMD() == int32(proto.StopRoomMessage) {
		json.Unmarshal(msg.GetData(), rm)

		log.Printf("%s room %d\n", rm.Command, rm.Rid)
		
		//更新设备状态为占用状态
		tc.db.DeviceUpdate("", proto.DeviceInfo{
			Using:        util.NewBool(true),
			UsingBeganAt: time.Now(),
		})
		args := make([]string,5)
		args[0] = "c:\\room.py"
		args[1] = "-c"
		args[2] = fmt.Sprintf("%s",rm.Command)
		args[3] = "-r"
		args[4] = fmt.Sprintf("%d",rm.Rid)
		c := exec.Command("python",args...)
		if err := c.Run();err != nil {
			fmt.Println("Error: ",err)
		}
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
