#!/bin/sh


function valid(){
	if [ $2 = "all" ];then
		echo 1
	elif [ `echo $1 | grep -e $2` ];then
		echo 1
	else 
		echo 0
	fi
}

projBin=$GOPATH/src/z/app
server=$2
tags='jsoniter prod'
mode="dev"

if [ $1 = "prod" ];then
    tags='jsoniter prod'
	mode="prod"
elif [ $1 = "test" ];then
	mode="test"
	tags='jsoniter test'
fi

echo build with mode ${mode}
echo build with tags ${tags}

sources=`find $projBin -name "main.go"`
for source in $sources;do
    dir=`dirname $source`
    
    ok=`valid "${source}" "${server}"`
    if [ $ok -gt 0 ];then
    	name=z-`basename $dir`-`basename server`
	    if go build -tags "${tags}" -o $dir/$name $source;then
	       echo "$name build success"
	    else
	       echo "$name build faild"
	       exit 1
	    fi
    fi 
done
