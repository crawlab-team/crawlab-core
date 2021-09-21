package plugin

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/constants"
	errors2 "github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/delegate"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/process"
	"github.com/crawlab-team/crawlab-core/sys_exec"
	"github.com/crawlab-team/crawlab-core/utils"
	vcs "github.com/crawlab-team/crawlab-vcs"
	"github.com/crawlab-team/go-trace"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/dig"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Service struct {
	// settings variables
	fsPathBase      string
	monitorInterval time.Duration
	pluginBaseUrl   string

	// dependencies
	modelSvc service.ModelService

	// internals
	daemonMap sync.Map
	stopped   bool
}

func (svc *Service) Init() (err error) {
	return nil
}

func (svc *Service) Start() {
	svc.initPlugins()
}

func (svc *Service) Wait() {
	utils.DefaultWait()
}

func (svc *Service) Stop() {
	// do nothing
}

func (svc *Service) SetFsPathBase(path string) {
	svc.fsPathBase = path
}

func (svc *Service) SetMonitorInterval(interval time.Duration) {
	svc.monitorInterval = interval
}

func (svc *Service) SetPluginBaseUrl(baseUrl string) {
	svc.pluginBaseUrl = baseUrl
}

func (svc *Service) InstallPlugin(id primitive.ObjectID) (err error) {
	// plugin
	p, err := svc.modelSvc.GetPluginById(id)
	if err != nil {
		return err
	}

	// save status (installing)
	p.Status = constants.PluginStatusInstalling
	p.Error = ""
	_ = delegate.NewModelDelegate(p).Save()

	// install
	switch p.InstallType {
	case constants.PluginInstallTypeName:
		p, err = svc.installName(p)
	case constants.PluginInstallTypeGit:
		p, err = svc.installGit(p)
	case constants.PluginInstallTypeLocal:
		p, err = svc.installLocal(p)
	default:
		err = errors2.ErrorPluginNotImplemented
	}
	if err != nil {
		p.Status = constants.PluginStatusInstallError
		p.Error = err.Error()
		return trace.TraceError(err)
	}

	// save status (stopped)
	p.Status = constants.PluginStatusStopped
	p.Error = ""
	_ = delegate.NewModelDelegate(p).Save()

	return nil
}

func (svc *Service) UninstallPlugin(id primitive.ObjectID) (err error) {
	// TODO: implement
	panic("implement me")
}

func (svc *Service) RunPlugin(id primitive.ObjectID) (err error) {
	// plugin
	p, err := svc.modelSvc.GetPluginById(id)
	if err != nil {
		return err
	}

	// save pid
	p.Pid = 0
	p.Status = constants.PluginStatusRunning
	_ = delegate.NewModelDelegate(p).Save()

	// fs service
	fsSvc, err := NewPluginFsService(id)
	if err != nil {
		return err
	}

	// sync to workspace
	if err := fsSvc.GetFsService().SyncToWorkspace(); err != nil {
		return err
	}

	// process daemon
	d := process.NewProcessDaemon(svc.getNewCmdFn(p, fsSvc))

	// add to daemon map
	svc.addDaemon(id, d)

	// run (async)
	go func() {
		// start (async)
		go func() {
			if err := d.Start(); err != nil {
				svc.handleCmdError(p, err)
				return
			}
		}()

		// listening to signal from daemon
		stopped := false
		for {
			if stopped {
				break
			}
			ch := d.GetCh()
			sig := <-ch
			switch sig {
			case process.SignalStart:
				// save pid
				p.Pid = d.GetCmd().Process.Pid
				_ = delegate.NewModelDelegate(p).Save()
			case process.SignalStopped, process.SignalReachedMaxErrors:
				// break for loop
				stopped = true
			default:
				continue
			}
		}

		// stopped
		p.Status = constants.PluginStatusStopped
		p.Pid = 0
		p.Error = ""
		_ = delegate.NewModelDelegate(p).Save()
	}()

	return nil
}

func (svc *Service) StopPlugin(id primitive.ObjectID) (err error) {
	var d interfaces.ProcessDaemon
	if d = svc.getDaemon(id); d == nil {
		return trace.TraceError(errors2.ErrorPluginNotExists)
	}
	d.Stop()
	svc.deleteDaemon(id)
	return nil
}

func (svc *Service) installName(p interfaces.Plugin) (_p *models.Plugin, err error) {
	p.SetInstallUrl(fmt.Sprintf("%s%s", svc.pluginBaseUrl, p.GetName()))
	return svc.installGit(p)
}

func (svc *Service) installGit(p interfaces.Plugin) (_p *models.Plugin, err error) {
	log.Infof("git installing %s", p.GetInstallUrl())

	// git clone to temporary directory
	pluginPath := filepath.Join(os.TempDir(), uuid.New().String())
	gitClient, err := vcs.CloneGitRepo(pluginPath, p.GetInstallUrl())
	if err != nil {
		return nil, err
	}

	// sync to fs
	fsSvc, err := GetPluginFsService(p.GetId())
	if err != nil {
		return nil, err
	}
	if err := fsSvc.GetFsService().GetFs().SyncLocalToRemote(pluginPath, fsSvc.GetFsPath()); err != nil {
		return nil, err
	}

	// plugin.json
	_p, err = svc.getPluginFromJson(pluginPath)
	if err != nil {
		return nil, err
	}

	// fill plugin data and save to db
	_p.SetId(p.GetId())
	if err := delegate.NewModelDelegate(_p).Save(); err != nil {
		return nil, err
	}

	// dispose temporary directory
	if err := gitClient.Dispose(); err != nil {
		return nil, err
	}

	log.Infof("git installed %s", p.GetInstallUrl())
	return _p, nil
}

func (svc *Service) installLocal(p interfaces.Plugin) (_p *models.Plugin, err error) {
	log.Infof("local installing %s", p.GetInstallUrl())

	// plugin path
	var pluginPath string
	if strings.HasPrefix(p.GetInstallUrl(), "file://") {
		pluginPath = strings.Replace(p.GetInstallUrl(), "file://", "", 1)
		if !utils.Exists(pluginPath) {
			return nil, trace.TraceError(errors2.ErrorPluginPathNotExists)
		}
	}

	// plugin.json
	_p, err = svc.getPluginFromJson(pluginPath)
	if err != nil {
		return nil, err
	}

	// sync to fs
	fsSvc, err := GetPluginFsService(p.GetId())
	if err != nil {
		return nil, err
	}
	if err := fsSvc.GetFsService().GetFs().SyncLocalToRemote(pluginPath, fsSvc.GetFsPath()); err != nil {
		return nil, err
	}

	// fill plugin data and save to db
	_p.SetId(p.GetId())
	_p.SetInstallUrl(p.GetInstallUrl())
	if err := delegate.NewModelDelegate(_p).Save(); err != nil {
		return nil, err
	}

	log.Infof("local installed %s", p.GetInstallUrl())

	return _p, nil
}

func (svc *Service) getDaemon(id primitive.ObjectID) (d interfaces.ProcessDaemon) {
	res, ok := svc.daemonMap.Load(id)
	if !ok {
		return nil
	}
	d, ok = res.(interfaces.ProcessDaemon)
	if !ok {
		return nil
	}
	return d
}

func (svc *Service) addDaemon(id primitive.ObjectID, d interfaces.ProcessDaemon) {
	svc.daemonMap.Store(id, d)
}

func (svc *Service) deleteDaemon(id primitive.ObjectID) {
	svc.daemonMap.Delete(id)
}

func (svc *Service) handleCmdError(p *models.Plugin, err error) {
	trace.PrintError(err)
	p.Status = constants.PluginStatusError
	p.Pid = 0
	p.Error = err.Error()
	_ = delegate.NewModelDelegate(p).Save()
	svc.deleteDaemon(p.Id)
}

func (svc *Service) getNewCmdFn(p *models.Plugin, fsSvc interfaces.PluginFsService) func() (cmd *exec.Cmd) {
	return func() (cmd *exec.Cmd) {
		// command
		cmd = sys_exec.BuildCmd(p.Cmd)

		// working directory
		cmd.Dir = fsSvc.GetWorkspacePath()

		// inherit system envs
		for _, env := range os.Environ() {
			cmd.Env = append(cmd.Env, env)
		}

		// bind all viper keys to envs
		for _, key := range viper.AllKeys() {
			value := viper.Get(key)
			_, ok := value.(string)
			if !ok {
				continue
			}
			envName := fmt.Sprintf("%s_%s", "CRAWLAB", strings.ReplaceAll(strings.ToUpper(key), ".", "_"))
			envValue := viper.GetString(key)
			env := fmt.Sprintf("%s=%s", envName, envValue)
			cmd.Env = append(cmd.Env, env)
		}

		// logging
		stdout, _ := cmd.StdoutPipe()
		stderr, _ := cmd.StderrPipe()
		scannerStdout := bufio.NewScanner(stdout)
		scannerStderr := bufio.NewScanner(stderr)
		go func() {
			for scannerStdout.Scan() {
				line := fmt.Sprintf("[Plugin-%s] %s\n", p.GetName(), scannerStdout.Text())
				_, _ = os.Stdout.WriteString(line)
			}
		}()
		go func() {
			for scannerStderr.Scan() {
				line := fmt.Sprintf("[PLUGIN-%s] %s\n", p.GetName(), scannerStderr.Text())
				_, _ = os.Stderr.WriteString(line)
			}
		}()

		return cmd
	}
}

func (svc *Service) initPlugins() {
	plugins, err := svc.modelSvc.GetPluginList(nil, nil)
	if err != nil {
		trace.PrintError(err)
		return
	}
	for _, p := range plugins {
		if p.Restart {
			if err := svc.RunPlugin(p.Id); err != nil {
				trace.PrintError(err)
				continue
			}
		} else {
			if p.Status == constants.PluginStatusRunning {
				p.Error = errors2.ErrorPluginMissingProcess.Error()
				p.Pid = 0
				p.Status = constants.PluginStatusError
				_ = delegate.NewModelDelegate(&p).Save()
			}
		}
	}
}

func (svc *Service) getPluginFromJson(pluginPath string) (p *models.Plugin, err error) {
	pluginJsonPath := filepath.Join(pluginPath, "plugin.json")
	if !utils.Exists(pluginJsonPath) {
		return nil, trace.TraceError(errors2.ErrorPluginPluginJsonNotExists)
	}
	pluginJsonData, err := ioutil.ReadFile(pluginJsonPath)
	if err != nil {
		return nil, trace.TraceError(err)
	}
	var _p models.Plugin
	if err := json.Unmarshal(pluginJsonData, &_p); err != nil {
		return nil, trace.TraceError(err)
	}
	return &_p, nil
}

func NewPluginService(opts ...Option) (svc2 interfaces.PluginService, err error) {
	// service
	svc := &Service{
		fsPathBase:      DefaultPluginFsPathBase,
		monitorInterval: 15 * time.Second,
		pluginBaseUrl:   "https://github.com/crawlab-team/plugin-",
		daemonMap:       sync.Map{},
	}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.NewService); err != nil {
		return nil, err
	}
	if err := c.Invoke(func(
		modelSvc service.ModelService,
	) {
		svc.modelSvc = modelSvc
	}); err != nil {
		return nil, err
	}

	// initialize
	if err := svc.Init(); err != nil {
		return nil, err
	}

	return svc, nil
}

var store = sync.Map{}

func GetPluginService(path string, opts ...Option) (svc interfaces.PluginService, err error) {
	// return if service exists
	res, ok := store.Load(path)
	if ok {
		svc, ok = res.(interfaces.PluginService)
		if ok {
			return svc, nil
		}
	}

	// plugin name base url
	pluginBaseUrl := viper.GetString("plugin.baseUrl")
	if pluginBaseUrl != "" {
		opts = append(opts, WithPluginBaseUrl(pluginBaseUrl))
	}

	// service
	svc, err = NewPluginService(opts...)
	if err != nil {
		return nil, err
	}

	// save to cache
	store.Store(path, svc)

	return svc, nil
}

func ProvideGetPluginService(path string, opts ...Option) func() (svr interfaces.PluginService, err error) {
	return func() (svr interfaces.PluginService, err error) {
		return GetPluginService(path, opts...)
	}
}
