package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/brotherlogic/goserver/utils"
	pbrc "github.com/brotherlogic/recordcollection/proto"
)

func main() {
	ctx, cancel := utils.ManualContext("recordalerting_cli", time.Minute*30)
	defer cancel()

	conn, err := utils.LFDialServer(ctx, "recordalerting")
	if err != nil {
		log.Fatalf("Unable to dial: %v", err)
	}
	defer conn.Close()

	switch os.Args[1] {
	case "ping":
		id, err := strconv.Atoi(os.Args[2])
		sclient := pbrc.NewClientUpdateServiceClient(conn)
		_, err = sclient.ClientUpdate(ctx, &pbrc.ClientUpdateRequest{InstanceId: int32(id)})
		if err != nil {
			log.Fatalf("Error on GET: %v", err)
		}
	}
}
