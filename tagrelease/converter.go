package tagrelease

import log "github.com/sirupsen/logrus"

type Adapter interface {
	Version() *Version
	Revision() (string, error)
	Branch() (string, error)
}

type Version struct {
	major  int
	minor  int
	patch  int
	diff   int
	rev    string
	suffix string
}

type Converter struct {
	adapter  Adapter
	strategy Strategy
}

func NewConverter(adapter Adapter, strategy Strategy) *Converter {
	return &Converter{
		adapter:  adapter,
		strategy: strategy,
	}
}

var empty = Version{}

func (c *Converter) Detect() (v *Version) {
	v = c.adapter.Version()
	log.WithField("version", v).Debug("version detected")
	if *v == empty {
		log.Debug("empty version detected, use first release strategy")
		v.minor = 1
		return
	}
	//use increment strategy
	c.strategy(v)
	return
}

func among(elem string, stack []string) bool {
	for i := range stack {
		if stack[i] == elem {
			return true
		}
	}
	return false
}

func (c *Converter) ReleaseKind() string {
	branch, _ := c.adapter.Branch()
	var kind string
	switch {
	case among(branch, GlobalConfig.Branches.Master):
		kind = "rc"
	case among(branch, GlobalConfig.Branches.Trunk):
		kind = "b"
	default:
		kind = "a"
	}
	log.WithField("kind", kind).Debug("calculated release kind")
	return kind
}

func (c *Converter) Revision() string {
	r, err := c.adapter.Revision()
	if err != nil {
		log.WithError(err).Debug("failed to detect revision")
	}
	return r
}
