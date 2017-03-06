// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/ContinuumLLC/platform-asset-plugin/src/model (interfaces: HandlerDependencies,AssetService,AssetServiceFactory,AssetServiceDependencies,AssetDal,AssetDalFactory,AssetDalDependencies,ConfigServiceFactory,ConfigService,ConfigDal,ConfigDalDependencies,ConfigDalFactory,ConfigServiceDependencies)

package mock

import (
	io "io"

	asset "github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	model "github.com/ContinuumLLC/platform-asset-plugin/src/model"
	clar "github.com/ContinuumLLC/platform-common-lib/src/clar"
	env "github.com/ContinuumLLC/platform-common-lib/src/env"
	json "github.com/ContinuumLLC/platform-common-lib/src/json"
	protocol "github.com/ContinuumLLC/platform-common-lib/src/plugin/protocol"
	procParser "github.com/ContinuumLLC/platform-common-lib/src/procParser"
	gomock "github.com/golang/mock/gomock"
)

// Mock of HandlerDependencies interface
type MockHandlerDependencies struct {
	ctrl     *gomock.Controller
	recorder *_MockHandlerDependenciesRecorder
}

// Recorder for MockHandlerDependencies (not exported)
type _MockHandlerDependenciesRecorder struct {
	mock *MockHandlerDependencies
}

func NewMockHandlerDependencies(ctrl *gomock.Controller) *MockHandlerDependencies {
	mock := &MockHandlerDependencies{ctrl: ctrl}
	mock.recorder = &_MockHandlerDependenciesRecorder{mock}
	return mock
}

func (_m *MockHandlerDependencies) EXPECT() *_MockHandlerDependenciesRecorder {
	return _m.recorder
}

func (_m *MockHandlerDependencies) GetAssetDal(_param0 model.AssetDalDependencies) model.AssetDal {
	ret := _m.ctrl.Call(_m, "GetAssetDal", _param0)
	ret0, _ := ret[0].(model.AssetDal)
	return ret0
}

func (_mr *_MockHandlerDependenciesRecorder) GetAssetDal(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAssetDal", arg0)
}

func (_m *MockHandlerDependencies) GetAssetService(_param0 model.AssetServiceDependencies) model.AssetService {
	ret := _m.ctrl.Call(_m, "GetAssetService", _param0)
	ret0, _ := ret[0].(model.AssetService)
	return ret0
}

func (_mr *_MockHandlerDependenciesRecorder) GetAssetService(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAssetService", arg0)
}

func (_m *MockHandlerDependencies) GetConfigDal(_param0 model.ConfigDalDependencies) model.ConfigDal {
	ret := _m.ctrl.Call(_m, "GetConfigDal", _param0)
	ret0, _ := ret[0].(model.ConfigDal)
	return ret0
}

func (_mr *_MockHandlerDependenciesRecorder) GetConfigDal(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetConfigDal", arg0)
}

func (_m *MockHandlerDependencies) GetConfigService(_param0 model.ConfigServiceDependencies) model.ConfigService {
	ret := _m.ctrl.Call(_m, "GetConfigService", _param0)
	ret0, _ := ret[0].(model.ConfigService)
	return ret0
}

func (_mr *_MockHandlerDependenciesRecorder) GetConfigService(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetConfigService", arg0)
}

func (_m *MockHandlerDependencies) GetDeserializerJSON() json.DeserializerJSON {
	ret := _m.ctrl.Call(_m, "GetDeserializerJSON")
	ret0, _ := ret[0].(json.DeserializerJSON)
	return ret0
}

func (_mr *_MockHandlerDependenciesRecorder) GetDeserializerJSON() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetDeserializerJSON")
}

func (_m *MockHandlerDependencies) GetEnv() env.Env {
	ret := _m.ctrl.Call(_m, "GetEnv")
	ret0, _ := ret[0].(env.Env)
	return ret0
}

func (_mr *_MockHandlerDependenciesRecorder) GetEnv() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetEnv")
}

func (_m *MockHandlerDependencies) GetHandler(_param0 model.HandlerDependencies, _param1 *model.AssetPluginConfig) model.Handler {
	ret := _m.ctrl.Call(_m, "GetHandler", _param0, _param1)
	ret0, _ := ret[0].(model.Handler)
	return ret0
}

func (_mr *_MockHandlerDependenciesRecorder) GetHandler(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetHandler", arg0, arg1)
}

func (_m *MockHandlerDependencies) GetParser() procParser.Parser {
	ret := _m.ctrl.Call(_m, "GetParser")
	ret0, _ := ret[0].(procParser.Parser)
	return ret0
}

func (_mr *_MockHandlerDependenciesRecorder) GetParser() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetParser")
}

func (_m *MockHandlerDependencies) GetReader() io.Reader {
	ret := _m.ctrl.Call(_m, "GetReader")
	ret0, _ := ret[0].(io.Reader)
	return ret0
}

func (_mr *_MockHandlerDependenciesRecorder) GetReader() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetReader")
}

func (_m *MockHandlerDependencies) GetSerializerJSON() json.SerializerJSON {
	ret := _m.ctrl.Call(_m, "GetSerializerJSON")
	ret0, _ := ret[0].(json.SerializerJSON)
	return ret0
}

func (_mr *_MockHandlerDependenciesRecorder) GetSerializerJSON() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetSerializerJSON")
}

func (_m *MockHandlerDependencies) GetServer(_param0 io.Reader, _param1 io.Writer) protocol.Server {
	ret := _m.ctrl.Call(_m, "GetServer", _param0, _param1)
	ret0, _ := ret[0].(protocol.Server)
	return ret0
}

func (_mr *_MockHandlerDependenciesRecorder) GetServer(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetServer", arg0, arg1)
}

func (_m *MockHandlerDependencies) GetServiceInit() clar.ServiceInit {
	ret := _m.ctrl.Call(_m, "GetServiceInit")
	ret0, _ := ret[0].(clar.ServiceInit)
	return ret0
}

func (_mr *_MockHandlerDependenciesRecorder) GetServiceInit() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetServiceInit")
}

func (_m *MockHandlerDependencies) GetWriter() io.Writer {
	ret := _m.ctrl.Call(_m, "GetWriter")
	ret0, _ := ret[0].(io.Writer)
	return ret0
}

func (_mr *_MockHandlerDependenciesRecorder) GetWriter() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetWriter")
}

// Mock of AssetService interface
type MockAssetService struct {
	ctrl     *gomock.Controller
	recorder *_MockAssetServiceRecorder
}

// Recorder for MockAssetService (not exported)
type _MockAssetServiceRecorder struct {
	mock *MockAssetService
}

func NewMockAssetService(ctrl *gomock.Controller) *MockAssetService {
	mock := &MockAssetService{ctrl: ctrl}
	mock.recorder = &_MockAssetServiceRecorder{mock}
	return mock
}

func (_m *MockAssetService) EXPECT() *_MockAssetServiceRecorder {
	return _m.recorder
}

func (_m *MockAssetService) Process() (*asset.AssetCollection, error) {
	ret := _m.ctrl.Call(_m, "Process")
	ret0, _ := ret[0].(*asset.AssetCollection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockAssetServiceRecorder) Process() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Process")
}

// Mock of AssetServiceFactory interface
type MockAssetServiceFactory struct {
	ctrl     *gomock.Controller
	recorder *_MockAssetServiceFactoryRecorder
}

// Recorder for MockAssetServiceFactory (not exported)
type _MockAssetServiceFactoryRecorder struct {
	mock *MockAssetServiceFactory
}

func NewMockAssetServiceFactory(ctrl *gomock.Controller) *MockAssetServiceFactory {
	mock := &MockAssetServiceFactory{ctrl: ctrl}
	mock.recorder = &_MockAssetServiceFactoryRecorder{mock}
	return mock
}

func (_m *MockAssetServiceFactory) EXPECT() *_MockAssetServiceFactoryRecorder {
	return _m.recorder
}

func (_m *MockAssetServiceFactory) GetAssetService(_param0 model.AssetServiceDependencies) model.AssetService {
	ret := _m.ctrl.Call(_m, "GetAssetService", _param0)
	ret0, _ := ret[0].(model.AssetService)
	return ret0
}

func (_mr *_MockAssetServiceFactoryRecorder) GetAssetService(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAssetService", arg0)
}

// Mock of AssetServiceDependencies interface
type MockAssetServiceDependencies struct {
	ctrl     *gomock.Controller
	recorder *_MockAssetServiceDependenciesRecorder
}

// Recorder for MockAssetServiceDependencies (not exported)
type _MockAssetServiceDependenciesRecorder struct {
	mock *MockAssetServiceDependencies
}

func NewMockAssetServiceDependencies(ctrl *gomock.Controller) *MockAssetServiceDependencies {
	mock := &MockAssetServiceDependencies{ctrl: ctrl}
	mock.recorder = &_MockAssetServiceDependenciesRecorder{mock}
	return mock
}

func (_m *MockAssetServiceDependencies) EXPECT() *_MockAssetServiceDependenciesRecorder {
	return _m.recorder
}

func (_m *MockAssetServiceDependencies) GetAssetDal(_param0 model.AssetDalDependencies) model.AssetDal {
	ret := _m.ctrl.Call(_m, "GetAssetDal", _param0)
	ret0, _ := ret[0].(model.AssetDal)
	return ret0
}

func (_mr *_MockAssetServiceDependenciesRecorder) GetAssetDal(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAssetDal", arg0)
}

func (_m *MockAssetServiceDependencies) GetEnv() env.Env {
	ret := _m.ctrl.Call(_m, "GetEnv")
	ret0, _ := ret[0].(env.Env)
	return ret0
}

func (_mr *_MockAssetServiceDependenciesRecorder) GetEnv() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetEnv")
}

func (_m *MockAssetServiceDependencies) GetParser() procParser.Parser {
	ret := _m.ctrl.Call(_m, "GetParser")
	ret0, _ := ret[0].(procParser.Parser)
	return ret0
}

func (_mr *_MockAssetServiceDependenciesRecorder) GetParser() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetParser")
}

// Mock of AssetDal interface
type MockAssetDal struct {
	ctrl     *gomock.Controller
	recorder *_MockAssetDalRecorder
}

// Recorder for MockAssetDal (not exported)
type _MockAssetDalRecorder struct {
	mock *MockAssetDal
}

func NewMockAssetDal(ctrl *gomock.Controller) *MockAssetDal {
	mock := &MockAssetDal{ctrl: ctrl}
	mock.recorder = &_MockAssetDalRecorder{mock}
	return mock
}

func (_m *MockAssetDal) EXPECT() *_MockAssetDalRecorder {
	return _m.recorder
}

func (_m *MockAssetDal) GetAssetData() (*asset.AssetCollection, error) {
	ret := _m.ctrl.Call(_m, "GetAssetData")
	ret0, _ := ret[0].(*asset.AssetCollection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockAssetDalRecorder) GetAssetData() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAssetData")
}

func (_m *MockAssetDal) GetOSInfo() (*asset.AssetOs, error) {
	ret := _m.ctrl.Call(_m, "GetOSInfo")
	ret0, _ := ret[0].(*asset.AssetOs)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockAssetDalRecorder) GetOSInfo() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetOSInfo")
}

func (_m *MockAssetDal) GetSystemInfo() (*asset.AssetSystem, error) {
	ret := _m.ctrl.Call(_m, "GetSystemInfo")
	ret0, _ := ret[0].(*asset.AssetSystem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockAssetDalRecorder) GetSystemInfo() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetSystemInfo")
}

// Mock of AssetDalFactory interface
type MockAssetDalFactory struct {
	ctrl     *gomock.Controller
	recorder *_MockAssetDalFactoryRecorder
}

// Recorder for MockAssetDalFactory (not exported)
type _MockAssetDalFactoryRecorder struct {
	mock *MockAssetDalFactory
}

func NewMockAssetDalFactory(ctrl *gomock.Controller) *MockAssetDalFactory {
	mock := &MockAssetDalFactory{ctrl: ctrl}
	mock.recorder = &_MockAssetDalFactoryRecorder{mock}
	return mock
}

func (_m *MockAssetDalFactory) EXPECT() *_MockAssetDalFactoryRecorder {
	return _m.recorder
}

func (_m *MockAssetDalFactory) GetAssetDal(_param0 model.AssetDalDependencies) model.AssetDal {
	ret := _m.ctrl.Call(_m, "GetAssetDal", _param0)
	ret0, _ := ret[0].(model.AssetDal)
	return ret0
}

func (_mr *_MockAssetDalFactoryRecorder) GetAssetDal(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAssetDal", arg0)
}

// Mock of AssetDalDependencies interface
type MockAssetDalDependencies struct {
	ctrl     *gomock.Controller
	recorder *_MockAssetDalDependenciesRecorder
}

// Recorder for MockAssetDalDependencies (not exported)
type _MockAssetDalDependenciesRecorder struct {
	mock *MockAssetDalDependencies
}

func NewMockAssetDalDependencies(ctrl *gomock.Controller) *MockAssetDalDependencies {
	mock := &MockAssetDalDependencies{ctrl: ctrl}
	mock.recorder = &_MockAssetDalDependenciesRecorder{mock}
	return mock
}

func (_m *MockAssetDalDependencies) EXPECT() *_MockAssetDalDependenciesRecorder {
	return _m.recorder
}

func (_m *MockAssetDalDependencies) GetEnv() env.Env {
	ret := _m.ctrl.Call(_m, "GetEnv")
	ret0, _ := ret[0].(env.Env)
	return ret0
}

func (_mr *_MockAssetDalDependenciesRecorder) GetEnv() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetEnv")
}

func (_m *MockAssetDalDependencies) GetParser() procParser.Parser {
	ret := _m.ctrl.Call(_m, "GetParser")
	ret0, _ := ret[0].(procParser.Parser)
	return ret0
}

func (_mr *_MockAssetDalDependenciesRecorder) GetParser() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetParser")
}

// Mock of ConfigServiceFactory interface
type MockConfigServiceFactory struct {
	ctrl     *gomock.Controller
	recorder *_MockConfigServiceFactoryRecorder
}

// Recorder for MockConfigServiceFactory (not exported)
type _MockConfigServiceFactoryRecorder struct {
	mock *MockConfigServiceFactory
}

func NewMockConfigServiceFactory(ctrl *gomock.Controller) *MockConfigServiceFactory {
	mock := &MockConfigServiceFactory{ctrl: ctrl}
	mock.recorder = &_MockConfigServiceFactoryRecorder{mock}
	return mock
}

func (_m *MockConfigServiceFactory) EXPECT() *_MockConfigServiceFactoryRecorder {
	return _m.recorder
}

func (_m *MockConfigServiceFactory) GetConfigService(_param0 model.ConfigServiceDependencies) model.ConfigService {
	ret := _m.ctrl.Call(_m, "GetConfigService", _param0)
	ret0, _ := ret[0].(model.ConfigService)
	return ret0
}

func (_mr *_MockConfigServiceFactoryRecorder) GetConfigService(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetConfigService", arg0)
}

// Mock of ConfigService interface
type MockConfigService struct {
	ctrl     *gomock.Controller
	recorder *_MockConfigServiceRecorder
}

// Recorder for MockConfigService (not exported)
type _MockConfigServiceRecorder struct {
	mock *MockConfigService
}

func NewMockConfigService(ctrl *gomock.Controller) *MockConfigService {
	mock := &MockConfigService{ctrl: ctrl}
	mock.recorder = &_MockConfigServiceRecorder{mock}
	return mock
}

func (_m *MockConfigService) EXPECT() *_MockConfigServiceRecorder {
	return _m.recorder
}

func (_m *MockConfigService) GetAssetPluginConfMap() (map[string]interface{}, error) {
	ret := _m.ctrl.Call(_m, "GetAssetPluginConfMap")
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockConfigServiceRecorder) GetAssetPluginConfMap() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAssetPluginConfMap")
}

func (_m *MockConfigService) GetAssetPluginConfig() (*model.AssetPluginConfig, error) {
	ret := _m.ctrl.Call(_m, "GetAssetPluginConfig")
	ret0, _ := ret[0].(*model.AssetPluginConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockConfigServiceRecorder) GetAssetPluginConfig() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAssetPluginConfig")
}

func (_m *MockConfigService) SetAssetPluginMap(_param0 map[string]interface{}) error {
	ret := _m.ctrl.Call(_m, "SetAssetPluginMap", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockConfigServiceRecorder) SetAssetPluginMap(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetAssetPluginMap", arg0)
}

// Mock of ConfigDal interface
type MockConfigDal struct {
	ctrl     *gomock.Controller
	recorder *_MockConfigDalRecorder
}

// Recorder for MockConfigDal (not exported)
type _MockConfigDalRecorder struct {
	mock *MockConfigDal
}

func NewMockConfigDal(ctrl *gomock.Controller) *MockConfigDal {
	mock := &MockConfigDal{ctrl: ctrl}
	mock.recorder = &_MockConfigDalRecorder{mock}
	return mock
}

func (_m *MockConfigDal) EXPECT() *_MockConfigDalRecorder {
	return _m.recorder
}

func (_m *MockConfigDal) GetAssetPluginConf() (*model.AssetPluginConfig, error) {
	ret := _m.ctrl.Call(_m, "GetAssetPluginConf")
	ret0, _ := ret[0].(*model.AssetPluginConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockConfigDalRecorder) GetAssetPluginConf() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAssetPluginConf")
}

func (_m *MockConfigDal) GetAssetPluginConfMap() (map[string]interface{}, error) {
	ret := _m.ctrl.Call(_m, "GetAssetPluginConfMap")
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockConfigDalRecorder) GetAssetPluginConfMap() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAssetPluginConfMap")
}

func (_m *MockConfigDal) SetAssetPluginMap(_param0 map[string]interface{}) error {
	ret := _m.ctrl.Call(_m, "SetAssetPluginMap", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockConfigDalRecorder) SetAssetPluginMap(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetAssetPluginMap", arg0)
}

// Mock of ConfigDalDependencies interface
type MockConfigDalDependencies struct {
	ctrl     *gomock.Controller
	recorder *_MockConfigDalDependenciesRecorder
}

// Recorder for MockConfigDalDependencies (not exported)
type _MockConfigDalDependenciesRecorder struct {
	mock *MockConfigDalDependencies
}

func NewMockConfigDalDependencies(ctrl *gomock.Controller) *MockConfigDalDependencies {
	mock := &MockConfigDalDependencies{ctrl: ctrl}
	mock.recorder = &_MockConfigDalDependenciesRecorder{mock}
	return mock
}

func (_m *MockConfigDalDependencies) EXPECT() *_MockConfigDalDependenciesRecorder {
	return _m.recorder
}

func (_m *MockConfigDalDependencies) GetDeserializerJSON() json.DeserializerJSON {
	ret := _m.ctrl.Call(_m, "GetDeserializerJSON")
	ret0, _ := ret[0].(json.DeserializerJSON)
	return ret0
}

func (_mr *_MockConfigDalDependenciesRecorder) GetDeserializerJSON() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetDeserializerJSON")
}

func (_m *MockConfigDalDependencies) GetEnv() env.Env {
	ret := _m.ctrl.Call(_m, "GetEnv")
	ret0, _ := ret[0].(env.Env)
	return ret0
}

func (_mr *_MockConfigDalDependenciesRecorder) GetEnv() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetEnv")
}

func (_m *MockConfigDalDependencies) GetSerializerJSON() json.SerializerJSON {
	ret := _m.ctrl.Call(_m, "GetSerializerJSON")
	ret0, _ := ret[0].(json.SerializerJSON)
	return ret0
}

func (_mr *_MockConfigDalDependenciesRecorder) GetSerializerJSON() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetSerializerJSON")
}

func (_m *MockConfigDalDependencies) GetServiceInit() clar.ServiceInit {
	ret := _m.ctrl.Call(_m, "GetServiceInit")
	ret0, _ := ret[0].(clar.ServiceInit)
	return ret0
}

func (_mr *_MockConfigDalDependenciesRecorder) GetServiceInit() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetServiceInit")
}

// Mock of ConfigDalFactory interface
type MockConfigDalFactory struct {
	ctrl     *gomock.Controller
	recorder *_MockConfigDalFactoryRecorder
}

// Recorder for MockConfigDalFactory (not exported)
type _MockConfigDalFactoryRecorder struct {
	mock *MockConfigDalFactory
}

func NewMockConfigDalFactory(ctrl *gomock.Controller) *MockConfigDalFactory {
	mock := &MockConfigDalFactory{ctrl: ctrl}
	mock.recorder = &_MockConfigDalFactoryRecorder{mock}
	return mock
}

func (_m *MockConfigDalFactory) EXPECT() *_MockConfigDalFactoryRecorder {
	return _m.recorder
}

func (_m *MockConfigDalFactory) GetConfigDal(_param0 model.ConfigDalDependencies) model.ConfigDal {
	ret := _m.ctrl.Call(_m, "GetConfigDal", _param0)
	ret0, _ := ret[0].(model.ConfigDal)
	return ret0
}

func (_mr *_MockConfigDalFactoryRecorder) GetConfigDal(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetConfigDal", arg0)
}

// Mock of ConfigServiceDependencies interface
type MockConfigServiceDependencies struct {
	ctrl     *gomock.Controller
	recorder *_MockConfigServiceDependenciesRecorder
}

// Recorder for MockConfigServiceDependencies (not exported)
type _MockConfigServiceDependenciesRecorder struct {
	mock *MockConfigServiceDependencies
}

func NewMockConfigServiceDependencies(ctrl *gomock.Controller) *MockConfigServiceDependencies {
	mock := &MockConfigServiceDependencies{ctrl: ctrl}
	mock.recorder = &_MockConfigServiceDependenciesRecorder{mock}
	return mock
}

func (_m *MockConfigServiceDependencies) EXPECT() *_MockConfigServiceDependenciesRecorder {
	return _m.recorder
}

func (_m *MockConfigServiceDependencies) GetConfigDal(_param0 model.ConfigDalDependencies) model.ConfigDal {
	ret := _m.ctrl.Call(_m, "GetConfigDal", _param0)
	ret0, _ := ret[0].(model.ConfigDal)
	return ret0
}

func (_mr *_MockConfigServiceDependenciesRecorder) GetConfigDal(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetConfigDal", arg0)
}

func (_m *MockConfigServiceDependencies) GetDeserializerJSON() json.DeserializerJSON {
	ret := _m.ctrl.Call(_m, "GetDeserializerJSON")
	ret0, _ := ret[0].(json.DeserializerJSON)
	return ret0
}

func (_mr *_MockConfigServiceDependenciesRecorder) GetDeserializerJSON() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetDeserializerJSON")
}

func (_m *MockConfigServiceDependencies) GetEnv() env.Env {
	ret := _m.ctrl.Call(_m, "GetEnv")
	ret0, _ := ret[0].(env.Env)
	return ret0
}

func (_mr *_MockConfigServiceDependenciesRecorder) GetEnv() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetEnv")
}

func (_m *MockConfigServiceDependencies) GetSerializerJSON() json.SerializerJSON {
	ret := _m.ctrl.Call(_m, "GetSerializerJSON")
	ret0, _ := ret[0].(json.SerializerJSON)
	return ret0
}

func (_mr *_MockConfigServiceDependenciesRecorder) GetSerializerJSON() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetSerializerJSON")
}

func (_m *MockConfigServiceDependencies) GetServiceInit() clar.ServiceInit {
	ret := _m.ctrl.Call(_m, "GetServiceInit")
	ret0, _ := ret[0].(clar.ServiceInit)
	return ret0
}

func (_mr *_MockConfigServiceDependenciesRecorder) GetServiceInit() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetServiceInit")
}
