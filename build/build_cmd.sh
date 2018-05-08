#!/bin/sh

PROJ_ROOT=`dirname $0`
echo "${PROJ_ROOT}"
function valid(){
	if [ $2 = "all" ];then
		echo 1
	elif [ `echo $1 | grep -e $2` ];then
		echo 1
	else 
		echo 0
	fi
}

PROJ_BIN=${PROJ_ROOT}/../cmd
SERVER_NAME=$2

tags='jsoniter'
mode="dev"

if [ $1 = "prod" ];then
    tags='jsoniter prod'
	mode="prod"
fi

echo build with mode ${mode}
echo build with tags ${tags}

sources=`find $PROJ_BIN -name "main.go"`
for source in $sources;do
    dir=`dirname $source`
    
    ok=`valid "${source}" "${SERVER_NAME}"`
    if [ $ok -gt 0 ];then
    	name=`basename ${dir}`
	    if go build -tags "${tags}" -o $dir/$name $source;then
	       echo "$name build success"
	    else
	       echo "$name build faild"
	       exit 1
	    fi
    fi 
done
