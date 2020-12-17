#!/bin/bash

# uncomment this to run the job, otherwise run this in another window after starting nomad
# nomad run grid-intensity-exporter.nomad

STATUS=1
check_for_export_job () {
  # check for the api, then check in the output for the existence of a job
  # being listed as running
  nomad job status -short grid-intensity-exporter  | grep "Status" | grep "running"
  STATUS=$?
}

# initial check to see if we can see that it's running:
while [ $STATUS -ne 0 ]

# loop until we get the successful, zero return status, meaning we can continue
do
  check_for_export_job
done

# we really want this in another step, but this is just to demonstrate
# waiting until we have a passing test
# uncomment to run it
# go test -v -tags=dockerrequired ./integration/test/docker
