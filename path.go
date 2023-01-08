package main

type Path interface {
	Parent() Path
	Path() string
	Label() string
	Final() bool
	HasParent() bool
}

type PathImpl struct {
	parent Path
	path   string
	label  string
	final  bool
}

func NewPath(parent Path, path, label string, final bool) Path {
	return &PathImpl{
		parent,
		path,
		label,
		final,
	}
}

func (p *PathImpl) Parent() Path {
	return p.parent
}

func (p *PathImpl) Path() string {
	return p.path
}

func (p *PathImpl) Label() string {
	return p.label
}

func (p *PathImpl) Final() bool {
	return p.final
}

func (p *PathImpl) HasParent() bool {
	return p.parent != nil
}
