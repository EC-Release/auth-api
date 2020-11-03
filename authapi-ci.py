__author__ = 'Apolo Yasuda <apolo.yasuda@ge.com>'

'''
   EC AUTH API deployment script
'''

import  os, json, base64, sys, threading, subprocess
from time import sleep

from common import Common

c=Common(__name__)

DIST="dist"
SDK="auth-api"
ART="api"
BINARY="{}/{}/{}".format(SDK,DIST,ART)
CKF = 'checksum.txt'
    

# CI_COMMIT_BRANCH=os.environ["CI_COMMIT_BRANCH"]
# CI_DEFAULT_BRANCH=os.environ["CI_DEFAULT_BRANCH"]
# CI_JOB_ID=os.environ["CI_JOB_ID"]
# GITHUB_TKN=os.environ["GITHUB_TKN"]
# DEPLOY_BRANCH="v1beta"
# API_BUILD_REV=os.environ["API_BUILD_REV"]

CI_COMMIT_BRANCH=sys.argv[1]
CI_DEFAULT_BRANCH=sys.argv[2]
CI_JOB_ID=sys.argv[3]
GITHUB_TKN=sys.argv[4]
DEPLOY_BRANCH="v1beta"
API_BUILD_REV=sys.argv[5]

EC_TAG=""

def main():

    if DEPLOY_BRANCH!=CI_COMMIT_BRANCH:
        print("discard the deployment step: not on the deployment branch: {}. Was on {}".format(DEPLOY_BRANCH,CI_COMMIT_BRANCH))
        return

    print("clonning external sdk..")
    os.system("git clone --depth 1 --branch {} https://{}@github.com/EC-Release/auth-api.git /{}".format(DEPLOY_BRANCH, GITHUB_TKN, SDK))

    print("remove existing dist..")
    os.system("rm /{}/*".format(BINARY))
    
    print("generate linux artifacts..")
    os.system("CGO_ENABLED=0 GOOS=linux GODEBUG=netdns=cgo GOARCH=amd64 go build -tags netgo -a -v -o /{}/{}_linux *.go".format(BINARY,ART))

    c.chksumgen('/{}'.format(BINARY),CKF)
    
    fl = os.listdir('/{}'.format(BINARY))
    for filename in fl:
        if filename==CKF:
            continue
        
        os.system('cd /{}; tar -czvf {}.tar.gz ./{}'.format(BINARY, filename, filename))
        os.system('rm /{}/{}'.format(BINARY,filename))
    
    os.system("ls -al /{}".format(BINARY))
    os.system("cd /{}; git add {}/{}".format(SDK,DIST,ART))
    os.system("cd /{}; git config user.name '{}' && git config user.email '{}'; git config core.fileMode false".format(BINARY,"EC.Bot","EC.Bot@ge.com"))
    os.system("cd /{}; git commit -m '{} job#{} checked-in'".format(BINARY, ART, CI_JOB_ID))

    #os.system("git tag {}".format(""))
    print("update the artifacts in auth-api..")
    os.system("cd /{}; git tag {}.auth-api.{}; git push -f origin {} --tags".format(BINARY,DEPLOY_BRANCH,API_BUILD_REV,DEPLOY_BRANCH))

    return

if __name__=='__main__':
    main()
