#!/bin/bash

# Copyright 2015 tsuru authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

status=0
out=`gofmt -s -l .`
if [ "${out}" != "" ]
then
    echo "ERROR: there are files that need to be formatted with gofmt"
    echo
    echo "Files:"
    for file in $out
    do
        echo "- ${file}"
    done
    echo
    status=1
fi

out=`goimports -l .`
if [ "${out}" != "" ]
then
    echo "ERROR: there are files that need to be formatted with goimports"
    echo
    echo "Files:"
    for file in $out
    do
        echo "- ${file}"
    done
    status=1
fi

out=`go tool vet -shadow -all . 2>&1`
if [ "${out}" != "" ]
then
    echo "ERROR: go vet failures:"
    echo
    cat <<END
${out}
END
    status=1
fi

exit $status
