# NOTE: THIS KS WILL BE THE GENERIC TEMPLATE FOR ALL CENTOS SERVER ROLES AFTER 
#       CLEANUP. THIS IS WIP.

# ===================================================================================
#                                     COMMON Setup
# ===================================================================================
# System authorization information
auth --useshadow --passalgo=sha512

# Reboot after installation
reboot

# SELinux configuration
selinux --disabled

# System services
services --enabled="systemd-networkd,systemd-resolved,sshd"

# Use cmdline install that is not interactive 
text

# Keyboard layouts
keyboard --vckeymap=us --xlayouts=''

# System language
lang en_US.UTF-8

# Use network installation media
nfs --server=nfs101 --dir=/persist2/root/centos7/CentOS-7-x86_64-DVD-1611.iso

# Root password
rootpw --plaintext arastra

# Do not configure the X Window System
skipx

# System timezone
timezone America/Los_Angeles --ntpservers=ntp.aristanetworks.com

# ===================================================================================
#                                     DISK Setup
# ===================================================================================
# Only use these disks to setup
# TODO NON-GENERIC - must be scriptified 
ignoredisk --only-use=sda,sdb,sdc

# System bootloader configuration
# TODO NON-GENERIC - must be scriptified 
bootloader --location=mbr --boot-drive=sda --timeout=5
bootloader --location=mbr --boot-drive=sdb --timeout=5
bootloader --location=mbr --boot-drive=sdc --timeout=5

# Clean up previously created partitions
# TODO NON-GENERIC - must be scriptified
clearpart --drives=sda,sdb,sdc --all --initlabel

# Create partitions on each disk
# TODO NON-GENERIC - must be scriptified
part biosboot --size=2 --ondisk=sda --fstype=biosboot 
part biosboot --size=2 --ondisk=sdb --fstype=biosboot 
part biosboot --size=2 --ondisk=sdc --fstype=biosboot

part raid.11 --asprimary --size=20000 --ondisk=sda --fstype=ext4
part raid.21 --asprimary --size=20000 --ondisk=sdb --fstype=ext4
part raid.31 --asprimary --size=20000 --ondisk=sdc --fstype=ext4

part raid.12 --asprimary --size=1 --grow --ondisk=sda --fstype=ext4
part raid.22 --asprimary --size=1 --grow --ondisk=sdb --fstype=ext4
part raid.32 --asprimary --size=1 --grow --ondisk=sdc --fstype=ext4

raid / --level=1 --fstype=ext4 --device=1 --label=/ raid.11 raid.21 raid.31
raid /persist --level=0 --fstype=ext4 --device=2 --label=/persist raid.12 raid.22 raid.32

# ===================================================================================
#                                  PACKAGES Section 
# ===================================================================================
%packages --nobase 
@core --nodefaults
-selinux-policy
-selinux-policy-targeted
-aic94xx-firmware*
-alsa-*
-biosdevname
-btrfs-progs*
-dracut-network
-iprutils
-ivtv*
-iwl*firmware
-libertas*
-NetworkManager*
-plymouth*
-postfix
-uboot-tools
-microcode_ctl
parted
nfs-utils
wget
lsof
kexec-tools
ntp
net-tools
%end

# ===================================================================================
#                                   ADD-ON section
# ===================================================================================
%addon com_redhat_kdump --disable --reserve-mb='128'
%end

# ===================================================================================
#                                   POST Section
# ===================================================================================
%post --log=/root/ks-post.log

echo "Following commands are being run from \%post section in KS"

# Temporarily mount NFS to get needed files
# TODO: This is a temporary workaround, should stop doing this once a better
# solution is found.
echo "Mounting nfs101:/persist2/milkyway to get needed files..."
mkdir -p /tmp/root
mount nfs101:/persist2/milkyway /tmp/root

# Minimum pkg requirements
# NOTE: Need this epel-release as a requirement for all other pkgs
echo "Installing packages..."
yum install -y epel-release
yum install -y systemd-networkd systemd-resolved
yum install -y cpufrequtils sysstat ipmitool pdsh vim-minimal
yum install -y htop dstat strace tcpdump xfsprogs tar lsof
yum install -y pciutils yum-utils mysql git
yum install -y ansible smartmontools bind-utils
yum install -y yum-cron

echo "Setting up default arastra and ansible accounts..."
# Arastra Default Account Setup
groupadd -g 10000 arastra
useradd -u 10000 -g 10000 -c "Arista Networks anonymous account" -G wheel,root arastra
mkdir -p ~arastra/.ssh
cp /tmp/root/arastra-ssh/* ~arastra/.ssh/
chown arastra:arastra -R ~arastra/.ssh
chmod 700 ~arastra/.ssh
chmod 600 ~arastra/.ssh/*

# Ansible Generic Account Setup
groupadd -g 20000 ansible
useradd -u 20000 -g 20000 -c "Ansible account" -G wheel,root ansible
mkdir -p ~ansible/.ssh
cp /tmp/root/ansible-ssh/authorized_keys ~ansible/.ssh/
chown ansible:ansible -R ~ansible/.ssh
chmod 700 ~ansible/.ssh
chmod 644 ~ansible/.ssh/authorized_keys

echo "Setting Arora18Release and ToolsV2 repos..."
# setup Arora18Release repo
cat <<EOF > /etc/yum.repos.d/Arora.repo
[Arora18Release]
name=AroraRelease
baseurl=http://dist/Abuild/Arora18.release/x86_64_18/latest/RPMS
enabled=1
gpgcheck=0
metadata_expire=2h
EOF

# Setup ToolsV2 repo
cat <<EOF > /etc/yum.repos.d/ToolsV2.repo
[ToolsV2]
name=ToolsV2 \$basearch
baseurl=http://tools/ToolsV2/repo/\$basearch/RPMS/
enabled=1
gpgcheck=0
metadata_expire=2h
exclude=scylla*
EOF

echo "Configuring network bonding..."
# Remove ifcfg scripts since we use systemd
rm /etc/sysconfig/network-scripts/ifcfg-eno1
rm /etc/sysconfig/network-scripts/ifcfg-eno2

export ENO1_MAC_ADDR=`ip link show dev eno1 | awk '/ether/ {print($2)}'`
/tmp/root/SetupSystemdBond.py $ENO1_MAC_ADDR

yum install -y lldpad
systemctl enable lldpad
lldptool -i eno1 -T -V sysName enableTx=yes
lldptool -i eno2 -T -V sysName enableTx=yes

# workaround for systemd-networkd bug (fixed in master)
# for setting hostname. This workaround is needed
# since CentOS 7.3 onwards, /etc/hostname is always
# localhost.localdomain and hostname depends on the
# transient hostname setup by the network owner - either
# NetworkManager or systemd-networkd. In CentOS 7.3,
# systemd-networkd is version v219 which lacks the fix
# refer to https://github.com/martinpitt/systemd/commit/e8c0de91271331ddbae872de63d0a267d4f71e12
# for more details
cat <<EOF > /etc/polkit-1/rules.d/51-set-hostname.rules
polkit.addRule(function(action, subject) {
      if (action.id == "org.freedesktop.hostname1.set-hostname" && subject.user == "systemd-network") {
      return polkit.Result.YES;
      }
      });
EOF

systemctl enable systemd-networkd
systemctl enable systemd-resolved

# Setup NTP Server Config
echo 'server ntp.aristanetworks.com' >> /etc/ntp.conf
systemctl enable ntpd

echo "Setting up serial console..."
# setup serial console
# must retain content added by anaconda
SERIAL_PORT=1
FILE="/etc/default/grub"
sed -i "s/GRUB_TERMINAL_OUTPUT=.*$/GRUB_TERMINAL_OUTPUT=\"serial console\"/" ${FILE}
sed -i "s/GRUB_SERIAL_COMMAND=.*/GRUB_SERIAL_COMMAND=\"serial --speed=115200 --unit=${SERIAL_PORT} --word=8 --parity=no --stop=1\"/" ${FILE}
# anconda seems to pick only one console from pxelinux cfg 
CONSOLE_TTY0=`cat ${FILE} | awk '/console=tty0/ {found=1;} END { if (!found) print "console=tty0"; }'`
CONSOLE_TTYSX=`cat ${FILE} | awk "/console=ttyS${SERIAL_PORT}/ {found=1;} END { if (!found) print \"console=ttyS${SERIAL_PORT},115200n8\"; }"`
# ensure tty0 is before ttySX
[[ ! -z ${CONSOLE_TTY0} ]] && sed -i "s/GRUB_CMDLINE_LINUX=\"\(.*\) console\(.*\)\"/GRUB_CMDLINE_LINUX=\"\1 ${CONSOLE_TTY0} console\2\"/" ${FILE}
[[ ! -z ${CONSOLE_TTYSX} ]] && sed -i "s/GRUB_CMDLINE_LINUX=\"\(.*\)\"/GRUB_CMDLINE_LINUX=\"\1 ${CONSOLE_TTYSX}\"/" ${FILE}

cat ${FILE}
/usr/sbin/grub2-mkconfig -o /boot/grub2/grub.cfg

cat <<EOF >> /etc/sudoers
Defaults !env_reset, !requiretty
%wheel ALL=(ALL) NOPASSWD: ALL
EOF

# Disable native audit/analysis tools
tuned-adm off
systemctl disable tuned

%end
