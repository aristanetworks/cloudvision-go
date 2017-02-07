# Ansible config for ops (r12s1 to r12s32 machines)

# Install ansible

```sh
brew install ansible@2.0
brew link --force ansible@2.0
```

# Initialize the remote coreos machines:

**This has to be done only once per machine. Usually you shouldn't need to do this.**

```sh
ansible-playbook bootstrap.yml
```

# Run ansible commands

You want to run your ansible command from this specific directory, so the `ansible.cfg` file in this directory is used.

## Basic commands

### Check if all the machines are ok
```sh
ansible all -m ping
```

### Ping one specific machine

```sh
ansible r12s1 -m ping
```

### Generate cloud-config.yml file for the ops.git/coreos repo/folder

The following command will generate all the cloud-config-IPADDRESS.yml files for each machine in the coreos group cluster.
The generated files will be placed in ~/git/ops/coreos, because the default values for the ops repo path is `~/git/ops`.

```sh
ansible-playbook playbook.yml -l localhost
```

You can override the ops repo path by using the OPS_REPO env var.
For instance, if your ops repo is `~/projects/common/ops`, the command will be:

```sh
OPS_REPO=~/projects/common/ops ansible-playbook playbook.yml -l localhost
```
