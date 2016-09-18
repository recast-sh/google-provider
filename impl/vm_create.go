package impl

import (
	"time"

	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
	"recast.sh/v0/core/log"
)

type CreateVMRequest struct {
	ProjectID string
	Zone      string
	*compute.Instance
}

func (req *CreateVMRequest) Exec(service *compute.Service) {
	_, err := service.Instances.Insert(req.ProjectID, req.Zone, req.Instance).Do()
	if err != nil {
		if alreadyExistsErr(err) {
			log.Debugf("[%s:%s] SKIPPING", req.Instance.Name, req.Zone)
			return
		}
		log.Debugf("Failed Instance: %s", req.Instance.Name)
		panic(err)
	}

	lastStatus := ""
	for {
		inst, err := service.Instances.Get(req.ProjectID, req.Zone, req.Instance.Name).Do()
		if err != nil {
			panic(err)
		}
		req.Instance = inst

		if req.Instance.Status == "RUNNING" {
			log.Infof("[%s:%s] RUNNING with ip %s", req.Instance.Name, req.Zone, req.GetExternalIP())
			return
		}
		if lastStatus != req.Instance.Status {
			lastStatus = req.Instance.Status
			log.Infof("[%s:%s] %s", req.Instance.Name, req.Zone, req.Instance.Status)
		}
		time.Sleep(1 * time.Second)
	}
}

func (req *CreateVMRequest) GetExternalIP() string {
	for _, network := range req.NetworkInterfaces {
		for _, ac := range network.AccessConfigs {
			if ac.NatIP != "" {
				return ac.NatIP
				break
			}
		}
	}
	return ""
}

func alreadyExistsErr(err error) bool {
	resp, ok := err.(*googleapi.Error)
	if !ok {
		return false
	}
	for _, e := range resp.Errors {
		if e.Reason == "alreadyExists" {
			return true
		}
	}
	return false
}
