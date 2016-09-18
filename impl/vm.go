package impl

import (
	"encoding/json"
	"fmt"

	compute "google.golang.org/api/compute/v1"
	"recast.sh/v0/core"
	"recast.sh/v0/core/log"
	ignition "recast.sh/v0/ignition/impl"
)

type VM struct {
	Name        string
	Description string
	Tags        []string
	Zone        string
	MachineType string
	Preemptable bool
	Disks       []Disk

	files    []*Path
	filesMap map[string]*Path

	units    []*Unit
	unitsMap map[string]*Unit

	users    []*User
	usersMap map[string]*User

	// TODO add group
}

func (vm *VM) AddFile(path Path) *Path {
	if vm.filesMap == nil {
		vm.filesMap = map[string]*Path{}
	} else {
		if _, exists := vm.filesMap[path.Path]; exists {
			panic(core.Errorf("Duplicate File `%s`", path.Path))
		}
	}
	vm.filesMap[path.Path] = &path
	vm.files = append(vm.files, &path)
	return &path
}

func (vm *VM) AddUnit(unit Unit) *Unit {
	if vm.unitsMap == nil {
		vm.unitsMap = map[string]*Unit{}
	} else {
		if _, exists := vm.unitsMap[unit.Name]; exists {
			panic(core.Errorf("Duplicate Unit `%s`", unit.Name))
		}
	}
	vm.unitsMap[unit.Name] = &unit
	vm.units = append(vm.units, &unit)
	return &unit
}

func (vm *VM) AddUser(user User) *User {
	if vm.usersMap == nil {
		vm.usersMap = map[string]*User{}
	} else {
		if _, exists := vm.usersMap[user.Name]; exists {
			panic(core.Errorf("Duplicate User `%s`", user.Name))
		}
	}
	vm.usersMap[user.Name] = &user
	vm.users = append(vm.users, &user)
	return &user
}

type Disk struct {
	AutoDelete bool
	Boot       bool
	Type       DiskType
	Image      string
	SizeGb     int64
}

func NewDisk() Disk {
	return Disk{
		AutoDelete: true,
		Boot:       true,
		Type:       PersistentDisk,
	}
}

type DiskType int

const (
	PersistentDisk DiskType = iota
	ScratchDisk
)

func (t DiskType) String() string {
	switch t {
	case PersistentDisk:
		return "PERSISTENT"
	case ScratchDisk:
		return "SCRATCH"
	default:
		return "" // TODO panic?
	}
}

func (vm *VM) Create(projectID string) CreateVMRequest {
	return CreateVMRequest{
		ProjectID: projectID,
		Zone:      vm.Zone,
		Instance: &compute.Instance{
			Name:              vm.Name,
			Description:       vm.Description,
			Tags:              &compute.Tags{Items: vm.Tags},
			MachineType:       fmt.Sprint("https://www.googleapis.com/compute/v1/projects/", projectID, "/zones/", vm.Zone, "/machineTypes/", vm.MachineType),
			Disks:             vm.disks(),
			NetworkInterfaces: vm.networks(),
			ServiceAccounts:   vm.serviceAccounts(),
			Scheduling:        vm.scheduling(),
			Metadata:          vm.metadata(),
		},
	}
}

func (vm *VM) scheduling() *compute.Scheduling {
	if vm.Preemptable {
		return &compute.Scheduling{
			Preemptible:      true,
			AutomaticRestart: false,
		}
	} else {
		return &compute.Scheduling{
			Preemptible:      false,
			AutomaticRestart: true,
		}
	}
}

func (vm *VM) serviceAccounts() []*compute.ServiceAccount {
	return []*compute.ServiceAccount{} // TODO?
}

func (vm *VM) metadata() *compute.Metadata {
	hasData := false
	userData := ignition.Config{
		Ignition: ignition.Ignition{Version: ignition.IgnitionVersion{Major: 2}},
	}
	if len(vm.files) > 0 {
		hasData = true
		log.Debugf("  Files:")
		for _, f := range vm.files {
			log.Debugf("  - %s", f.Path)
			userData.Storage.Files = append(userData.Storage.Files, f.ignitionFile())
		}
	}
	if len(vm.units) > 0 {
		hasData = true
		log.Debugf("  Units:")
		for _, u := range vm.units {
			log.Debugf("  - %s", u.Name)
			userData.Systemd.Units = append(userData.Systemd.Units, u.ignitionUnit())
		}
	}
	if len(vm.users) > 0 {
		hasData = true
		log.Debugf("  Users:")
		for _, u := range vm.users {
			log.Debugf("  - %s", u.Name)
			userData.Passwd.Users = append(userData.Passwd.Users, u.ignitionUser())
		}
	}

	metadata := compute.Metadata{
		Items: []*compute.MetadataItems{}, // TODO: add user data!
	}

	if hasData {
		value, err := json.Marshal(&userData)
		if err != nil {
			panic(err)
		}
		v := string(value)
		metadata.Items = append(metadata.Items, &compute.MetadataItems{
			Key:   "user-data",
			Value: &v,
		})
	}
	return &metadata
}

func (vm *VM) disks() []*compute.AttachedDisk {
	disks := make([]*compute.AttachedDisk, len(vm.Disks))
	for i, disk := range vm.Disks {
		disks[i] = &compute.AttachedDisk{
			AutoDelete: disk.AutoDelete,
			Boot:       disk.Boot,
			Type:       disk.Type.String(),
			InitializeParams: &compute.AttachedDiskInitializeParams{
				SourceImage: disk.Image,
				DiskSizeGb:  disk.SizeGb,
			},
		}
	}
	return disks
}

func (vm *VM) networks() []*compute.NetworkInterface {
	return []*compute.NetworkInterface{
		{
			AccessConfigs: []*compute.AccessConfig{
				&compute.AccessConfig{Type: "ONE_TO_ONE_NAT", Name: "External NAT"},
			},
		},
	}
}
