package coaster

import (
	"context"
	"log"

	"kom.com/m/v2/src/kom.com/grpcCoaster"
)

type CoasterGrpcServerPort struct {
	grpcCoaster.UnimplementedCoasterServiceServer
	service CoasterService
}

func NewCoasterGrpcServerPort(theService CoasterService) CoasterGrpcServerPort {
	return CoasterGrpcServerPort{service: theService}
}

func (grpcS *CoasterGrpcServerPort) GetCoasters(context.Context, *grpcCoaster.Empty) (*grpcCoaster.CoastersMessage, error) {
	theCoasters := grpcS.service.getCoasters()
	log.Printf("Coasters ermittelt: %v", theCoasters)

	cma := &grpcCoaster.CoastersMessage{}

	for _, theCoaster := range theCoasters {
		grpcc := grpcCoaster.CoasterMessage{}
		grpcc.Name = theCoaster.Name
		grpcc.Id = theCoaster.ID
		grpcc.Manufacture = theCoaster.Manufacture
		grpcc.Height = uint32(theCoaster.Height)
		cma.Coasters = append(cma.Coasters, &grpcc)
	}

	return cma, nil
}

func (grpcS *CoasterGrpcServerPort) GetCoaster(ctx context.Context, in *grpcCoaster.CoasterIDMessage) (*grpcCoaster.CoasterMessage, error) {
	log.Printf("GetCoaster %v", in)

	theCoaster, err := grpcS.service.getCoaster(in.Id)
	log.Printf("Erg := %v", theCoaster)

	grpcc := grpcCoaster.CoasterMessage{}

	if err == nil {
		grpcc.Name = theCoaster.Name
		grpcc.Id = theCoaster.ID
		grpcc.Manufacture = theCoaster.Manufacture
		grpcc.Height = uint32(theCoaster.Height)
	}

	return &grpcc, err
}

func (grpcS *CoasterGrpcServerPort) CreateCoaster(ctx context.Context, in *grpcCoaster.CoasterMessage) (*grpcCoaster.Empty, error) {
	log.Printf("CreateCoaster %v", in)
	theCoaster := Coaster{}
	theCoaster.Height = int(in.Height)
	theCoaster.ID = in.Id
	theCoaster.Manufacture = in.Manufacture
	theCoaster.Name = in.Name
	err := grpcS.service.createCoaster(theCoaster)
	return nil, err
}
