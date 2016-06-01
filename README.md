
# ardc-config
This repository is used to host playbooks, config files (e.g. hosts file like 'ansible\_hosts'), and test framework that is expected to be used by Ansible for datacenter server provisioning and maintenance.

<br />
<br />

## Ansible-push vs. Ansible-pull
The main design of Ansible is that there is a centralized server from which a set of "playbooks" or configuration modules are pushed to the remote servers. There is also an opposite design, called [Ansible-pull](http://docs.ansible.com/ansible/playbooks_intro.html#ansible-pull), which is essentially the opposite - instead of having a centralized server, the remote servers themselves pull from a git repo that contains the playbooks that need to be run, and run them independently by themselves. We support both models, although we rely more on Ansible-pull model than push. 

<br />
<br />

## How to Write Playbooks
Description of what a playbook is and how to write them are available in Ansible's own [documentation](http://docs.ansible.com/ansible/playbooks_intro.html). 

### local.yml
This playbook is the top-level master playbook that dictates which playbooks should be run, as well as the order in which they should run. If your newly written playbook is not included in this master playbook, it will not get run. *Order of playbooks does matter.* The playbooks included in this file will be played in sequential order. This playbook is also the playbook that will be run when ansible-pull is invoked. 

### If you are converting from AroraConfig
Having both Ansible and AroraConfig manage the same files is bad.  So we need an interlock that lets new ansible playbooks to know when it is safe to take over. The following steps show how to write the playbook so that it is compatible with changes made from AroraConfig.

#### Step 1. Write playbook with the following structure.
Replace the 'PACKAGE' with appropriate package name and 'PLAYBOOK' with the playbook name. When in doubt, refer to example playbooks already written in 'ardc-config/playbooks/'.

```
   ---
   - hosts: all
     become: yes
     become_user: root
     vars:
        aroraConfigPACKAGE:       "/var/lib/AroraConfig/PACKAGE"
        # other variables you want to define goes here...

     tasks:
        - stat: path="{{aroraConfigPACKAGE}}"
          register: fixPACKAGE

        - block:
           # tasks you want to write goes here...

           - name: "TEST"
             script: test/test_PLAYBOOK.py
             when: TEST is defined

          when: fixPACKAGE.stat.exists == True
   ...
```

#### Step 2. Gut the AroraConfig.spec, and replace/add post action to touch the sentinel file.
```
   %post
   # Create sentinel file used by Ansible playbook ( PACKAGE.yml ) to indicate
   # that AroraConfig no longer manages PACKAGE.
   mkdir -p /var/lib/AroraConfig
   touch /var/lib/AroraConfig/PACKAGE
```

#### Step 3. Write corresponding test file in 'ardc-config/playbooks/test/'.
Each playbook should have a corresponding test file. The name of the test file must be the name of the playbook you want to test prepended with "test\_" ( e.g.: test\_pkg.py for 'pkg.yml' )
You can write the test in any way you'd prefer (any file type, with any structure), but must:
    - Upon SUCCESSFUL completion, gracefully EXIT.
    - Upon FAIL, ABORT with error code.
With these tests what we care for the most is whether the expected post-conditions hold true after running the playbooks.

#### Step 4. Add your playbook to "ardc-config/local.yml" 
Add the relative path to your playbook below indicated line in the file. 
```
   - include: playbooks/PLAYBOOK.yml
```

### If you are not converting from AroraConfig, but writing a completely new config playbook
The structure of the playbook should look similar as above ( i.e.: hosts, become, become\_user, vars, tasks, TEST task ) except checking for the sentinel file. For example:

```
   ---
   - hosts: all
     become: yes
     become_user: root
     vars:
        # variables you want to define goes here...

     tasks:
        # tasks you want to write goes here...

        - name: "TEST"
          script: test/test_PLAYBOOK.py
          when: TEST is defined
   ...
```

### Getting your Playbook Code Reviewed
Copy and paste the contents of 'gitconfig-review' to .git/config. Add and commit the files (your playbooks) that you wish to be reviewed. Then do "git push review". This will automatically open a review ticket on gerrit and notify 'ardc-config' maintainers.

<br />
<br />

## Other Notable Files/Folders
### playbooks
Place where all the playbooks should be.

### playbooks/test
Place where playbook tests should be.

#### test/dockerfiles
Contains subdirectories and files required to create the appropriate Docker image needed for testing.

#### test/ping.yml
Main "ping" test playbook used only by the testing framework. This playbook is separate from local.yml and will not be run in deployment environment.

#### test/TestAnsiblePush.py
Test framework. In a nutshell, it simulates a datacenter environment with a central Ansible server and end nodes, and runs the playbooks against them so we can detect errors and observe outcomes of the playbooks. For more detail, please view the header comments section in this test file. The test framework models Ansible push instead of pull, but the end result of running the playbooks should be the same, so currently there are no plans on providing another framework to test Ansible pull.

<br />
<br />

## Dependent Packages
There are packages that ardc-config is dependent on, mainly only for testing purposes.

#### docker-library
'docker-library/library/ardc-config' holds the manifest file for the Docker image used to create the Docker containers for testing. This manifest file should only be updated when there is a critical need for the Docker file describing the Docker image needs to be changed. The manifest file builds and registers the custom Docker image to Arista's own Docker hub registry, from which the test framework pulls the needed Docker image from. If all of that went over your head, feel free to shoot an email to ren@arista.com or view this [README](http://gerrit/plugins/gitblit/docs/?r=docker-library.git&h=master).
