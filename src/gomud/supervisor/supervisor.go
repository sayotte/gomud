package supervisor

const NormalExit = "normal"
const FailsafeExit = "unknown"
const FailureExit = "failure"

type ExitStatus string

type Supervisable interface {
	Start()
	Stop()
	FailChan() <-chan ExitStatus
}

type Supervisor struct {
	strategy string
	children []Supervisable
	failChan chan ExitStatus
}

func NewSupervisor(strategy string) *Supervisor {
	s := &Supervisor{strategy: strategy}
	s.init()
	return s
}
func (s *Supervisor) init() {
	s.failChan = make(chan ExitStatus, 0)
}

// Supervisor should itself be a Supervisable, so we can make trees
// of them.
func (s *Supervisor) Start() {
	go s.supervisorLoop()
}
func (s *Supervisor) Stop() {
	return // stub
}
func (s *Supervisor) FailChan() <-chan ExitStatus {
	return s.failChan
}
func (s *Supervisor) supervisorLoop() {
	defer func() {
		s.failChan <- FailsafeExit
	}()

	// FIXME this needs to select over all of the
	// FIXME failchans belonging to its children.
	// FIXME It also needs a way to swallow a poison
	// FIXME pill, so we can force it to exit.
	for {
		select {
		case <-s.children[0].FailChan():
			continue
		}
	}
}
