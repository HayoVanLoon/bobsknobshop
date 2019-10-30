#!/usr/bin/env bash

# This file contains project-specific variables that might vary developer to
# developer. If they do not fit your configuration, copy them over to your
# personal-envs.sh.
#
# This is a shared file, DO NOT UPDATE WITH PERSONAL VALUES unless they make
# good defaults.

# General Variables
export PROJECT_ID=hayovanloon-0
export GOOGLE_CLOUD_PROJECT=hayovanloon-0

# Project Variables
export PROJECT_ORGANISATION=bobsknobshop
export PROJECT_BASE_VERSION=v1

# Expected to be set in personal-envs.sh
# Placed here for reference, DO NOT EDIT
export GOOGLE_ACCOUNT=
export GOOGLE_APPLICATION_CREDENTIALS=
export PYTHON27_EXEC=$(which python2.7)
export PYTHON35P_EXEC=$(which python3.6)
export VENV_EXEC=$(which virtualenv)
export PROTO_GOOGLE_APIS=
export MAKE=$(which make)
export PROTOC_EXEC=$(which protoc)


if [[ -f personal-envs.sh ]]; then
    . personal-envs.sh
fi

if [[ -z ${GOOGLE_ACCOUNT} ]]; then
    echo "Error: Update personal-envs.sh with your personal settings first."
else
    if [[ -z "${PROJECT_ID}" || \
            "${PROJECT_ID}" = "PLACEHOLDER_*" ]]; then
        echo "Error: PROJECT_ID not set"
    else
        gcloud config set account ${GOOGLE_ACCOUNT} &> /dev/null
        gcloud config set project ${PROJECT_ID} &> /dev/null
    fi
fi

