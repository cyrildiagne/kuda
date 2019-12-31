#!/bin/sh

{ # Prevent execution if this script was only partially downloaded

  set -e

  VERSION=0.3.3-preview

  green="\033[32m"
  red="\033[31m"
  reset="\033[0m"

  install_path='/usr/local/bin/kuda'

  OS=$(uname | tr '[:upper:]' '[:lower:]')
  ARCH=$(uname -m | tr '[:upper:]' '[:lower:]')

  cmd_exists() {
    command -v "$@" >/dev/null 2>&1
  }

  latestURL=https://github.com/cyrildiagne/kuda/releases/download/v$VERSION

  case "$OS" in
  darwin)
    URL=${latestURL}/kuda_${VERSION}_Darwin_x86_64.tar.gz
    ;;
  linux)
    case "$ARCH" in
    x86_64)
      URL=${latestURL}/kuda_${VERSION}_Linux_x86_64.tar.gz
      ;;
    amd64)
      URL=${latestURL}/kuda_${VERSION}_Linux_x86_64.tar.gz
      ;;
    *)
      printf "$red> The architecture (${ARCH}) is not supported by this installation script.$reset\n"
      exit 1
      ;;
    esac
    ;;
  *)
    printf "$red> The OS (${OS}) is not supported by this installation script.$reset\n"
    exit 1
    ;;
  esac

  sh_c='sh -c'
  if [ ! -w $install_path ]; then
    # use sudo if $user doesn't have write access to the path
    if [ "$user" != 'root' ]; then
      if cmd_exists sudo; then
        sh_c='sudo -E sh -c'
      elif cmd_exists su; then
        sh_c='su -c'
      else
        echo 'This script requires to run commands as sudo. We are unable to find either "sudo" or "su".'
        exit 1
      fi
    fi
  fi

  printf "> Installing Kuda...\n"
  printf "  from $URL\n"
  printf "  to $install_path\n\n"

  $sh_c "rm -f $install_path"
  $sh_c "curl -fSL $URL -o kuda-download.tar.gz"
  $sh_c "tar -xf kuda-download.tar.gz kuda"
  $sh_c "mv kuda $install_path"
  $sh_c "chmod +x $install_path"
  $sh_c "rm kuda-download.tar.gz"

  printf "$green> Kuda successfully installed!\n$reset"

} # End of wrapping
