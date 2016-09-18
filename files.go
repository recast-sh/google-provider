package google

import (
	"net/url"

	"recast.sh/v0/core"

	google "recast.sh/v0/provider/google/impl"
)

var filesScope *googleVM

func (vm *googleVM) Files(fn func()) core.BaseVM {
	if filesScope != nil {
		panic(core.ErrFilesNested)
	}
	filesScope = vm
	defer func() {
		filesScope = nil
	}()
	fn()
	return vm
}

func PathFromFile(p, file string) core.BaseFile {
	return Path(p).Contents(core.File(file))
}

func PathFromURL(p, u string) core.BaseURLFile {
	value, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	inner := filesScope.AddFile(google.Path{
		Path: p,
		Mode: 0644,
		Uid:  0,
		Gid:  0,
		Contents: core.URLValue{
			URL: *value,
		},
	})
	return &googleURL{
		Path: inner,
	}
}

func Path(p string) core.BaseFile {
	return &googlePath{
		Path: filesScope.AddFile(google.Path{
			Path: p,
			Mode: 0644,
			Uid:  0,
			Gid:  0,
		}),
	}
}

type googlePath struct {
	*google.Path
}

// TODO think about adding FileSystem ie "root" support

func (p *googlePath) GetPath() string {
	return p.Path.Path
}

func (p *googlePath) Mode(m uint32) core.BaseFile {
	p.Path.Mode = m
	return p
}

func (p *googlePath) Uid(i int) core.BaseFile {
	p.Path.Uid = i
	return p
}

func (p *googlePath) Gid(i int) core.BaseFile {
	p.Path.Gid = i
	return p
}

func (p *googlePath) Contents(c interface{}) core.BaseFile {
	var ok bool
	if p.Path.Contents, ok = c.(core.Value); !ok {
		switch v := c.(type) {
		case string:
			p.Path.Contents = core.StringValue(v)
		default:
			panic(core.ErrExpectedValue)
		}
	}
	return p
}

func (p *googlePath) Filter(fn core.ValueFilter) core.BaseFile {
	p.Path.Filters = append(p.Path.Filters, fn)
	return p
}

type googleURL struct {
	*google.Path
}

func (p *googleURL) GetPath() string {
	return p.Path.Path
}

func (p *googleURL) Mode(m uint32) core.BaseFile {
	p.Path.Mode = m
	return p
}

func (p *googleURL) Uid(i int) core.BaseFile {
	p.Path.Uid = i
	return p
}

func (p *googleURL) Gid(i int) core.BaseFile {
	p.Path.Gid = i
	return p
}

func (p *googleURL) Compression(t string) core.BaseURLFile {
	c := p.Path.Contents.(core.URLValue)
	c.Compression = t
	p.Path.Contents = c
	return p
}
func (p *googleURL) Verification(fn, sum string) core.BaseURLFile {
	c := p.Path.Contents.(core.URLValue)
	c.VerificationHashFunction = fn
	c.VerificationHashSum = sum
	p.Path.Contents = c
	return p
}

func (p *googleURL) Contents(c interface{}) core.BaseFile {
	var ok bool
	if p.Path.Contents, ok = c.(core.Value); !ok {
		switch v := c.(type) {
		case string:
			p.Path.Contents = core.StringValue(v)
		default:
			panic(core.ErrExpectedValue)
		}
	}
	return p
}

func (p *googleURL) Filter(fn core.ValueFilter) core.BaseFile {
	p.Path.Filters = append(p.Path.Filters, fn)
	return p
}
