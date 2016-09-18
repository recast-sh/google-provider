package impl

import (
	"recast.sh/v0/core"
	ignition "recast.sh/v0/ignition/impl"
)

type Unit struct {
	Name     string
	Enable   bool
	Mask     bool
	Contents core.Value
	Filters  []core.ValueFilter

	dropIns    []*UnitDropIn
	dropInsMap map[string]*UnitDropIn
}

func (unit *Unit) AddUnitDropIn(dropIn UnitDropIn) *UnitDropIn {
	if unit.dropInsMap == nil {
		unit.dropInsMap = map[string]*UnitDropIn{}
	} else {
		if _, exists := unit.dropInsMap[unit.Name]; exists {
			panic(core.Errorf("Duplicate Unit DropIn `%s`", unit.Name))
		}
	}
	unit.dropInsMap[unit.Name] = &dropIn
	unit.dropIns = append(unit.dropIns, &dropIn)
	return &dropIn
}

func (u *Unit) ignitionUnit() ignition.SystemdUnit {
	su := ignition.SystemdUnit{
		Name:   ignition.SystemdUnitName(u.Name),
		Enable: u.Enable,
		Mask:   u.Mask,
	}

	value := u.Contents
	for _, f := range u.Filters {
		value = f(value)
	}
	if value != nil {
		su.Contents = value.String()
	}
	for _, d := range u.dropIns {
		su.DropIns = append(su.DropIns, d.ignitionDropIn())
	}
	return su
}

type UnitDropIn struct {
	Name     string
	Contents core.Value
	Filters  []core.ValueFilter
}

func (u *UnitDropIn) ignitionDropIn() ignition.SystemdUnitDropIn {
	value := u.Contents
	for _, f := range u.Filters {
		value = f(value)
	}
	return ignition.SystemdUnitDropIn{
		Name:     ignition.SystemdUnitDropInName(u.Name),
		Contents: value.String(),
	}
}
