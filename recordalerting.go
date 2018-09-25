package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/brotherlogic/goserver"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pbgh "github.com/brotherlogic/githubcard/proto"
	pbgd "github.com/brotherlogic/godiscogs"
	pbg "github.com/brotherlogic/goserver/proto"
	"github.com/brotherlogic/goserver/utils"
	pbrc "github.com/brotherlogic/recordcollection/proto"
	pbro "github.com/brotherlogic/recordsorganiser/proto"
)

type rc interface {
	getRecordsInPurgatory() ([]*pbrc.Record, error)
	getLibraryRecords() ([]*pbrc.Record, error)
	getSaleRecords() ([]*pbrc.Record, error)
}

type prodRC struct{}

func (gh *prodRC) getLibraryRecords() ([]*pbrc.Record, error) {
	host, port, err := utils.Resolve("recordsorganiser")

	if err != nil {
		return []*pbrc.Record{}, err
	}

	conn, err := grpc.Dial(host+":"+strconv.Itoa(int(port)), grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		return []*pbrc.Record{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	client := pbro.NewOrganiserServiceClient(conn)
	resp, err := client.GetOrganisation(ctx, &pbro.GetOrganisationRequest{Locations: []*pbro.Location{&pbro.Location{Name: "Library Records"}}})

	if err != nil {
		return []*pbrc.Record{}, err
	}

	if len(resp.GetLocations()) != 1 {
		return []*pbrc.Record{}, fmt.Errorf("Too many locations returned: %v", len(resp.GetLocations()))
	}

	recs := make([]*pbrc.Record, 0)
	for _, loc := range resp.GetLocations()[0].GetReleasesLocation() {
		rec, err := gh.getRecord(loc.GetInstanceId())
		if err != nil {
			return []*pbrc.Record{}, err
		}
		recs = append(recs, rec)
	}

	return recs, nil
}

func (gh *prodRC) getRecordsInPurgatory() ([]*pbrc.Record, error) {
	host, port, err := utils.Resolve("recordcollection")

	if err != nil {
		return []*pbrc.Record{}, err
	}

	conn, err := grpc.Dial(host+":"+strconv.Itoa(int(port)), grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		return []*pbrc.Record{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	client := pbrc.NewRecordCollectionServiceClient(conn)
	recs, err := client.GetRecords(ctx, &pbrc.GetRecordsRequest{Filter: &pbrc.Record{Release: &pbgd.Release{FolderId: 1362206}}})

	if err != nil {
		return []*pbrc.Record{}, err
	}

	return recs.GetRecords(), nil
}

func (gh *prodRC) getSaleRecords() ([]*pbrc.Record, error) {
	host, port, err := utils.Resolve("recordcollection")

	if err != nil {
		return []*pbrc.Record{}, err
	}

	conn, err := grpc.Dial(host+":"+strconv.Itoa(int(port)), grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		return []*pbrc.Record{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	client := pbrc.NewRecordCollectionServiceClient(conn)
	recs, err := client.GetRecords(ctx, &pbrc.GetRecordsRequest{Filter: &pbrc.Record{Release: &pbgd.Release{FolderId: 488127}}})

	if err != nil {
		return []*pbrc.Record{}, err
	}

	return recs.GetRecords(), nil
}

func (gh *prodRC) getRecord(instanceID int32) (*pbrc.Record, error) {
	host, port, err := utils.Resolve("recordcollection")
	if err != nil {
		return &pbrc.Record{}, err
	}

	conn, err := grpc.Dial(host+":"+strconv.Itoa(int(port)), grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		return &pbrc.Record{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	client := pbrc.NewRecordCollectionServiceClient(conn)
	recs, err := client.GetRecords(ctx, &pbrc.GetRecordsRequest{Filter: &pbrc.Record{Release: &pbgd.Release{InstanceId: instanceID}}})

	if err != nil {
		return &pbrc.Record{}, err
	}

	if len(recs.GetRecords()) == 0 {
		return &pbrc.Record{}, fmt.Errorf("No records found %v", instanceID)
	}

	return recs.GetRecords()[0], nil
}

type gh interface {
	alert(r *pbrc.Record, text string) error
}

type prodGh struct{}

func (gh *prodGh) alert(r *pbrc.Record, text string) error {
	host, port, err := utils.Resolve("githubcard")

	if err != nil {
		return err
	}

	conn, err := grpc.Dial(host+":"+strconv.Itoa(int(port)), grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	client := pbgh.NewGithubClient(conn)
	if r != nil {
		_, err = client.AddIssue(ctx, &pbgh.Issue{Title: "Problematic Record", Body: fmt.Sprintf("%v - %v", text, r.GetRelease().Title), Service: "recordcollection"})
	} else {
		_, err = client.AddIssue(ctx, &pbgh.Issue{Title: "Problematic Record", Body: fmt.Sprintf("%v", text), Service: "recordcollection"})
	}
	return err
}

//Server main server type
type Server struct {
	*goserver.GoServer
	rc rc
	gh gh
}

// Init builds the server
func Init() *Server {
	s := &Server{GoServer: &goserver.GoServer{}}
	s.gh = &prodGh{}
	s.rc = &prodRC{}
	return s
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {
	// Do nothing
}

// ReportHealth alerts if we're not healthy
func (s *Server) ReportHealth() bool {
	return true
}

// Mote promotes/demotes this server
func (s *Server) Mote(ctx context.Context, master bool) error {
	return nil
}

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {
	return []*pbg.State{}
}

func main() {
	var quiet = flag.Bool("quiet", false, "Show all output")
	flag.Parse()

	//Turn off logging
	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	server := Init()
	server.PrepServer()
	server.Register = server
	server.RegisterRepeatingTask(server.alertForPurgatory, "alert_for_purgatory", time.Hour)
	server.RegisterRepeatingTask(server.alertForMisorderedMPI, "alert_for_misordered_mpi", time.Hour)
	server.RegisterServer("recordalerting", false)
	server.Log("Starting!")
	fmt.Printf("%v", server.Serve())
}
