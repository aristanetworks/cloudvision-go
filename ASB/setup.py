#!/usr/bin/env python
# Copyright (c) 2016 Arista Networks, Inc.  All rights reserved.
# Arista Networks, Inc. Confidential and Proprietary.

from setuptools import setup, find_packages

py_modules=[]
install_requires = []
scripts = []
tests_requires = []
packages = find_packages()

setup(
   name='ASB',
   version='1.0',
   description='Automatic Server Bringup',
   author='Ren Lee, Leo Kam, Adam Setters',
   author_email='ren@arista.com',
   license='Arista Networks',
   url='http://gerrit/ardc-config',

   py_modules=py_modules,
   install_requires=install_requires,
   packages=packages,
   scripts=scripts,
   test_suite='test',
   tests_require=tests_requires

)
