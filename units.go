package google

import (
	"path"

	"recast.sh/v0/core"

	google "recast.sh/v0/provider/google/impl"
)

var unitsScope *googleVM

func (vm *googleVM) Units(fn func()) core.BaseVM {
	if unitsScope != nil {
		panic(core.ErrUnitsNested)
	}
	unitsScope = vm
	defer func() {
		unitsScope = nil
	}()
	fn()
	return vm
}

func Unit(name string) core.BaseUnit {
	return &googleUnit{
		Unit: unitsScope.AddUnit(google.Unit{
			Name:   name,
			Enable: true,
		}),
	}
}

func UnitFromFile(file string) core.BaseUnit {
	return Unit(path.Base(file)).
		Contents(core.File(file))
}

func UnitsFromFile(files ...string) []core.BaseUnit {
	units := make([]core.BaseUnit, len(files))
	for i, file := range files {
		units[i] = UnitFromFile(file)
	}
	return units
}

type googleUnit struct {
	*google.Unit
}

func (u *googleUnit) GetName() string {
	return u.Unit.Name
}

func (u *googleUnit) Enable(e bool) core.BaseUnit {
	u.Unit.Enable = e
	return u
}

func (u *googleUnit) Mask(m bool) core.BaseUnit {
	u.Unit.Mask = m
	return u
}

func (u *googleUnit) Contents(c interface{}) core.BaseUnit {
	var ok bool
	if u.Unit.Contents, ok = c.(core.Value); !ok {
		switch v := c.(type) {
		case string:
			u.Unit.Contents = core.StringValue(v)
		default:
			panic(core.ErrExpectedValue)
		}
	}
	return u
}

func (u *googleUnit) Filter(fn core.ValueFilter) core.BaseUnit {
	u.Unit.Filters = append(u.Unit.Filters, fn)
	return u
}

var dropInsScope *googleUnit

func (unit *googleUnit) DropIns(fn func()) core.BaseUnit {
	if dropInsScope != nil {
		panic(core.ErrUnitsNested)
	}
	dropInsScope = unit
	defer func() {
		dropInsScope = nil
	}()
	fn()
	return unit
}

func UnitDropIn(name string) core.BaseUnitDropIn {
	return &googleUnitDropIn{
		UnitDropIn: dropInsScope.AddUnitDropIn(google.UnitDropIn{
			Name: name,
		}),
	}
}

type googleUnitDropIn struct {
	*google.UnitDropIn
}

func (u *googleUnitDropIn) GetName() string {
	return u.UnitDropIn.Name
}

func (u *googleUnitDropIn) Contents(c interface{}) core.BaseUnitDropIn {
	var ok bool
	if u.UnitDropIn.Contents, ok = c.(core.Value); !ok {
		switch v := c.(type) {
		case string:
			u.UnitDropIn.Contents = core.StringValue(v)
		default:
			panic(core.ErrExpectedValue)
		}
	}
	return u
}

func (u *googleUnitDropIn) Filter(fn core.ValueFilter) core.BaseUnitDropIn {
	u.UnitDropIn.Filters = append(u.UnitDropIn.Filters, fn)
	return u
}
