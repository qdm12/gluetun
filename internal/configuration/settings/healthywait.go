package settings

import (
	"time"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gotree"
)

type HealthyWait struct {
	// Initial is the initial duration to wait for the program
	// to be healthy before taking action.
	// It cannot be nil in the internal state.
	Initial *time.Duration
	// Addition is the duration to add to the Initial duration
	// after Initial has expired to wait longer for the program
	// to be healthy.
	// It cannot be nil in the internal state.
	Addition *time.Duration
}

func (h HealthyWait) validate() (err error) {
	return nil
}

func (h *HealthyWait) copy() (copied HealthyWait) {
	return HealthyWait{
		Initial:  gosettings.CopyPointer(h.Initial),
		Addition: gosettings.CopyPointer(h.Addition),
	}
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (h *HealthyWait) overrideWith(other HealthyWait) {
	h.Initial = gosettings.OverrideWithPointer(h.Initial, other.Initial)
	h.Addition = gosettings.OverrideWithPointer(h.Addition, other.Addition)
}

func (h *HealthyWait) setDefaults() {
	const initialDurationDefault = 6 * time.Second
	const additionDurationDefault = 5 * time.Second
	h.Initial = gosettings.DefaultPointer(h.Initial, initialDurationDefault)
	h.Addition = gosettings.DefaultPointer(h.Addition, additionDurationDefault)
}

func (h HealthyWait) String() string {
	return h.toLinesNode("Health").String()
}

func (h HealthyWait) toLinesNode(kind string) (node *gotree.Node) {
	node = gotree.New(kind + " wait durations:")
	node.Appendf("Initial duration: %s", *h.Initial)
	node.Appendf("Additional duration: %s", *h.Addition)
	return node
}

func (h *HealthyWait) read(r *reader.Reader) (err error) {
	h.Initial, err = r.DurationPtr(
		"HEALTH_VPN_DURATION_INITIAL",
		reader.RetroKeys("HEALTH_OPENVPN_DURATION_INITIAL"))
	if err != nil {
		return err
	}

	h.Addition, err = r.DurationPtr(
		"HEALTH_VPN_DURATION_ADDITION",
		reader.RetroKeys("HEALTH_OPENVPN_DURATION_ADDITION"))
	if err != nil {
		return err
	}

	return nil
}
