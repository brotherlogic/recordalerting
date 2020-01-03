package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/brotherlogic/goserver"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pbgh "github.com/brotherlogic/githubcard/proto"
	pbg "github.com/brotherlogic/goserver/proto"
	pbrc "github.com/brotherlogic/recordcollection/proto"
	pbro "github.com/brotherlogic/recordsorganiser/proto"
)

type ro interface {
	getLocation(ctx context.Context, name string) (*pbro.Location, error)
}

type prodRO struct {
	dial func(server string) (*grpc.ClientConn, error)
}

func (gh *prodRO) getLocation(ctx context.Context, name string) (*pbro.Location, error) {
	conn, err := gh.dial("recordsorganiser")
	if err != nil {
		return &pbro.Location{}, err
	}
	defer conn.Close()

	client := pbro.NewOrganiserServiceClient(conn)
	resp, err := client.GetOrganisation(ctx, &pbro.GetOrganisationRequest{Locations: []*pbro.Location{&pbro.Location{Name: name}}})

	if err != nil {
		return &pbro.Location{}, err
	}

	if len(resp.GetLocations()) != 1 {
		return &pbro.Location{}, fmt.Errorf("Too many locations returned: %v", len(resp.GetLocations()))
	}

	return resp.GetLocations()[0], nil
}

type rc interface {
	getRecord(ctx context.Context, instanceID int32) (*pbrc.Record, error)
	getLibraryRecords(ctx context.Context) ([]*pbrc.Record, error)
	getRecordsInFolder(ctx context.Context, folder int32) ([]int32, error)
}

type prodRC struct {
	dial func(server string) (*grpc.ClientConn, error)
}

func (gh *prodRC) getRecord(ctx context.Context, i int32) (*pbrc.Record, error) {
	conn, err := gh.dial("recordcollection")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pbrc.NewRecordCollectionServiceClient(conn)
	recs, err := client.GetRecord(ctx, &pbrc.GetRecordRequest{InstanceId: i})

	if err != nil {
		return nil, err
	}

	return recs.GetRecord(), nil
}

func (gh *prodRC) getLibraryRecords(ctx context.Context) ([]*pbrc.Record, error) {
	conn, err := gh.dial("recordsorganiser")
	if err != nil {
		return []*pbrc.Record{}, err
	}
	defer conn.Close()

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
		rec, err := gh.getRecord(ctx, loc.GetInstanceId())
		if err != nil {
			return []*pbrc.Record{}, err
		}
		recs = append(recs, rec)
	}

	return recs, nil
}

func (gh *prodRC) getRecordsInFolder(ctx context.Context, folder int32) ([]int32, error) {
	conn, err := gh.dial("recordcollection")
	if err != nil {
		return []int32{}, err
	}
	defer conn.Close()

	client := pbrc.NewRecordCollectionServiceClient(conn)
	recs, err := client.QueryRecords(ctx, &pbrc.QueryRecordsRequest{Query: &pbrc.QueryRecordsRequest_FolderId{folder}})

	if err != nil {
		return []int32{}, err
	}

	return recs.GetInstanceIds(), nil
}

type gh interface {
	alert(ctx context.Context, r *pbrc.Record, text string) error
}

type prodGh struct {
	dial func(server string) (*grpc.ClientConn, error)
}

func (gh *prodGh) alert(ctx context.Context, r *pbrc.Record, text string) error {
	conn, err := gh.dial("githubcard")
	if err != nil {
		return err
	}
	defer conn.Close()

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
	rc             rc
	gh             gh
	ro             ro
	invalidRecords int
}

// Init builds the server
func Init() *Server {
	s := &Server{GoServer: &goserver.GoServer{}}
	s.gh = &prodGh{s.DialMaster}
	s.rc = &prodRC{s.DialMaster}
	s.ro = &prodRO{s.DialMaster}
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

// Shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {
	return nil
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
	server.RegisterRepeatingTask(server.alertForMissingSaleID, "alert_for_missing_sale_id", time.Hour)
	server.RegisterRepeatingTask(server.alertForPurgatory, "alert_for_purgatory", time.Hour)
	server.RegisterRepeatingTask(server.alertForMisorderedMPI, "alert_for_misordered_mpi", time.Hour)
	server.RegisterRepeatingTask(server.alertForOldListeningBoxRecord, "alert_for_old_listening_box_record", time.Hour)
	server.RegisterRepeatingTask(server.alertForOldListeningPileRecord, "alert_for_old_listening_pile_record", time.Hour)
	server.RegisterRepeatingTask(server.validateRecords, "validate_records", time.Hour)

	err := server.RegisterServerV2("recordalerting", false, false)
	if err != nil {
		return
	}

	fmt.Printf("%v", server.Serve())
}
