package plugin

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/cenkalti/backoff/v4"
	"github.com/crawlab-team/crawlab-core/constants"
	errors2 "github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models/client"
	"github.com/crawlab-team/crawlab-core/models/delegate"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/node/config"
	"github.com/crawlab-team/crawlab-core/process"
	"github.com/crawlab-team/crawlab-core/sys_exec"
	"github.com/crawlab-team/crawlab-core/utils"
	vcs "github.com/crawlab-team/crawlab-vcs"
	"github.com/crawlab-team/go-trace"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	cfgSvc                     interfaces.NodeConfigService
	modelSvc                   service.ModelService
	clientModelSvc             interfaces.GrpcClientModelService
	clientModelNodeSvc         interfaces.GrpcClientModelNodeService
	clientModelPluginSvc       interfaces.GrpcClientModelPluginService
	clientModelPluginStatusSvc interfaces.GrpcClientModelPluginStatusService

	// internals
	daemonMap sync.Map
	stopped   bool
	n         *models.Node // current node
}

func (svc *Service) Init() (err error) {
	return nil
}

func (svc *Service) Start() {
	// get current node
	if err := svc.getCurrentNode(); err != nil {
		panic(err)
	}

	// get global settings
	if err := svc.getPluginBaseUrl(); err != nil {
		panic(err)
	}

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
	p, err := svc.getPluginById(id)
	if err != nil {
		return err
	}

	// plugin status
	ps, err := svc.getPluginStatus(id)
	if err != nil {
		return err
	}

	// save status (installing)
	ps.Status = constants.PluginStatusInstalling
	ps.Error = ""
	_ = svc.savePluginStatus(ps)

	// get plugin base url
	if err := svc.getPluginBaseUrl(); err != nil {
		return err
	}

	// install
	switch p.InstallType {
	case constants.PluginInstallTypeName:
		err = svc.installName(p)
	case constants.PluginInstallTypeGit:
		err = svc.installGit(p)
	case constants.PluginInstallTypeLocal:
		err = svc.installLocal(p)
	default:
		err = errors2.ErrorPluginNotImplemented
	}
	if err != nil {
		ps.Status = constants.PluginStatusInstallError
		ps.Error = err.Error()
		return trace.TraceError(err)
	}

	// wait
	for i := 0; i < 10; i++ {
		query := bson.M{
			"plugin_id": id,
			"node_id":   svc.n.Id,
		}
		ps, err = svc.modelSvc.GetPluginStatus(query, nil)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if p.AutoStart {
		// start plugin
		_ = svc.StartPlugin(p.Id)
	} else {
		// save status (stopped)
		ps.Status = constants.PluginStatusStopped
		ps.Error = ""
		_ = svc.savePluginStatus(ps)
	}

	return nil
}

func (svc *Service) UninstallPlugin(id primitive.ObjectID) (err error) {
	// plugin
	_, err = svc.getPluginById(id)
	if err != nil {
		return err
	}

	// plugin status
	ps, err := svc.getPluginStatus(id)
	if err != nil {
		return err
	}

	// stop
	if ps.Status == constants.PluginStatusRunning {
		if err := svc.StopPlugin(id); err != nil {
			return err
		}
	}

	// delete fs (master)
	if svc.cfgSvc.IsMaster() {
		// fs service
		fsSvc, err := NewPluginFsService(id)
		if err != nil {
			return err
		}

		// delete fs
		fsPath := fsSvc.GetFsPath()
		if err := fsSvc.GetFsService().Delete(fsPath); err != nil {
			return err
		}
	}

	return nil
}

func (svc *Service) StartPlugin(id primitive.ObjectID) (err error) {
	// plugin
	p, err := svc.getPluginById(id)
	if err != nil {
		return err
	}

	// plugin status
	ps, err := svc.getPluginStatus(id)
	if err != nil {
		return err
	}

	// save pid
	ps.Pid = 0
	ps.Status = constants.PluginStatusRunning
	_ = svc.savePluginStatus(ps)

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
				svc.handleCmdError(p, ps, err)
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
				ps.Pid = d.GetCmd().Process.Pid
				_ = svc.savePluginStatus(ps)
			case process.SignalStopped, process.SignalReachedMaxErrors:
				// break for loop
				stopped = true
			default:
				continue
			}
		}

		// stopped
		ps.Status = constants.PluginStatusStopped
		ps.Pid = 0
		ps.Error = ""
		if _, err := svc.getPluginStatus(p.GetId()); err != nil {
			if err.Error() != mongo.ErrNoDocuments.Error() {
				trace.PrintError(err)
			}
			return
		}
		_ = svc.savePluginStatus(ps)
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

func (svc *Service) installName(p interfaces.Plugin) (err error) {
	p.SetInstallUrl(fmt.Sprintf("%s%s", svc.pluginBaseUrl, p.GetName()))
	return svc.installGit(p)
}

func (svc *Service) installGit(p interfaces.Plugin) (err error) {
	log.Infof("git installing %s", p.GetInstallUrl())

	// git clone to temporary directory
	pluginPath := filepath.Join(os.TempDir(), uuid.New().String())
	gitClient, err := vcs.CloneGitRepo(pluginPath, p.GetInstallUrl())
	if err != nil {
		return err
	}

	// sync to fs
	fsSvc, err := GetPluginFsService(p.GetId())
	if err != nil {
		return err
	}
	if err := fsSvc.GetFsService().GetFs().SyncLocalToRemote(pluginPath, fsSvc.GetFsPath()); err != nil {
		return err
	}

	// plugin.json
	_p, err := svc.getPluginFromJson(pluginPath)
	if err != nil {
		return err
	}

	// fill plugin data and save to db
	if svc.cfgSvc.IsMaster() {
		_p.SetId(p.GetId())
		if err := svc.savePlugin(_p); err != nil {
			return err
		}
	}

	// dispose temporary directory
	if err := gitClient.Dispose(); err != nil {
		return err
	}

	log.Infof("git installed %s", p.GetInstallUrl())
	return nil
}

func (svc *Service) installLocal(p interfaces.Plugin) (err error) {
	log.Infof("local installing %s", p.GetInstallUrl())

	// plugin path
	pluginPath := p.GetInstallUrl()
	if strings.HasPrefix(p.GetInstallUrl(), "file://") {
		pluginPath = strings.Replace(p.GetInstallUrl(), "file://", "", 1)
		if !utils.Exists(pluginPath) {
			return trace.TraceError(errors2.ErrorPluginPathNotExists)
		}
	}

	// plugin.json
	_p, err := svc.getPluginFromJson(pluginPath)
	if err != nil {
		return err
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
	if svc.cfgSvc.IsMaster() {
		_p.SetId(p.GetId())
		_p.SetInstallUrl(p.GetInstallUrl())
		_p.SetInstallType(p.GetInstallType())
		if err := svc.savePlugin(_p); err != nil {
			return err
		}
	}

	log.Infof("local installed %s", p.GetInstallUrl())

	return nil
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

func (svc *Service) handleCmdError(p *models.Plugin, ps *models.PluginStatus, err error) {
	trace.PrintError(err)
	ps.Status = constants.PluginStatusError
	ps.Pid = 0
	ps.Error = err.Error()
	_ = svc.savePluginStatus(ps)
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
	// reset plugin status
	psList, err := svc.modelSvc.GetPluginStatusList(nil, nil)
	for _, ps := range psList {
		if ps.Status == constants.PluginStatusRunning {
			ps.Status = constants.PluginStatusError
			ps.Error = errors2.ErrorPluginMissingProcess.Error()
			ps.Pid = 0
			_ = svc.savePluginStatus(&ps)
		}
	}

	// plugins
	plugins, err := svc.modelSvc.GetPluginList(nil, nil)
	if err != nil {
		trace.PrintError(err)
		return
	}

	// restart plugins that need restart
	for _, p := range plugins {
		if p.AutoStart {
			if err := svc.StartPlugin(p.Id); err != nil {
				trace.PrintError(err)
				continue
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

func (svc *Service) getPluginById(id primitive.ObjectID) (p *models.Plugin, err error) {
	if svc.cfgSvc.IsMaster() {
		p, err = svc.modelSvc.GetPluginById(id)
		if err != nil {
			return nil, err
		}
		return p, nil
	} else {
		_p, err := svc.clientModelPluginSvc.GetPluginById(id)
		if err != nil {
			return nil, err
		}
		p, ok := _p.(*models.Plugin)
		if !ok {
			return nil, trace.TraceError(errors2.ErrorPluginInvalidType)
		}
		return p, nil
	}
}

func (svc *Service) getPluginStatus(pluginId primitive.ObjectID) (ps *models.PluginStatus, err error) {
	if svc.cfgSvc.IsMaster() {
		ps, err = svc.modelSvc.GetPluginStatus(bson.M{
			"plugin_id": pluginId,
			"node_id":   svc.n.Id,
		}, nil)
		if err != nil {
			// add if not exists
			if strings.Contains(err.Error(), mongo.ErrNoDocuments.Error()) {
				return svc.addPluginStatus(pluginId, svc.n.Id)
			}

			// error
			return nil, err
		}
		return ps, nil
	} else {
		_ps, err := svc.clientModelPluginStatusSvc.GetPluginStatus(bson.M{
			"plugin_id": pluginId,
			"node_id":   svc.n.Id,
		}, nil)
		if err != nil {
			// add if not exists
			if strings.Contains(err.Error(), mongo.ErrNoDocuments.Error()) {
				return svc.addPluginStatus(pluginId, svc.n.Id)
			}
		}
		ps, ok := _ps.(*models.PluginStatus)
		if !ok {
			return nil, trace.TraceError(errors2.ErrorPluginInvalidType)
		}
		return ps, nil
	}
}

func (svc *Service) savePlugin(p *models.Plugin) (err error) {
	if svc.cfgSvc.IsMaster() {
		return delegate.NewModelDelegate(p).Save()
	} else {
		return client.NewModelDelegate(p).Save()
	}
}

func (svc *Service) savePluginStatus(ps interfaces.PluginStatus) (err error) {
	if svc.cfgSvc.IsMaster() {
		return delegate.NewModelDelegate(ps).Save()
	} else {
		return client.NewModelDelegate(ps).Save()
	}
}

func (svc *Service) addPluginStatus(pid primitive.ObjectID, nid primitive.ObjectID) (ps *models.PluginStatus, err error) {
	ps = &models.PluginStatus{
		PluginId: pid,
		NodeId:   nid,
	}
	if svc.cfgSvc.IsMaster() {
		if err := delegate.NewModelDelegate(ps).Add(); err != nil {
			return nil, err
		}
	} else {
		psD := client.NewModelDelegate(ps)
		if err := psD.Add(); err != nil {
			return nil, err
		}
		ps = psD.GetModel().(*models.PluginStatus)
	}
	return ps, nil
}

func (svc *Service) getCurrentNode() (err error) {
	return backoff.RetryNotify(svc._getNode, backoff.NewExponentialBackOff(), utils.BackoffErrorNotify("plugin service get node"))
}

func (svc *Service) _getNode() (err error) {
	if svc.cfgSvc.IsMaster() {
		svc.n, err = svc.modelSvc.GetNodeByKey(svc.cfgSvc.GetNodeKey(), nil)
		if err != nil {
			return err
		}
	} else {
		_n, err := svc.clientModelNodeSvc.GetNodeByKey(svc.cfgSvc.GetNodeKey())
		if err != nil {
			return err
		}
		n, ok := _n.(*models.Node)
		if !ok {
			return trace.TraceError(errors2.ErrorPluginInvalidType)
		}
		svc.n = n
	}
	return nil
}

func (svc *Service) getPluginBaseUrl() (err error) {
	return backoff.RetryNotify(svc._getPluginBaseUrl, backoff.NewExponentialBackOff(), utils.BackoffErrorNotify("plugin service get global settings"))
}

func (svc *Service) _getPluginBaseUrl() (err error) {
	if svc.cfgSvc.IsMaster() {
		s, err := svc.modelSvc.GetSettingByKey(constants.SettingPlugin, nil)
		if err != nil {
			if err.Error() == mongo.ErrNoDocuments.Error() {
				value := bson.M{}
				value[constants.SettingPluginBaseUrl] = constants.DefaultSettingPluginBaseUrl
				s := &models.Setting{
					Key:   constants.SettingPlugin,
					Value: value,
				}
				if err := delegate.NewModelDelegate(s).Add(); err != nil {
					return err
				}
				svc.pluginBaseUrl = constants.DefaultSettingPluginBaseUrl
				return nil
			}
			return err
		}
		res, ok := s.Value[constants.SettingPluginBaseUrl]
		if ok {
			svc.pluginBaseUrl, _ = res.(string)
		}
		return nil
	} else {
		var settingModelSvc interfaces.GrpcClientModelBaseService
		if err := backoff.Retry(func() error {
			if svc.clientModelSvc == nil {
				return errors.New("clientModelSvc is nil")
			}
			settingModelSvc, err = svc.clientModelSvc.NewBaseServiceDelegate(interfaces.ModelIdSetting)
			if err != nil {
				return err
			}
			return nil
		}, backoff.NewConstantBackOff(1*time.Second)); err != nil {
			return trace.TraceError(err)
		}
		_s, err := settingModelSvc.Get(bson.M{"key": constants.SettingPluginBaseUrl}, nil)
		if err != nil {
			return err
		}
		s, ok := _s.(*models.Setting)
		if !ok {
			return trace.TraceError(errors2.ErrorPluginInvalidType)
		}
		res, ok := s.Value[constants.SettingPluginBaseUrl]
		if ok {
			svc.pluginBaseUrl, _ = res.(string)
		}
		return nil
	}
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
	if err := c.Provide(config.NewNodeConfigService); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(service.GetService); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(client.NewNodeServiceDelegate); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(client.NewPluginServiceDelegate); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(client.NewPluginStatusServiceDelegate); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Invoke(func(
		cfgSvc interfaces.NodeConfigService,
		modelSvc service.ModelService,
		clientModelNodeSvc interfaces.GrpcClientModelNodeService,
		clientModelPluginSvc interfaces.GrpcClientModelPluginService,
		clientModelPluginStatusSvc interfaces.GrpcClientModelPluginStatusService,
	) {
		svc.cfgSvc = cfgSvc
		svc.modelSvc = modelSvc
		svc.clientModelNodeSvc = clientModelNodeSvc
		svc.clientModelPluginSvc = clientModelPluginSvc
		svc.clientModelPluginStatusSvc = clientModelPluginStatusSvc
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
