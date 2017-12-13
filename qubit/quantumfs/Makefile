version := $(shell git describe || echo "dev-`git rev-parse HEAD`")
RPM_VERSION := $(shell echo "$(version)" | sed -e "s/^v//" -e "s/-/_/g")

all: rpmQuantumfsDevConf rpmQuantumfsProdConf rpmQuantumfsCiabConf

define DESCRIPTION_QUANTUMFS_PROD_CONF
This rpm contains following:
/etc/quantumfs.conf - QuantumFS config file to connect to the Production
                      Scylla Cluster
endef
export DESCRIPTION_QUANTUMFS_PROD_CONF

define DESCRIPTION_QUANTUMFS_DEV_CONF
This rpm contains following:
/etc/quantumfs.conf - QuantumFS config file to connect to the Dev
                      Scylla Cluster
endef
export DESCRIPTION_QUANTUMFS_DEV_CONF

define DESCRIPTION_QUANTUMFS_CIAB_CONF
This rpm contains following:
/etc/quantumfs.conf - QuantumFS config file to connect to the Dev
                      Scylla Cluster's CIAB keyspaces.
endef
export DESCRIPTION_QUANTUMFS_CIAB_CONF

YUM_FPM_PKGS := ruby-devel rubygems
setupFpm:
	@type fpm > /dev/null 2>&1 || (echo "Installing fpm" && \
	sudo yum install -y $(FPM_YUM_PKGS) && sudo gem install fpm)

RPM_BASENAME := QuantumFS-Conf
RPM_BASENAME_PROD := $(RPM_BASENAME)-Prod
RPM_BASENAME_DEV := $(RPM_BASENAME)-Dev
RPM_BASENAME_CIAB := $(RPM_BASENAME)-Ciab

RPM_FILES_TOOLSV2_X86_64 += $(RPM_BASENAME_PROD)-$(RPM_VERSION).x86_64.rpm
RPM_FILES_TOOLSV2_X86_64 += $(RPM_BASENAME_DEV)-$(RPM_VERSION).x86_64.rpm
RPM_FILES_TOOLSV2_X86_64 += $(RPM_BASENAME_CIAB)-$(RPM_VERSION).x86_64.rpm

version:
	@echo $(RPM_VERSION)

rpmQuantumfsProdConf: setupFpm
	fpm -f -s dir -t rpm -m 'ether-dev@arista.com' -n $(RPM_BASENAME_PROD) \
      		--license='Arista Proprietary' \
		--vendor='Arista Networks' \
		--url http://gerrit/ether \
		--rpm-summary="Production Scylla configuration for QuantumFS" \
		--description="$$DESCRIPTION_QUANTUMFS_PROD_CONF" \
		--version $(RPM_VERSION) \
		./prod-qubit1=/etc/quantumfs.conf

rpmQuantumfsDevConf: setupFpm
	fpm -f -s dir -t rpm -m 'ether-dev@arista.com' -n $(RPM_BASENAME_DEV) \
      		--license='Arista Proprietary' \
		--vendor='Arista Networks' \
		--url http://gerrit/ether \
		--rpm-summary="Dev Cluster configuration for QuantumFS" \
		--description="$$DESCRIPTION_QUANTUMFS_DEV_CONF" \
		--version $(RPM_VERSION) \
		./dev-dev_qubit1=/etc/quantumfs.conf

rpmQuantumfsCiabConf: setupFpm
	fpm -f -s dir -t rpm -m 'ether-dev@arista.com' -n $(RPM_BASENAME_CIAB) \
      		--license='Arista Proprietary' \
		--vendor='Arista Networks' \
		--url http://gerrit/ether \
		--rpm-summary="Ciab Cluster configuration for QuantumFS" \
		--description="$$DESCRIPTION_QUANTUMFS_CIAB_CONF" \
		--version $(RPM_VERSION) \
		./dev-dev_ninja0=/etc/quantumfs.conf

push-rpms: rpmQuantumfsCiabConf rpmQuantumfsProdConf rpmQuantumfsDevConf
	a4 scp $(RPM_FILES_TOOLSV2_X86_64) dist:/dist/release/ToolsV2/repo/x86_64/RPMS
	a4 ssh dist /usr/bin/createrepo --update /dist/release/ToolsV2/repo/x86_64/RPMS
	@echo
	@echo "If you're refreshing existing RPMs, then on machines which use this repo you should:"
	@echo "   sudo yum clean all"
	@echo "   sudo yum makecache"

clean:
	rm -rf *.rpm
