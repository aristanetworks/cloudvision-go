# Ansible config for ops (r123s1 to r123s32 machines)

# LIMITATIONS

* **k8s MASTER**: We only support ONE kubernetes master instance in the cluster right now.
* **ETCD2**: ETCD2 instances are NOT managed by ansible. They are managed manually. The role defined on ansible is just a basic role to have the etcd2 config correctly setup on the host running etcd2 on reinstall/restart/ansible reapply.


# Install ansible + cfssl

Ansible 2.4 minimun is required.

```sh
brew install ansible
# Optionally, if needed:
brew link --force ansible

# Install PyYAML as well (needed to deploy k8s services)
pip install PyYAML

# Install cfssl
go get -u -tags nopkcs11 github.com/cloudflare/cfssl/cmd/cfssl
go get -u github.com/cloudflare/cfssl/cmd/cfssljson
```

# Install the remote coreos machines:

```sh
ansible-playbook cluster.yml -i inventories/dev/hosts -l r123sXX.*
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
ansible r123s1.sjc.aristanetworks.com -i inventories/dev/hosts -m ping
```

### Generate cloud-config.yml file for the ops.git/coreos repo/folder

The following command will generate all the cloud-config-IPADDRESS.yml files for each machine in the coreos group cluster.
The generated files will be placed in ~/git/ops/coreos, because the default values for the ops repo path is `~/git/ops`.

```sh
ansible-playbook cluster.yml -i inventories/dev/hosts -l localhost
```

You can override the ops repo path by using the OPS_REPO env var.
For instance, if your ops repo is `~/projects/common/ops`, the command will be:

```sh
OPS_REPO=~/projects/common/ops ansible-playbook cluster.yml -i inventories/dev/hosts -l localhost
```

# Add a new ssh key / Update an existing ssh

The ssh keys are now in the file `group_vars/coreos/ssh_keys`.
After adding/updating/removing an ssh key from this file, the cloud-config files need to be generated again (see previous section), and the ssh keys can be deployed again to the entire cluster by using the following command:

```sh
ansible-playbook cluster.yml -i inventories/dev/hosts
```

# k8s services for a cluster

k8s services are created by the role **k8s_services**.

TODO: A shared list of services are defined, and they will be deployed for all the clusters.

TODO: Each inventory is having its own list of services as well.

In order to run this role on your local machine, you have to run `kubectl proxy` so ansible can access the k8s master of the cluster being updated.

To run the role against localhost, use this command:

```sh
ansible-playbook cluster.yml -l localhost -i inventories/dev/hosts
```

## Init a new cluster

A new cluster initialization needs to go through some one time initialization actions.
In order to not have those tasks executed by default, an env var needs to be passed in the command line for those actions to be executed.
There are two init actions so far:

* Format the namenode: The env var and its value are `NAMENODE_FORMAT=yesformatnamenode`. This command is dangerous and *MUST NOT* be executed on existing cluster.
In order to protect this action to not be executed even when the env var is set, ansible will pause and ask for a randomly generated code which will be to be copied as the input of the prompt. Only when the code is correct, namenode formatting will be processed.

* Initialize the aeris db schema: The env var and its value are `AERIS_INITDBSCHEMA=yes`.
