#!/usr/bin/env bash

declare -a arr=( # http://stackoverflow.com/questions/8880603/loop-through-array-of-strings-in-bash-script
#	"amount"
#	"balance"
	"user_contact_json"
#	"user"
#	"bill_member_json"
#	"group"
#    "transfer_fromto"
)

for f in "${arr[@]}"
do
	echo "Removing ${f}_ffjson.go...."
	rm "${f}_ffjson.go"
done

for f in "${arr[@]}"
do
   echo "Regenerating ${f}..."
   ffjson "${f}.go"
done