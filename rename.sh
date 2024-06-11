#!/usr/bin/env sh

export CUR="github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/pkg/logger"
export NEW="github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/kit/logger"
go mod edit -module ${NEW}
find . -type f -name '*.go' -exec perl -pi -e 's/$ENV{CUR}/$ENV{NEW}/g' {} \;

