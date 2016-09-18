package google

import (
	"sync"

	"golang.org/x/net/context"
	compute "google.golang.org/api/compute/v1"

	"recast.sh/v0/core"
	"recast.sh/v0/core/log"
	"recast.sh/v0/provider/google/impl"
)

func (p *googlePlan) Apply() {
	plan.l.Lock()
	if p.state != statePlanned {
		panic(core.ErrPlanAlreadyExecuted)
	}
	p.state = stateExecuting
	plan.l.Unlock()
	defer func() {
		p.state = stateExecuted
	}()

	if auth == nil {
		p.state = stateFailed
		panic(core.ErrMissingAuthentication)
	}

	if len(p.vms) > 0 {
		auth.Requires(impl.ComputeScope, impl.DevstorageFullControlScope)
	}

	if len(p.dnsZones) > 0 {
		auth.Requires(impl.CloudPlatformScope)
	}

	s, err := compute.New(auth.Config.Client(context.Background()))
	if err != nil {
		panic(err)
	}

	msg := "apply: Start"
	if core.DryRun {
		msg += " (dry run)"
	}
	log.Notice(msg) // return time taken

	wg := &sync.WaitGroup{}
	requests := make(chan ComputeRequest)
	for i := 0; i < core.MaxWorkers; i++ {
		computeWorker(s, requests, wg)
	}
	for _, vm := range p.vms {
		if vm.VM.Description == "" {
			log.Noticef(" Create [%s:%s]", vm.VM.Name, vm.VM.Zone)
		} else {
			log.Noticef(" Create [%s:%s] '%s'", vm.VM.Name, vm.VM.Zone, vm.VM.Description)
		}
		// TODO validate?

		r := vm.Create(auth.ProjectID)
		requests <- &r
	}
	close(requests)
	wg.Wait()

	msg = "apply: Finished"
	if core.DryRun {
		msg += " (dry run)"
	}
	log.Notice(msg) // return time taken
}

func computeWorker(service *compute.Service, requests chan ComputeRequest, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if r, more := <-requests; more {
				if core.DryRun {
					continue
				}
				r.Exec(service)
			} else {
				return
			}
		}
	}()
}

type ComputeRequest interface {
	Exec(service *compute.Service)
}
