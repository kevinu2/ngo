#!/bin/bash

function git_credential {
        local git_home=${1:-./}
        local git_user=${2:-${M_USER}}
        local git_pass=${3:-${M_PASSWD}}
        local timeout=${4:-6000}
        local return_path=`pwd`
        cd ${git_home}
        cat << EOF > .git.tmp
protocol=https
host=github.com
username=${git_user}
password=${git_pass}
EOF
        git config --global credential.helper "cache --timeout=${timeout}"
        cat .git.tmp | git credential approve
        rm -f .git.tmp
        cd ${return_path}
}
function _get_users ()
{
	local _users=`git log --pretty='%aN' | sort | uniq -c | sort -k1 -n -r | awk '{print $2}'`
	echo ${_users}
}

function _get_dates ()
{
	local _date=`git log |awk '/Date:/ {print $2"/"$3"/"$4}' |uniq`
	echo ${_date}
}

function _user_stat ()
{
	local _user=${1}
	local _date=`echo ${2} |sed -e 's#/# #g'`
	if [ ${#_date} -gt 0 ]; then
		local _stat=`git log --author="${_user}" --pretty=tformat: --numstat \
			--since="${_date} 00:00:00" --until="${_date} 23:59:59" \
			| awk '{ add += $1; subs += $2; loc += $1 - $2 } \
			END { printf "added lines: %s, removed lines: %s, total lines: %s ", add, subs, loc }' -`
	else
		local _stat=`git log --author="${_user}" --pretty=tformat: --numstat \
			| awk '{ add += $1; subs += $2; loc += $1 - $2 } \
			END { printf "added lines: %s, removed lines: %s, total lines: %s ", add, subs, loc }' -`
	fi
	echo ${_stat}
}

function _per_user_stat_total ()
{
	local _users=(`_get_users`)
	echo "Total: " 
	for _user in ${_users[@]}; do
		echo -n -e " User: ${_user} Stat: \t"
		_user_stat ${_user} ${_date}
	done
}

function _per_user_stat_daily ()
{
	local _users=(`_get_users`)
	for _date in `_get_dates`; do
		echo ${_date} |sed -e 's#/# #g'
		for _user in ${_users[@]}; do
			echo -n -e " User: ${_user} Stat: \t"
			_user_stat ${_user} ${_date}
		done
	done
}

case ${1} in
	"login")
		echo -n "GIT USER: "; read M_USER
		echo -n "GIT PASSWD: " ; read M_PASSWD
		git_credential
	;;
	"user")
		echo -n "User for Git: "; read user
		echo -n "User for E-mail: "; read email
		git config --global user.name "${user}"
		git config --global user.email "${email}"
	;;
	"push")
		git add -A
		git commit -am "RUN git.sh push"
		git push -f origin `git branch |awk '/^*/ {print $2}'`
	;;
	"stat")
		echo "Top:"
		echo -n " "
		git log --pretty='%aN' | sort | uniq -c | sort -k1 -n -r
		_per_user_stat_total
		_per_user_stat_daily |grep -v 'added lines: ,'	
	;;
	"build")
		git pull
		make dev
	;;
	"release")
		git pull
		make
	;;
	*)
		echo
		echo "Useage:"
		echo -e "\t${0} user"
		echo -e "\t${0} login"
		echo -e "\t${0} push"
		echo -e "\t${0} stat"
		echo
	;;
esac
