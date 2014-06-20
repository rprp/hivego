#!/usr/bin/python2.7
# Filename: build.py

import os
import string

curdir=os.getcwd()

oldpath=os.getenv("GOPATH")
oldbin=os.getenv("GOBIN")
oldarch=os.getenv("GOARCH")
oldos=os.getenv("GOOS")

os.environ['GOBIN'] = curdir + '/bin'
os.environ['GOPATH'] = curdir + ':' + oldpath

os.environ['GOARCH'] = 'amd64'
os.environ['GOOS'] = 'linux' 
os.system('go install')

os.environ['GOARCH'] = '386'
os.environ['GOOS'] = 'linux' 
os.system('go install')

os.environ['GOARCH'] = 'amd64'
os.environ['GOOS'] = 'windows' 
os.system('go install')

os.environ['GOARCH'] = '386'
os.environ['GOOS'] = 'windows' 
os.system('go install')

os.environ['GOARCH'] = 'amd64'
os.environ['GOOS'] = 'darwin' 
os.system('go install')

os.system('cp *.toml bin/')

if oldpath is not None:
    os.environ['GOPATH']= oldpath
else:
    os.environ['GOPATH']= ""

if oldbin is not None:
    os.environ['GOBIN']= oldbin
else:
    os.environ['GOBIN']= ""

if oldarch is not None:
    os.environ['GOARCH']= oldarch
else:
    os.environ['GOARCH']= ""

if oldos is not None:
    os.environ['GOOS']= oldos
else:
    os.environ['GOOS']= ""
