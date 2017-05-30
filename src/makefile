# defining environment variables
TARGETNAME = platform-asset-plugin
ifeq ($(OS), Windows_NT)
	TARGET = $(TARGETNAME).exe
else
	TARGET = $(TARGETNAME)
endif

TARGETPATH=$(GOPATH)/bin/$(TARGET)

PLUGINPATH=$(GOPATH)/src/github.com/ContinuumLLC/$(TARGETNAME)/services
JSONCFGPATH=$(GOPATH)/src/github.com/ContinuumLLC/$(TARGETNAME)/src

JSONCFGNAME="ctm_asset_agent_plugin_cfg.json"


#Linux
LINUXPLUGINPACKPATH=$(PLUGINPATH)/continuum_plugin_linux
LINUXPLUGINBINPATH=$(LINUXPLUGINPACKPATH)/asset
LINUXPLUGINPACKNAME="continuum_asset_plugin_linux.zip"
LINUXPLUGINPACTARGET="continuum_plugin_linux"


#windows64
WINDOWSPLUGIN64PACKPATH=$(PLUGINPATH)/continuum_plugin_windows64
WINDOWSPLUGIN64BINPATH=$(WINDOWSPLUGIN64PACKPATH)/asset
WINDOWS64PLUGINPACKNAME="continuum_asset_plugin_windows64.zip"
WINDOWSPLUGIN64PACTARGET="continuum_plugin_windows64"

#windows32
WINDOWSPLUGIN32PACKPATH=$(PLUGINPATH)/continuum_plugin_windows32
WINDOWSPLUGIN32BINPATH=$(WINDOWSPLUGIN32PACKPATH)/asset
WINDOWS32PLUGINPACKNAME="continuum_asset_plugin_windows32.zip"
WINDOWSPLUGIN32PACTARGET="continuum_plugin_windows32"
 
#darwin
MACPLUGINPACKPATH=$(PLUGINPATH)/continuum_plugin_darwin
MACPLUGINBINPATH=$(MACPLUGINPACKPATH)/asset
MACPLUGINPACKNAME="continuum_asset_plugin_darwin.zip
MACPLUGINPACTARGET="continuum_plugin_darwin"


# global command
all: dependencies platform-asset-plugin test cover pack-service

linux: dependencies platform-asset-plugin-linux test cover pack-service-linux

windows: dependencies platform-asset-plugin-win64 platform-asset-plugin-win32 test cover pack-service-windows

mac: dependencies platform-asset-plugin-mac test cover pack-service-darwin

local: dependencies localbuild test code-quality
	
localbuild: 
	go build -o $(TARGETPATH) ./asset-plugin

pack-service: pack-service-linux pack-service-windows pack-service-darwin

platform-asset-plugin: platform-asset-plugin-win64 platform-asset-plugin-win32 platform-asset-plugin-linux #platform-asset-plugin-mac

platform-asset-plugin-linux:
	GOOS=linux GOARCH=amd64 go build -o $(GOPATH)/bin/linux_amd64/$(TARGETNAME) ./asset-plugin

platform-asset-plugin-win64:
	GOOS=windows GOARCH=amd64 go build -o $(GOPATH)/bin/windows_amd64/$(TARGETNAME).exe ./asset-plugin

platform-asset-plugin-win32:
	GOOS=windows GOARCH=386 go build -o $(GOPATH)/bin/windows_386/$(TARGETNAME).exe ./asset-plugin	

platform-asset-plugin-mac:
	GOOS=darwin GOARCH=amd64 go build -o $(GOPATH)/bin/darwin_amd64/$(TARGETNAME) ./asset-plugin

dependencies:
	glide install

test:
	go test -coverprofile ./services/cover.out -covermode=count ./services
	go test -coverprofile ./dal/cover.out -covermode=count ./dal
	go test -coverprofile ./msgl/cover.out -covermode=count ./msgl

cover:
	echo "mode: count" > cover-all.out
	tail -n +2 ./services/cover.out >> cover-all.out
	tail -n +2 ./dal/cover.out >> cover-all.out
	gocover-cobertura < cover-all.out > cover-cobertura.xml

code-quality:
	gometalinter --vendor --tests --skip=mock --disable=gotype --disable=errcheck --disable=gas --disable=dupl --deadline=1500s --checkstyle --sort=linter ./... > static-analysis.xml

mockgen:
	mockgen -package mock github.com/ContinuumLLC/platform-asset-plugin/src/model HandlerDependencies,AssetService,AssetServiceFactory,AssetServiceDependencies,AssetDal,AssetDalFactory,AssetDalDependencies,ConfigServiceFactory,ConfigService,ConfigDal,ConfigDalDependencies,ConfigDalFactory,ConfigServiceDependencies > model/mock/mocks.go

pack-service-linux:
	mkdir -p $(LINUXPLUGINBINPATH)
	cp $(GOPATH)/bin/linux_amd64/$(TARGETNAME) $(LINUXPLUGINBINPATH)
	cp $(JSONCFGPATH)/$(JSONCFGNAME) $(LINUXPLUGINBINPATH)
	cd $(PLUGINPATH)/$(LINUXPLUGINPACTARGET); rm -f $(LINUXPLUGINPACKNAME); zip -r ../$(LINUXPLUGINPACKNAME) *

pack-service-windows: pack-service-windows64 pack-service-windows32

pack-service-windows64:
	mkdir -p $(WINDOWSPLUGIN64BINPATH)
	cp $(GOPATH)/bin/windows_amd64/$(TARGETNAME).exe $(WINDOWSPLUGIN64BINPATH)
	cp $(JSONCFGPATH)/$(JSONCFGNAME) $(WINDOWSPLUGIN64BINPATH)
	cd $(PLUGINPATH)/$(WINDOWSPLUGIN64PACTARGET); rm -f $(WINDOWS64PLUGINPACKNAME); zip -r ../$(WINDOWS64PLUGINPACKNAME) *

pack-service-windows32:
	mkdir -p $(WINDOWSPLUGIN32BINPATH)
	cp $(GOPATH)/bin/windows_386/$(TARGETNAME).exe $(WINDOWSPLUGIN32BINPATH)
	cp $(JSONCFGPATH)/$(JSONCFGNAME) $(WINDOWSPLUGIN32BINPATH)
	cd $(PLUGINPATH)/$(WINDOWSPLUGIN32PACTARGET); rm -f $(WINDOWS32PLUGINPACKNAME); zip -r ../$(WINDOWS32PLUGINPACKNAME) *

pack-service-darwin:
	#mkdir -p $(MACPLUGINBINPATH)
	#cp $(GOPATH)/bin/darwin_amd64/$(TARGETNAME) $(MACPLUGINBINPATH)
	#cp $(JSONCFGPATH)/$(JSONCFGNAME) $(MACPLUGINBINPATH)
	#cd $(PLUGINPATH)/$(MACPLUGINPACTARGET); rm -f $(MACPLUGINPACKNAME); zip -r ../$(MACPLUGINPACKNAME) *

