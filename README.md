# ardc-config
This repository is used to host playbooks and config files (e.g. hosts
file like 'ansible_hosts') that is expected to be used by Ansible-pull
for datacenter server provisioning and maintenance.

Interactions with AroraConfig
------------------------------------------------------------------------------
Having both Ansible and AroraConfig manage the same files is bad.  So we
need an interlock that lets new ansible playbooks and use to know then
it is safe to take over.

Step 1. Write playbook with the following structure.

  vars:
     aroraConfigPostfix:       "/var/lib/AroraConfig/postfix"
  tasks:
     - stat: path="{{aroraConfigPostfix}}"
       register: fixPostfix
     - block:
        - name: "Do stuff"
          command: "echo Hello"

       when: fixPostfix.stat.exists == True

Step 2. Gut the AroraConfig.spec, and replace/add post action to touch the
        sentinal file.

%post
# Create sentinal file used by Ansible playbook (postfix.yml) to indicate
# that AroraConfig no longer manages postfix.
mkdir -p /var/lib/AroraConfig
touch /var/lib/AroraConfig/postfix

## Folder/File Details
### ansible_hosts
TODO

### playbooks
TODO

### test
TODO

