#!/bin/bash

packages=($(find . -name "go.mod" -print0 | xargs -0 -n1 dirname | sort --unique))
SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]:-$0}"; )" &> /dev/null && pwd 2> /dev/null; )";


function cmd_on_packages() {
	echo "$1 packages..."
	for package in "${packages[@]}"; do
		pushd $SCRIPT_DIR/$package &> /dev/null
		echo -e "\n${1}: $(go list -m)"
		for cmd in "${@:2}"; do
			($cmd)
		done
		popd &> /dev/null
	done
}

function get_package_name() {
	pushd $SCRIPT_DIR/$1 &> /dev/null
	package=$(go list -m)
	popd &> /dev/null
	echo $package
}

function input_number_with_max() {
	local max=$1

	while true; do
		read  -n 1 -p "Choose a number(1-$1): " input
		echo -e "" > /dev/tty

		if ((input >= 1 && input <= $1)); then
			break;
		fi
	done

	echo $max
}


function release() {
	# This explains releases of subpackages in go
	# https://github.com/go-modules-by-example/index/blob/master/009_submodules/README.md

	echo "Which package would you want to release"

	local package_names=()
	for i in "${!packages[@]}"; do
		package=${packages[$i]}
		package_name=$(get_package_name ${packages[$i]})
		package_names+=($package_name)

		printf "%s\t%s\t%s\n" "$(($i+1))" "$package" "$package_name"
	done

	# input_index=$(($(input_number_with_max ${#packages[@]} "Choose a number")-1))
	input_index=0

	echo -e "\nFinding tags associated with: ${package_names[$input_index]}"
	# echo ${package_names[$input_index]}
	cut_string="${SCRIPT_DIR##*/}/"
	# echo $cut_string
	tag_prefix="$(echo ${package_names[$input_index]} | awk -F $cut_string '{print $2}')/v"

	tags=$(git tag --list | grep $tag_prefix)
	echo "Tags found: $tags"

	# TODO finish
}

## Choose which function to use by argument
case $1 in
	t | -t | test | --test)
		cmd_on_packages "Testing" \
			"go clean -testcache" \
			"go test ./..."
		;;
	d | -d | tidy | --tidy)
		cmd_on_packages "Tidying" "go mod tidy"
		;;
	u | -u | update |--update)
		cmd_on_packages "Updating" \
			"go get github.com/mvndaai/ctxerr" \
			"go get -u ./..." \
			"go mod tidy"
		;;
	r | -r | release | --release)
		release
		;;
	*)
		echo "Usage: $(basename $0) [OPTIONS]

Options:
-t, --test          Runs 'go test ./...' for ctxerr and all subpackages
-d, --tidy          Runs 'go mod tidy ./...' for ctxerr and all subpackages
-u, --update        Updates and tidyies ctxerr and all subpackages
-r, --release       Releases a package by adding a version tag"
		;;
esac