package tasks

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

// See http://www.linuxfoundation.org/collaborate/workgroups/networking/netem
type ControlNetOptions struct {
	Type    string
	Timeout string // Times may be suffixed with ms,s,m,h

	// slow: tc qdisc add dev eth0 root netem delay 50ms 10ms distribution normal
	Delay          string
	DelayVariation string

	// flaky: tc qdisc add dev eth0 root netem loss 20% 75%
	Loss            string
	LossCorrelation string

	// reset: tc qdisc del dev eth0 root
}

func (ControlNetOptions) _private() {}

type ControlNetTask struct {
	cmdRunner boshsys.CmdRunner
	opts      ControlNetOptions
}

func NewControlNetTask(cmdRunner boshsys.CmdRunner, opts ControlNetOptions, _ boshlog.Logger) ControlNetTask {
	return ControlNetTask{cmdRunner, opts}
}

func (t ControlNetTask) Execute(stopCh chan struct{}) error {
	timeoutCh, err := NewOptionalTimeoutCh(t.opts.Timeout)
	if err != nil {
		return err
	}

	if len(t.opts.Delay) == 0 && len(t.opts.Loss) == 0 {
		return bosherr.Error("Must specify delay or loss")
	}

	ifaceNames, err := NonLocalIfaceNames()
	if err != nil {
		return err
	}

	opts := make([]string, 0, 16)

	if len(t.opts.Delay) > 0 {
		variation := t.opts.DelayVariation

		if len(variation) == 0 {
			variation = "10ms"
		}

		opts = append(opts, "delay", t.opts.Delay, variation, "distribution", "normal")
	}

	if len(t.opts.Loss) > 0 {
		correlation := t.opts.LossCorrelation

		if len(correlation) == 0 {
			correlation = "75%"
		}

		opts = append(opts, "loss", t.opts.Loss, correlation)
	}

	for _, ifaceName := range ifaceNames {
		err := t.configureInterface(ifaceName, opts)
		if err != nil {
			return err
		}
	}

	select {
	case <-timeoutCh:
	case <-stopCh:
	}

	for _, ifaceName := range ifaceNames {
		err := t.resetIface(ifaceName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t ControlNetTask) configureInterface(ifaceName string, opts []string) error {
	args := []string{"qdisc", "add", "dev", ifaceName, "root", "netem"}
	args = append(args, opts...)

	_, _, _, err := t.cmdRunner.RunCommand("tc", args...)
	if err != nil {
		return bosherr.WrapError(err, "Shelling out to tc netem")
	}

	return nil
}

func (t ControlNetTask) resetIface(ifaceName string) error {
	_, _, _, err := t.cmdRunner.RunCommand("tc", "qdisc", "del", "dev", ifaceName, "root")
	if err != nil {
		return bosherr.WrapError(err, "Resetting tc")
	}

	return nil
}
