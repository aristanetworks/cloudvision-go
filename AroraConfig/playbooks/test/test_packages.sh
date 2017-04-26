# Tests packages are installed
rpm -qa telegraf-Linux | grep telegraf-Linux

# Tests packages are removed
! rpm -qa prelink | grep prelink
