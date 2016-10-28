# Ansible config for ops (r12s1 to r12s32 machines)

# Install ansible

```sh
brew install ansible
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
