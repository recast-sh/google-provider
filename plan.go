package google

import (
	"sync"

	"recast.sh/v0/core"
	"recast.sh/v0/core/log"

	"recast.sh/v0/provider/google/impl"
)

type planState int

const (
	stateReady planState = iota
	statePlanning
	statePlanned
	stateExecuting
	stateExecuted
	stateFailed
)

var plan = googlePlan{
	state:    stateReady,
	vmsNames: map[string]*googleVM{},
}

type googlePlan struct {
	l           sync.Mutex
	state       planState
	description string
	vms         []*googleVM
	vmsNames    map[string]*googleVM

	dnsZones []*googleDNSZone
}

func Plan(run func()) core.BasePlan {
	plan.l.Lock()
	if plan.state != stateReady {
		panic(core.ErrPlanNested)
	}
	plan.state = statePlanning
	plan.l.Unlock()
	defer func() {
		plan.l.Lock()
		plan.state = statePlanned
		plan.l.Unlock()
	}()

	// TODO authenticate here
	Authenticate()

	run()

	log.Debug("plan: done") // TODO finished in ...?

	return &plan
}

func (p *googlePlan) addVM(vm googleVM) GoogleVM {
	p.l.Lock()
	if p.state != statePlanning {
		panic(core.ErrMustUsePlan)
	}
	defer func() {
		p.l.Unlock()
	}()
	if _, exists := p.vmsNames[vm.Name]; exists {
		panic(core.Errorf("Duplicate VM `%s`", vm.Name))
	}
	p.vms = append(p.vms, &vm)
	p.vmsNames[vm.Name] = &vm
	return &vm
}

func (p *googlePlan) updateState(from, to planState) {
	plan.l.Lock()
	if plan.state != from {
		panic(core.ErrPlanNested)
	}
	plan.state = statePlanning
}

func (p *googlePlan) readClient() {
	if auth == nil {
		p.state = stateFailed
		panic(core.ErrMissingAuthentication)
	}

	auth.Requires(impl.CloudPlatformScope, impl.ComputeScope, impl.DevstorageFullControlScope)

	// s, err := compute.New(auth.Config.Client(context.Background()))
	// if err != nil {
	// 	panic(err)
	// }
}
