package main

type Path interface {
	Parent() Path
	LocalPath() string
	GlobalPath() string
	Label() string
	Final() bool
	HasParent() bool
}

type PathImpl struct {
	parent     Path
	localPath  string
	globalPath string
	label      string
	final      bool
}

func NewPath(parent Path, localPath, globalPath, label string, final bool) Path {
	return &PathImpl{
		parent,
		localPath,
		globalPath,
		label,
		final,
	}
}

func (p *PathImpl) Parent() Path {
	return p.parent
}

func (p *PathImpl) LocalPath() string {
	return p.localPath
}

func (p *PathImpl) GlobalPath() string {
	return p.globalPath
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
