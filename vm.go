package google

import (
	"sync"

	"recast.sh/v0/core"

	google "recast.sh/v0/provider/google/impl"
)

func VM(name string) GoogleVM {
	return plan.addVM(googleVM{
		VM: google.VM{
			Name:  name,
			Disks: []google.Disk{google.NewDisk()},
		},
	})
}

type GoogleVM interface {
	core.BaseVM
	Preemptable(preemptable bool) GoogleVM
	Tags(tags ...string) GoogleVM
	Zone(zone string) GoogleVM
	MachineType(t string) GoogleVM
	Image(url string, size int64) GoogleVM
	// Metadata(n string, s func() string) GoogleVM // TODO can this be generic?
}

type googleVM struct {
	l sync.Mutex
	google.VM
}

func (vm *googleVM) GetName() string {
	return vm.Name
}

func (vm *googleVM) Description(description string) core.BaseVM {
	vm.VM.Description = description
	return vm
}

func (vm *googleVM) Preemptable(preemptable bool) GoogleVM {
	vm.VM.Preemptable = preemptable
	return vm
}

func (vm *googleVM) Tags(tags ...string) GoogleVM {
	vm.VM.Tags = append(vm.VM.Tags, tags...)
	return vm
}

func (vm *googleVM) Zone(zone string) GoogleVM {
	vm.VM.Zone = zone
	return vm
}

func (vm *googleVM) MachineType(machineType string) GoogleVM {
	vm.VM.MachineType = machineType
	return vm
}

func (vm *googleVM) Image(url string, size int64) GoogleVM {
	vm.VM.Disks[0].Image = url
	vm.VM.Disks[0].SizeGb = size
	return vm
}

// func (vm *googleVM) Networks(fn func()) core.BaseVM {
// 	fn()
// 	return vm
// }

// // Google only?
// func (vm *googleVM) Disks(fn func()) core.BaseVM {
// 	fn()
// 	return vm
// }
