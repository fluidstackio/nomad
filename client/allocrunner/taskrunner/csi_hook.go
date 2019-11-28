package taskrunner

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad/client/allocrunner/interfaces"
	"github.com/hashicorp/nomad/client/consul"
	"github.com/hashicorp/nomad/nomad/structs"
	"github.com/pkg/errors"
)

var _ interfaces.TaskPrestartHook = (*csiHook)(nil)

const (
	// the name of this hook, used in logs
	csiHookName = "consul_si"

	// csiBackoffBaseline is the baseline time for exponential backoff when
	// attempting to retrieve a Consul SI token
	csiBackoffBaseline = 5 * time.Second

	// csiBackoffLimit is the limit of the exponential backoff when attempting
	// to retrieve a Consul SI token
	csiBackoffLimit = 3 * time.Minute

	// csiTokenFile is the name of the file holding the Consul SI token inside
	// the task's secret directory
	csiTokenFile = "si_token"
)

type csiHookConfig struct {
	alloc    *structs.Allocation
	task     *structs.Task
	siClient consul.ServiceIdentityAPI
	logger   log.Logger
}

// Consul Service Identity hook for managing SI tokens of connect enabled tasks.
type csiHook struct {
	alloc    *structs.Allocation
	taskName string
	siClient consul.ServiceIdentityAPI
	logger   log.Logger

	lock     sync.Mutex
	firstRun bool
}

func newCSIHook(c csiHookConfig) *csiHook {
	return &csiHook{
		alloc:    c.alloc,
		taskName: c.task.Name,
		siClient: c.siClient,
		logger:   c.logger.Named(csiHookName),
		firstRun: true,
	}
}

func (h *csiHook) Name() string {
	return csiHookName
}

func (h *csiHook) Prestart(
	ctx context.Context,
	req *interfaces.TaskPrestartRequest,
	resp *interfaces.TaskPrestartResponse) error {

	h.lock.Lock()
	defer h.lock.Unlock()

	if h.earlyExit() {
		return nil
	}

	// try to recover token from disk
	token, err := h.recoverToken(req.TaskDir.SecretsDir)
	if err != nil {
		return err
	}

	// actually use h.siClient to get tokens and things
	fmt.Printf("## csiHook.Prestart, alloc: %s, task: %s\n", h.alloc.ID, h.taskName)
	fmt.Printf("##   recovered token: %s\n", token)

	if token == "" {
		// we need to derive a token
	}

	return nil
}

// earlyExit returns true if the Prestart hook has already been executed during
// the lifecycle of this task runner.
//
// assumes h is locked
func (h *csiHook) earlyExit() bool {
	if h.firstRun {
		h.firstRun = false
		return false
	}
	return true
}

// recoverToken returns the token saved to disk in the secrets directory for the
// task if it exists, or the empty string if the file does not exist. an error
// is returned only for some other (e.g. disk IO) error.
func (h *csiHook) recoverToken(dir string) (string, error) {
	tokenPath := filepath.Join(dir, csiTokenFile)
	token, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", errors.Wrap(err, "failed to recover SI token")
		}
		return "", nil // token file does not exist yet
	}
	return string(token), nil
}

func (h *csiHook) deriveSIToken(ctx context.Context) (string, error) {

	getCh := make(chan string)

	// keep trying to get the token in the background
	go func(ch chan<- string) {
		for {
			tokens, err := h.siClient.DeriveSITokens(h.alloc, []string{h.taskName})
			if err != nil {
				// log an error?
				continue // backoff
			}
			ch <- tokens[h.taskName]
			return
		}
	}(getCh)

	// wait until we get a token, or we get a signal to quit
	for {
		select {
		case token := <-getCh:
			return token, nil
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}
