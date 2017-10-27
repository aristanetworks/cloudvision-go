push:
	# update dist
	a4 ssh dist git --git-dir=/dist/storage/ardc-config/.git --work-tree=/dist/storage/ardc-config pull --rebase

jenkins:
	$(MAKE) -C ansible

.PHONY: jenkins
