#!/bin/sh

projBin=$GOPATH/src/z/app

sources=`find $projBin -name "server"`
for source in $sources;do
	echo "rm "$source
	rm $source
done
