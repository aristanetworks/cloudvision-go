# ardc-config
This repository is used to host playbooks, config files (e.g. hosts file like 'ansible\_hosts'), and test framework that is expected to be used by Ansible for datacenter server provisioning and maintenance.


## Conversion from AroraConfig
Having both Ansible and AroraConfig manage the same files is bad.  So we need an interlock that lets new ansible playbooks and use to know then it is safe to take over.

### Step 1. Write playbook with the following structure.
Make sure to write the appropriate package name you are converting between the brackets. When in doubt, refer to example playbooks already written under 'ardc-config/playbooks/'.

```
  vars:
     aroraConfig[ package name ]:       "/var/lib/AroraConfig/[ package name ]"
  tasks:
     - stat: path="{{aroraConfig[ package name ]}}"
       register: fix[ package name ]

     - block:
        [  stuff you want to do goes here ]

        - name: "TEST"
          script: test/test_[ this playbook file name ]
          when: TEST is defined

       when: fix[ package name ].stat.exists == True

```

### Step 2. Gut the AroraConfig.spec, and replace/add post action to touch the sentinel file.
```
%post
# Create sentinel file used by Ansible playbook (postfix.yml) to indicate
# that AroraConfig no longer manages postfix.
mkdir -p /var/lib/AroraConfig
touch /var/lib/AroraConfig/postfix
```

### Step 3. Write corresponding test file in 'ardc-config/playbooks/test/'.
Make sure that the test file you write has a name in the form of 'test\_[ corresponding playbook name ]'.
You can write the test in any way you'd prefer, but must:
    * Upon SUCCESSFUL completion, gracefully EXIT.
    * Upon FAIL, ABORT with error code.


## Folder/File Details
### ansible\_hosts
TODO

### playbooks
TODO

### test
TODO

