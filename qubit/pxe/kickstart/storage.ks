#version=DEVEL
# System authorization information
auth --useshadow --passalgo=sha512
# Reboot after installation
reboot
# SELinux configuration
selinux --disabled
# System services
services --enabled="systemd-networkd,systemd-resolved"

# Use text mode install
text
# Use network installation media
url --url="http://pxe.aristanetworks.com/tftpboot/max/dist/"

ignoredisk --only-use=sda,sdb,sdc

# Keyboard layouts
keyboard --vckeymap=us --xlayouts=''
# System language
lang en_US.UTF-8

# Root password
rootpw --plaintext arastra
# Do not configure the X Window System
skipx
# System timezone
timezone America/Los_Angeles --ntpservers=ntp.aristanetworks.com
# System bootloader configuration
bootloader --location=mbr --boot-drive=sda --timeout=5

# Auto partitioning sets up swap which we don't want
#autopart --type=lvm
part raid.11 --asprimary --fstype="mdmember" --ondisk=sda --size=20000
part raid.21 --asprimary --fstype="mdmember" --ondisk=sdb --size=20000
part raid.31 --asprimary --fstype="mdmember" --ondisk=sdc --size=20000
part raid.12 --asprimary --fstype="mdmember" --ondisk=sda --size=1 --grow
part raid.22 --asprimary --fstype="mdmember" --ondisk=sdb --size=1 --grow
part raid.32 --asprimary --fstype="mdmember" --ondisk=sdc --size=1 --grow
#
raid / --device=1 --fstype="xfs" --level=raid1 raid.11 raid.21 raid.31
raid /persist --device=2 --fstype="xfs" --level=raid0 raid.12 raid.22 raid.32

# Partition clearing information
clearpart --all --initlabel --drives=sda,sdb,sdc


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
-kexec-tools
-NetworkManager*
-plymouth*
-postfix
-uboot-tools
-parted
nfs-utils
wget
lsof
kexec-tools
ntp
net-tools
%end

%addon com_redhat_kdump --disable --reserve-mb='128'

%end

%post

# logs all post install actions to the console
set -x
echo "Following commands are being run from \%post section in KS"

# base useful packages
yum install -y systemd-networkd systemd-resolved
yum install -y epel-release vim-minimal
yum install -y cpufrequtils sysstat hwloc-gui.x86 ipmitool pdsh
yum install -y dstat htop strace tcpdump gdb xfsprogs tar lsof
yum install -y pciutils
# yum-utils contains debuginfo-install used to setup debuginfo for Scylla
yum install -y yum-utils
# for Scylla debugging
yum install -y perf

# Configure bonding
mkdir /tmp/root
mount nfs101:/persist2/milkyway /tmp/root
# remove ifcfg scripts since we use systemd
rm /etc/sysconfig/network-scripts/ifcfg-eno1
rm /etc/sysconfig/network-scripts/ifcfg-eno2

export ENO1_MAC_ADDR=`ip link show dev eno1 | awk '/ether/ {print($2)}'`
/tmp/root/SetupSystemdBond.py $ENO1_MAC_ADDR

systemctl enable systemd-networkd
systemctl enable systemd-resolved
systemctl disable firewalld

# ASB support
yum install -y git mysql ansible smartmontools

# Arastra Default Account Setup
group --name=arastra --gid=10000
user --name=arastra --plaintext --password=arastra --uid=10000 --gid=10000 --gecos="Arista Networks anonymous account" --groups=wheel,root
mkdir -p ~arastra/.ssh
cp /tmp/root/arastra-ssh/* ~arastra/.ssh/
chown arastra:arastra -R ~arastra/.ssh
chmod 700 ~arastra/.ssh
chmod 600 ~arastra/.ssh/*

# Ansible Generic Account Setup
group --name=ansible --gid=20000
user --name=ansible --uid=20000 --gid=20000 --gecos="Ansible Account" --groups=wheel,root
mkdir -p ~ansible/.ssh
cp /tmp/root/ansible-ssh/authorized_keys ~ansible/.ssh/
chown ansible:ansible -R ~ansible/.ssh
chmod 700 ~ansible/.ssh
chmod 644 ~ansible/.ssh/authorized_keys

cat <<EOF >> /etc/sudoers
Defaults !env_reset, !requiretty
%wheel ALL=(ALL) NOPASSWD: ALL
EOF

# setup ToolsV2 repo
cat <<EOF > /etc/yum.repos.d/ToolsV2.repo
[ToolsV2]
name=ToolsV2 \$basearch
baseurl=http://tools/ToolsV2/repo/\$basearch/RPMS/
enabled=1
gpgcheck=0
metadata_expire=2h
EOF

# enable nightly update of ToolsV2 repo
# we don't enable automatic nightly update for any other repo
cp /tmp/root/ToolsV2-update /etc/cron.daily/
chmod a+x /etc/cron.daily/ToolsV2-update
chmod a+x /etc/cron.daily/logrotate

yum install -y lldpad
systemctl enable lldpad
lldptool -i eno1 -T -V sysName enableTx=yes
lldptool -i eno2 -T -V sysName enableTx=yes

echo 'server ntp.aristanetworks.com' >> /etc/ntp.conf
systemctl enable ntpd

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

# since all milkyway servers use SSDs
# enable TRIM support on raid0
cp /tmp/root/raid0.conf /etc/modprobe.d/
cp /tmp/root/fstrim.scylla /etc/cron.daily/fstrim
chmod a+x /etc/cron.daily/fstrim

tuned-adm off
systemctl disable tuned
%end
