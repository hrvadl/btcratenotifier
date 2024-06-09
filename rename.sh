#!/usr/bin/env sh

export CUR="github.com/hrvadl/converter"
export NEW="github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl"
go mod edit -module ${NEW}
find . -type f -name '*.go' -exec perl -pi -e 's/$ENV{CUR}/$ENV{NEW}/g' {} \;

