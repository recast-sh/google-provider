package impl

import (
	"encoding/base64"

	"recast.sh/v0/core"
	ignition "recast.sh/v0/ignition/impl"
)

// TODO wrong place!
type Path struct {
	Path     string
	Mode     uint32
	Uid      int
	Gid      int
	Contents core.Value
	Filters  []core.ValueFilter
}

func (p *Path) ignitionFile() ignition.File {
	var f = ignition.File{
		Filesystem: "root", // TODO configure!
		Path:       ignition.Path(p.Path),
		Mode:       ignition.FileMode(p.Mode),
		User:       ignition.FileUser{Id: p.Uid},
		Group:      ignition.FileGroup{Id: p.Gid},
	}

	if v, ok := p.Contents.(core.URLValue); ok {
		f.Contents.Source = ignition.Url(v.URL)
		f.Contents.Compression = ignition.Compression(v.Compression)
		if v.VerificationHashSum != "" {
			f.Contents.Verification.Hash = &ignition.Hash{
				Function: v.VerificationHashFunction,
				Sum:      v.VerificationHashSum,
			}
		}
	} else {
		value := p.Contents
		for _, f := range p.Filters {
			value = f(value)
		}
		if value != nil {
			if v := value.String(); v != "" {
				f.Contents.Source = ignition.Url{
					Scheme: "data",
					Opaque: ";base64," + base64.StdEncoding.EncodeToString([]byte(v)),
				}
			}
		}
	}
	return f
}
