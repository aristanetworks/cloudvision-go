
# ardc-config
This repository is used to host playbooks, config files (e.g. hosts file like 'ansible\_hosts'), and test framework that is expected to be used by Ansible for datacenter server provisioning and maintenance.

## Ansible-push vs. Ansible-pull
The main design of Ansible is that there is a centralized server from which a set of "playbooks" or configuration modules are pushed to the remote servers. There is also an opposite design, called [Ansible-pull](http://docs.ansible.com/ansible/playbooks_intro.html#ansible-pull), which is essentially the opposite - instead of having a centralized server, the remote servers themselves pull from a git repo that contains the playbooks that need to be run, and run them independently by themselves. We support both models, although we rely more on Ansible-pull model than push. 

## How to Write Playbooks
Description of what a playbook is and how to write them are available in Ansible's own [documentation](http://docs.ansible.com/ansible/playbooks_intro.html). 

