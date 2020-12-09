#
#  Copyright (c) 2016 General Electric Company. All rights reserved.
#
#  The copyright to the computer software herein is the property of
#  General Electric Company. The software may be used and/or copied only
#  with the written permission of General Electric Company or in accordance
#  with the terms and conditions stipulated in the agreement/contract
#  under which the software has been supplied.
#
#  author: apolo.yasuda@ge.com
#

ecauthapi=authapi

.DEFAULT_GOAL: $(ecauthapi)

$(ecauthapi): authapi-lint authapi-test authapi-build

pre-install:
	@ls -la

authapi-build:
	@echo creating artifact..
	@go build -o ./agent .

authapi-sast:
	@echo begining SAST scanning..
	@gosec -exclude-dir=src ./...

authapi-lint:
	echo begining LINT checking..
	@golint ./...

authapi-test:
	@echo begining test..
	@go test -vet=off
	
.PHONY: install
install:
	ls -al

