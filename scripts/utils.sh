bold="\033[1m"
blue="\e[36m"
green="\033[32m"
red="\033[31m"
reset="\033[0m"

function error() {
  printf "${red}ERROR:${reset} $1\n"
}

function assert_set() {
  var_name=$1
  var_value=$2
  if [ -z "$var_value" ]; then
    error "Missing required env variable $var_name"
    exit 1
  else
    printf "$var_name: ${blue}$var_value${reset}\n"
  fi
}