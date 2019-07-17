#!/usr/bin/env bash

set -e

# Turn colors in this script off by setting the NO_COLOR variable in your
# environment to any value:

NO_COLOR=${NO_COLOR:-""}
if [ -z "$NO_COLOR" ]; then
  header=$'\e[1;33m'
  reset=$'\e[0m'
else
  header=''
  reset=''
fi

while [[ $# -gt 0 ]]
do
key="$1"

case $key in
    -n|--namespace)
    namespace="$2"
    shift # past argument
    shift # past value
    ;;
    -p|--pod)
    pod="$2"
    shift # past argument
    shift # past value
    ;;
    -t|--type)
    type="$2"
    shift # past argument
    shift # past value
    ;;
esac
done


function header_text {
  echo "$header$*$reset"
}


echo "Using namespace/Project $namespace"
oc project $namespace

echo "Setting Label for $namespace"
oc label namespace $namespace knative-eventing-injection=enabled --overwrite

if [ "$type" == "prod" ]; then
    header_text "Starting Production deployment on OpenShift!"
    for filename in deploy/templates/*.yaml; do
        header_text "Applying configuration $filename on $namespace"
        kubectl -n $namespace apply -f $filename
    done
fi
if [ "$type" == "staging" ]; then
    header_text "Starting Staging deployment on OpenShift!"
    for filename in deploy/templates/*.yaml; do
        header_text "Applying configuration $filename on $namespace"
        sed 's/prod/staging/ ; s/PROD/STAGING/' $filename | kubectl -n $namespace apply -f -
    done
fi


header_text "Starting rollout of $type release"
#oc rollout latest dc/$pod