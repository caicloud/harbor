# Makefile for a harbor project
#	
# Targets:
#
# build: 	build harbor photon images
# clean:	clean adminserver, ui and jobservice harbor images

# common
SHELL := /bin/bash
BUILDPATH=$(CURDIR)
MAKEPATH=$(BUILDPATH)/make
MAKEDEVPATH=$(MAKEPATH)/dev
SRCPATH=./src
TOOLSPATH=$(CURDIR)/tools
SEDCMD=$(shell which sed)
WGET=$(shell which wget)

# docker parameters
DOCKERCMD=$(shell which docker)
DOCKERBUILD=$(DOCKERCMD) build --pull
DOCKERRMIMAGE=$(DOCKERCMD) rmi
DOCKERIMASES=$(DOCKERCMD) images

# binary 
ADMINSERVERSOURCECODE=$(SRCPATH)/adminserver
ADMINSERVERBINARYPATH=$(MAKEDEVPATH)/adminserver
ADMINSERVERBINARYNAME=harbor_adminserver
CORE_SOURCECODE=$(SRCPATH)/core
CORE_BINARYPATH=$(MAKEDEVPATH)/core
CORE_BINARYNAME=harbor_core
JOBSERVICESOURCECODE=$(SRCPATH)/jobservice
JOBSERVICEBINARYPATH=$(MAKEDEVPATH)/jobservice
JOBSERVICEBINARYNAME=harbor_jobservice

# photon dockerfile
DOCKERFILEPATH=$(MAKEPATH)/photon

DOCKERFILEPATH_ADMINSERVER=$(DOCKERFILEPATH)/adminserver
DOCKERFILENAME_ADMINSERVER=Dockerfile
DOCKERFILENAME_ADMINSERVER_ARM64=Dockerfile.arm64
DOCKERIMAGENAME_ADMINSERVER=goharbor/harbor-adminserver
DOCKERIMAGENAME_ADMINSERVER_ARM64=arm64/harbor-adminserver

DOCKERFILEPATH_PORTAL=$(DOCKERFILEPATH)/portal
DOCKERFILENAME_PORTAL=Dockerfile
DOCKERFILENAME_PORTAL_ARM64=Dockerfile.arm64
DOCKERIMAGENAME_PORTAL=goharbor/harbor-portal
DOCKERIMAGENAME_PORTAL_ARM64=arm64/harbor-portal

DOCKERFILEPATH_CORE=$(DOCKERFILEPATH)/core
DOCKERFILENAME_CORE=Dockerfile
DOCKERFILENAME_CORE_ARM64=Dockerfile.arm64
DOCKERIMAGENAME_CORE=goharbor/harbor-core
DOCKERIMAGENAME_CORE_ARM64=arm64/harbor-core

DOCKERFILEPATH_JOBSERVICE=$(DOCKERFILEPATH)/jobservice
DOCKERFILENAME_JOBSERVICE=Dockerfile
DOCKERFILENAME_JOBSERVICE_ARM64=Dockerfile.arm64
DOCKERIMAGENAME_JOBSERVICE=goharbor/harbor-jobservice
DOCKERIMAGENAME_JOBSERVICE_ARM64=arm64/harbor-jobservice

DOCKERFILEPATH_LOG=$(DOCKERFILEPATH)/log
DOCKERFILENAME_LOG=Dockerfile
DOCKERFILENAME_LOG_ARM64=Dockerfile.arm64
DOCKERIMAGENAME_LOG=goharbor/harbor-log
DOCKERIMAGENAME_LOG_ARM64=arm64/harbor-log

DOCKERFILEPATH_DB=$(DOCKERFILEPATH)/db
DOCKERFILENAME_DB=Dockerfile
DOCKERFILENAME_DB_ARM64=Dockerfile.arm64
DOCKERIMAGENAME_DB=goharbor/harbor-db
DOCKERIMAGENAME_DB_ARM64=arm64/harbor-db

DOCKERFILEPATH_POSTGRESQL=$(DOCKERFILEPATH)/postgresql
DOCKERFILENAME_POSTGRESQL=Dockerfile
DOCKERIMAGENAME_POSTGRESQL=goharbor/postgresql-photon

DOCKERFILEPATH_CLAIR=$(DOCKERFILEPATH)/clair
DOCKERFILENAME_CLAIR=Dockerfile
DOCKERIMAGENAME_CLAIR=goharbor/clair-photon

DOCKERFILEPATH_NGINX=$(DOCKERFILEPATH)/nginx
DOCKERFILENAME_NGINX=Dockerfile
DOCKERIMAGENAME_NGINX=goharbor/nginx-photon

DOCKERFILEPATH_REG=$(DOCKERFILEPATH)/registry
DOCKERFILENAME_REG=Dockerfile
DOCKERFILENAME_REG_ARM64=Dockerfile.arm64
DOCKERIMAGENAME_REG=goharbor/registry-photon
DOCKERIMAGENAME_REG_ARM64=arm64/registry-photon

DOCKERFILEPATH_REGISTRYCTL=$(DOCKERFILEPATH)/registryctl
DOCKERFILENAME_REGISTRYCTL=Dockerfile
DOCKERFILENAME_REGISTRYCTL_ARM64=Dockerfile.arm64
DOCKERIMAGENAME_REGISTRYCTL=goharbor/harbor-registryctl
DOCKERIMAGENAME_REGISTRYCTL_ARM64=arm64/harbor-registryctl

DOCKERFILEPATH_NOTARY=$(DOCKERFILEPATH)/notary
DOCKERFILENAME_NOTARYSIGNER=signer.Dockerfile
DOCKERIMAGENAME_NOTARYSIGNER=goharbor/notary-signer-photon
DOCKERFILENAME_NOTARYSERVER=server.Dockerfile
DOCKERIMAGENAME_NOTARYSERVER=goharbor/notary-server-photon

DOCKERFILEPATH_REDIS=$(DOCKERFILEPATH)/redis
DOCKERFILENAME_REDIS=Dockerfile
DOCKERFILENAME_REDIS_ARM64=Dockerfile.arm64
DOCKERIMAGENAME_REDIS=goharbor/redis-photon
DOCKERIMAGENAME_REDIS_ARM64=arm64/redis-photon

DOCKERFILEPATH_MIGRATOR=$(TOOLSPATH)/migration
DOCKERFILENAME_MIGRATOR=Dockerfile
DOCKERIMAGENAME_MIGRATOR=goharbor/harbor-migrator

# for chart server (chartmuseum)
DOCKERFILEPATH_CHART_SERVER=$(DOCKERFILEPATH)/chartserver
DOCKERFILENAME_CHART_SERVER=Dockerfile
CHART_SERVER_CODE_BASE=github.com/helm/chartmuseum
CHART_SERVER_MAIN_PATH=cmd/chartmuseum
CHART_SERVER_BIN_NAME=chartm

_build_db:
	@echo "building db container for photon..."
    ifeq ($(ARCH), "arm64")
	    @$(DOCKERBUILD) -f $(DOCKERFILEPATH_DB)/$(DOCKERFILENAME_DB) -t $(DOCKERIMAGENAME_DB_ARM64):$(VERSIONTAG) .
    else
	    @$(DOCKERBUILD) -f $(DOCKERFILEPATH_DB)/$(DOCKERFILENAME_DB) -t $(DOCKERIMAGENAME_DB):$(VERSIONTAG) .
    endif
	@echo "Done."

_build_adminserver:
	@echo "building adminserver container for photon..."
    ifeq ($(ARCH), "arm64")
	    @$(DOCKERBUILD) -f $(DOCKERFILEPATH_ADMINSERVER)/$(DOCKERFILENAME_ADMINSERVER) -t $(DOCKERIMAGENAME_ADMINSERVER_ARM64):$(VERSIONTAG) .
    else
	    @$(DOCKERBUILD) -f $(DOCKERFILEPATH_ADMINSERVER)/$(DOCKERFILENAME_ADMINSERVER) -t $(DOCKERIMAGENAME_ADMINSERVER):$(VERSIONTAG) .
    endif
	@echo "Done."

_build_portal:
	@echo "building portal container for photon..."
    ifeq ($(ARCH), "arm64")
	    $(DOCKERBUILD) -f $(DOCKERFILEPATH_PORTAL)/$(DOCKERFILENAME_PORTAL) -t $(DOCKERIMAGENAME_PORTAL_ARM64):$(VERSIONTAG) .
    else
	    $(DOCKERBUILD) -f $(DOCKERFILEPATH_PORTAL)/$(DOCKERFILENAME_PORTAL) -t $(DOCKERIMAGENAME_PORTAL):$(VERSIONTAG) .
    endif
	@echo "Done."

_build_core:	
	@echo "building core container for photon..."
    ifeq ($(ARCH), "arm64")
	    @$(DOCKERBUILD) -f $(DOCKERFILEPATH_CORE)/$(DOCKERFILENAME_CORE) -t $(DOCKERIMAGENAME_CORE_ARM64):$(VERSIONTAG) .
    else
	    @$(DOCKERBUILD) -f $(DOCKERFILEPATH_CORE)/$(DOCKERFILENAME_CORE) -t $(DOCKERIMAGENAME_CORE):$(VERSIONTAG) .
    endif
	@echo "Done."
	
_build_jobservice:		
	@echo "building jobservice container for photon..."
    ifeq ($(ARCH), "arm64")
	    @$(DOCKERBUILD) -f $(DOCKERFILEPATH_JOBSERVICE)/$(DOCKERFILENAME_JOBSERVICE) -t $(DOCKERIMAGENAME_JOBSERVICE_ARM64):$(VERSIONTAG) .
    else
	    @$(DOCKERBUILD) -f $(DOCKERFILEPATH_JOBSERVICE)/$(DOCKERFILENAME_JOBSERVICE) -t $(DOCKERIMAGENAME_JOBSERVICE):$(VERSIONTAG) .
    endif
	@echo "Done."

_build_log:	
	@echo "building log container for photon..."
    ifeq ($(ARCH), "arm64")
	    $(DOCKERBUILD) -f $(DOCKERFILEPATH_LOG)/$(DOCKERFILENAME_LOG) -t $(DOCKERIMAGENAME_LOG_ARM64):$(VERSIONTAG) .
    else
	    $(DOCKERBUILD) -f $(DOCKERFILEPATH_LOG)/$(DOCKERFILENAME_LOG) -t $(DOCKERIMAGENAME_LOG):$(VERSIONTAG) .
    endif
	@echo "Done."

_build_clair:
	@if [ "$(CLAIRFLAG)" = "true" ] ; then \
		if [ "$(BUILDBIN)" != "true" ] ; then \
			rm -rf $(DOCKERFILEPATH_CLAIR)/binary && mkdir -p $(DOCKERFILEPATH_CLAIR)/binary && \
			$(call _get_binary, https://storage.googleapis.com/harbor-builds/bin/clair/release2.0-$(CLAIRVERSION)/clair, $(DOCKERFILEPATH_CLAIR)/binary/clair); \
		else \
			cd $(DOCKERFILEPATH_CLAIR) && $(DOCKERFILEPATH_CLAIR)/builder $(CLAIRVERSION) && cd - ; \
		fi ; \
		echo "building clair container for photon..." ; \
		$(DOCKERBUILD) -f $(DOCKERFILEPATH_CLAIR)/$(DOCKERFILENAME_CLAIR) -t $(DOCKERIMAGENAME_CLAIR):$(CLAIRVERSION)-$(VERSIONTAG) . ; \
		rm -rf $(DOCKERFILEPATH_CLAIR)/binary; \
		echo "Done." ; \
	fi

_build_chart_server:
	@if [ "$(CHARTFLAG)" = "true" ] ; then \
		if [ "$(BUILDBIN)" != "true" ] ; then \
			rm -rf $(DOCKERFILEPATH_CHART_SERVER)/binary && mkdir -p $(DOCKERFILEPATH_CHART_SERVER)/binary && \
			$(call _get_binary, https://storage.googleapis.com/harbor-builds/bin/chartmuseum/release-$(CHARTMUSEUMVERSION)/chartm, $(DOCKERFILEPATH_CHART_SERVER)/binary/chartm); \
		else \
			cd $(DOCKERFILEPATH_CHART_SERVER) && $(DOCKERFILEPATH_CHART_SERVER)/builder $(GOBUILDIMAGE) $(CHART_SERVER_CODE_BASE) $(CHARTMUSEUMVERSION) $(CHART_SERVER_MAIN_PATH) $(CHART_SERVER_BIN_NAME) && cd - ; \
		fi ; \
		echo "building chartmuseum container for photon..." ; \
		$(DOCKERBUILD) -f $(DOCKERFILEPATH_CHART_SERVER)/$(DOCKERFILENAME_CHART_SERVER) -t $(DOCKERIMAGENAME_CHART_SERVER):$(CHARTMUSEUMVERSION)-$(VERSIONTAG) . ; \
		rm -rf $(DOCKERFILEPATH_CHART_SERVER)/binary; \
		echo "Done." ; \
	fi
		
_build_nginx:
	@echo "building nginx container for photon..."
	@$(DOCKERBUILD) -f $(DOCKERFILEPATH_NGINX)/$(DOCKERFILENAME_NGINX) -t $(DOCKERIMAGENAME_NGINX):$(NGINXVERSION) .
	@echo "Done."

_build_notary:
	@if [ "$(NOTARYFLAG)" = "true" ] ; then \
		if [ "$(BUILDBIN)" != "true" ] ; then \
			rm -rf $(DOCKERFILEPATH_NOTARY)/binary && mkdir -p $(DOCKERFILEPATH_NOTARY)/binary && \
			$(call _get_binary, https://storage.googleapis.com/harbor-builds/bin/notary/release-$(NOTARYVERSION)/binary-bundle.tgz, $(DOCKERFILEPATH_NOTARY)/binary-bundle.tgz); \
			cd $(DOCKERFILEPATH_NOTARY) && tar -zvxf binary-bundle.tgz && cd - ; \
		else \
			cd $(DOCKERFILEPATH_NOTARY) && $(DOCKERFILEPATH_NOTARY)/builder $(NOTARYVERSION) $(NOTARYMIGRATEVERSION) && cd - ; \
		fi ; \
		echo "building notary container for photon..."; \
		chmod 655 $(DOCKERFILEPATH_NOTARY)/binary/notary-signer && $(DOCKERBUILD) -f $(DOCKERFILEPATH_NOTARY)/$(DOCKERFILENAME_NOTARYSIGNER) -t $(DOCKERIMAGENAME_NOTARYSIGNER):$(NOTARYVERSION)-$(VERSIONTAG) . ; \
		chmod 655 $(DOCKERFILEPATH_NOTARY)/binary/notary-server && $(DOCKERBUILD) -f $(DOCKERFILEPATH_NOTARY)/$(DOCKERFILENAME_NOTARYSERVER) -t $(DOCKERIMAGENAME_NOTARYSERVER):$(NOTARYVERSION)-$(VERSIONTAG) . ; \
		rm -rf $(DOCKERFILEPATH_NOTARY)/binary; \
		echo "Done."; \
	fi
		
_build_registry:
	@if [ "$(BUILDBIN)" != "true" ] ; then \
		rm -rf $(DOCKERFILEPATH_REG)/binary && mkdir -p $(DOCKERFILEPATH_REG)/binary && \
		$(call _get_binary, https://storage.googleapis.com/harbor-builds/bin/registry, $(DOCKERFILEPATH_REG)/binary/registry); \
	else \
		cd $(DOCKERFILEPATH_REG) && $(DOCKERFILEPATH_REG)/builder $(REGISTRYVERSION) && cd - ; \
	fi
	@echo "building registry container for photon..."
    ifeq ($(ARCH), "arm64")
	    @chmod 655 $(DOCKERFILEPATH_REG)/binary/registry && $(DOCKERBUILD) -f $(DOCKERFILEPATH_REG)/$(DOCKERFILENAME_REG_ARM64) -t $(DOCKERIMAGENAME_REG):$(REGISTRYVERSION)-$(VERSIONTAG) .
    else
	    @chmod 655 $(DOCKERFILEPATH_REG)/binary/registry && $(DOCKERBUILD) -f $(DOCKERFILEPATH_REG)/$(DOCKERFILENAME_REG) -t $(DOCKERIMAGENAME_REG):$(REGISTRYVERSION)-$(VERSIONTAG) .
    endif
	@echo "Done."

_build_registryctl:		
	@echo "building registry controller for photon..."
    ifeq ($(ARCH), "arm64")
	    @$(DOCKERBUILD) -f $(DOCKERFILEPATH_REGISTRYCTL)/$(DOCKERFILENAME_REGISTRYCTL) -t $(DOCKERIMAGENAME_REGISTRYCTL_ARM64):$(VERSIONTAG) .
    else
	    @$(DOCKERBUILD) -f $(DOCKERFILEPATH_REGISTRYCTL)/$(DOCKERFILENAME_REGISTRYCTL) -t $(DOCKERIMAGENAME_REGISTRYCTL):$(VERSIONTAG) .
    endif
	@rm -rf $(DOCKERFILEPATH_REG)/binary
	@echo "Done."

_build_redis:
	@echo "building redis container for photon..."
    ifeq ($(ARCH), "arm64")
	    @$(DOCKERBUILD) -f $(DOCKERFILEPATH_REDIS)/$(DOCKERFILENAME_REDIS) -t $(DOCKERIMAGENAME_REDIS_ARM64):$(REDISVERSION) .
    else
	    @$(DOCKERBUILD) -f $(DOCKERFILEPATH_REDIS)/$(DOCKERFILENAME_REDIS) -t $(DOCKERIMAGENAME_REDIS):$(REDISVERSION) .
    endif
	@echo "Done."

_build_migrator:
	@echo "building db migrator container for photon..."
	@cd $(DOCKERFILEPATH_MIGRATOR) && $(DOCKERBUILD) -f $(DOCKERFILEPATH_MIGRATOR)/$(DOCKERFILENAME_MIGRATOR) -t $(DOCKERIMAGENAME_MIGRATOR):$(MIGRATORVERSION) .
	@echo "Done."

define _get_binary
	$(WGET) --timeout 30 --no-check-certificate $1 -O $2
endef

build: _build_db _build_adminserver _build_portal _build_core _build_jobservice _build_log _build_nginx _build_registry _build_registryctl _build_notary _build_clair _build_redis _build_migrator _build_chart_server

cleanimage:
	@echo "cleaning image for photon..."
	- $(DOCKERRMIMAGE) -f $(DOCKERIMAGENAME_ADMINSERVER):$(VERSIONTAG)
	- $(DOCKERRMIMAGE) -f $(DOCKERIMAGENAME_PORTAL):$(VERSIONTAG)
	- $(DOCKERRMIMAGE) -f $(DOCKERIMAGENAME_CORE):$(VERSIONTAG)
	- $(DOCKERRMIMAGE) -f $(DOCKERIMAGENAME_JOBSERVICE):$(VERSIONTAG)
	- $(DOCKERRMIMAGE) -f $(DOCKERIMAGENAME_LOG):$(VERSIONTAG)
	
.PHONY: clean
clean: cleanimage

