# Ansible config for ops (r12s1 to r12s32 machines)

# LIMITATIONS

* **k8s MASTER**: We only support ONE kubernetes master instance in the cluster right now.
* **ETCD2**: ETCD2 instances are NOT managed by ansible. They are managed manually. The role defined on ansible is just a basic role to have the etcd2 config correctly setup on the host running etcd2 on reinstall/restart/ansible reapply.


# Install ansible

```sh
brew install ansible@2.0
brew link --force ansible@2.0
```

# Install the remote coreos machines:

```sh
ansible-playbook playbook.yml -i inventories/dev/hosts -l r12sXX
```

# Run ansible commands

You want to run your ansible command from this specific directory, so the `ansible.cfg` file in this directory is used.

## Basic commands

### Check if all the machines are ok
```sh
ansible all -i inventories/dev/hosts -m ping
```

### Ping one specific machine

```sh
ansible r12s1 -i inventories/dev/hosts -m ping
```

### Generate cloud-config.yml file for the ops.git/coreos repo/folder

The following command will generate all the cloud-config-IPADDRESS.yml files for each machine in the coreos group cluster.
The generated files will be placed in ~/git/ops/coreos, because the default values for the ops repo path is `~/git/ops`.

```sh
ansible-playbook playbook.yml -i inventories/dev/hosts -l localhost
```

You can override the ops repo path by using the OPS_REPO env var.
For instance, if your ops repo is `~/projects/common/ops`, the command will be:

```sh
OPS_REPO=~/projects/common/ops ansible-playbook playbook.yml -i inventories/dev/hosts -l localhost
```

# Add a new ssh key / Update an existing ssh

The ssh keys are now in the file `group_vars/coreos/ssh_keys`.
After adding/updating/removing an ssh key from this file, the cloud-config files need to be generated again (see previous section), and the ssh keys can be deployed again to the entire cluster by using the following command:

```sh
ansible-playbook playbook.yml -i inventories/dev/hosts
```
