#!/usr/bin/env bash
###
# deploy all the apps
###
app=$1
branch=$2
respository=$3
projectsRoot="/home/kaifa/apps"

if [ ! -d "$projectsRoot" ];then
    mkdir -p "$projectsRoot"
fi

# shellcheck disable=SC2164
cd "$projectsRoot"

appDir="${projectsRoot}/${app}"
if [ ! -d "$appDir"  ];then
    clone="git clone -b  ${branch} ${respository} ${app}"
    echo `$clone`
fi

# shellcheck disable=SC2164
cd "$appDir"

git checkout "${branch}"

git pull --rebase

# shellcheck disable=SC2164
sh "${appDir}/start.sh"

# usage: sh deploy.sh avian develop http://oauth2:bYqRf6xE6UDC_yLtuPKj@gitlab.gosccap.cn/bourse/avian.git
