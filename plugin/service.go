package plugin

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/crawlab-team/crawlab-core/constants"
	errors2 "github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/delegate"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/process"
	"github.com/crawlab-team/crawlab-core/sys_exec"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/crawlab-team/go-trace"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/dig"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type Service struct {
	// settings variables
	fsPathBase string

	// dependencies
	modelSvc service.ModelService

	// internals
	daemonMap sync.Map
}

func (svc *Service) Init() (err error) {
	plugins, err := svc.modelSvc.GetPluginList(bson.M{"status": constants.PluginStatusRunning}, nil)
	if err != nil {
		return err
	}
	for _, p := range plugins {
		p.Error = errors2.ErrorPluginMissingProcess.Error()
		p.Pid = 0
		p.Status = constants.PluginStatusError
		if err := delegate.NewModelDelegate(&p).Save(); err != nil {
			return err
		}
	}
	return nil
}

func (svc *Service) SetFsPathBase(path string) {
	svc.fsPathBase = path
}

func (svc *Service) InstallPlugin(id primitive.ObjectID) (err error) {
	// plugin
	p, err := svc.modelSvc.GetPluginById(id)
	if err != nil {
		return err
	}

	// install url type
	installUrlType := svc.getInstallUrlType(p)

	// install
	switch installUrlType {
	case constants.PluginInstallUrlTypePluginName:
		return svc.installPluginName(p)
	case constants.PluginInstallUrlTypeGithub:
		return svc.installGithub(p)
	case constants.PluginInstallUrlTypeGitee:
		return svc.installGitee(p)
	case constants.PluginInstallUrlTypeFile:
		return svc.installFile(p)
	case constants.PluginInstallUrlTypeGeneralUrl:
		return svc.installGeneralUrl(p)
	default:
		return trace.TraceError(errors2.ErrorPluginNotImplemented)
	}
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

func (svc *Service) getInstallUrlType(p interfaces.Plugin) (installUrlType string) {
	if p.GetName() != "" {
		return constants.PluginInstallUrlTypePluginName
	}

	url := p.GetInstallUrl()
	if strings.Contains(url, "github.com") {
		return constants.PluginInstallUrlTypeGithub
	} else if strings.Contains(url, "gitee.com") {
		return constants.PluginInstallUrlTypeGitee
	} else if strings.HasPrefix(url, "file:///") {
		return constants.PluginInstallUrlTypeFile
	} else {
		return constants.PluginInstallUrlTypeGeneralUrl
	}
}

func (svc *Service) installPluginName(p interfaces.Plugin) (err error) {
	p.SetInstallUrl(fmt.Sprintf("https://github.com/crawlab-team/plugin-%s", p.GetName()))
	return svc.installGithub(p)
}

func (svc *Service) installGithub(p interfaces.Plugin) (err error) {
	// TODO: implement
	panic("not implemented")
}

func (svc *Service) installGitee(p interfaces.Plugin) (err error) {
	// TODO: implement
	panic("not implemented")
}

func (svc *Service) installFile(p interfaces.Plugin) (err error) {
	// plugin path
	pluginPath := strings.Replace(p.GetInstallUrl(), "file://", "", 1)
	if !utils.Exists(pluginPath) {
		return trace.TraceError(errors2.ErrorPluginPathNotExists)
	}

	// plugin.json
	pluginJsonPath := filepath.Join(pluginPath, "plugin.json")
	if !utils.Exists(pluginJsonPath) {
		return trace.TraceError(errors2.ErrorPluginPluginJsonNotExists)
	}
	pluginJsonData, err := ioutil.ReadFile(pluginJsonPath)
	if err != nil {
		return trace.TraceError(err)
	}
	var _p models.Plugin
	if err := json.Unmarshal(pluginJsonData, &_p); err != nil {
		return trace.TraceError(err)
	}

	// sync to fs
	fsSvc, err := GetPluginFsService(p.GetId())
	if err != nil {
		return err
	}
	if err := fsSvc.GetFsService().GetFs().SyncLocalToRemote(pluginPath, fsSvc.GetFsPath()); err != nil {
		return err
	}

	// fill plugin data and save to db
	_p.SetId(p.GetId())
	_p.SetInstallUrl(p.GetInstallUrl())
	if err := delegate.NewModelDelegate(&_p).Save(); err != nil {
		return err
	}

	return nil
}

func (svc *Service) installGeneralUrl(p interfaces.Plugin) (err error) {
	// TODO: implement
	panic("not implemented")
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

func (svc *Service) monitorCmd() {
	//for {
	//	plugins, err := svc.modelSvc.GetPluginList(nil, nil)
	//	if err != nil {
	//		trace.PrintError(err)
	//		continue
	//	}
	//	for _, p := range plugins {
	//		if p.Status == constants.PluginStatusRunning {
	//		} else {
	//		}
	//	}
	//}
}

func NewPluginService(opts ...Option) (svc2 interfaces.PluginService, err error) {
	// service
	svc := &Service{
		fsPathBase: DefaultPluginFsPathBase,
		daemonMap:  sync.Map{},
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
	res, ok := store.Load(path)
	if ok {
		svc, ok = res.(interfaces.PluginService)
		if ok {
			return svc, nil
		}
	}
	svc, err = NewPluginService(opts...)
	if err != nil {
		return nil, err
	}
	store.Store(path, svc)
	return svc, nil
}

func ProvideGetPluginService(path string, opts ...Option) func() (svr interfaces.PluginService, err error) {
	return func() (svr interfaces.PluginService, err error) {
		return GetPluginService(path, opts...)
	}
}
