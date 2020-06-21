#!/bin/bash -eu

## GO related functions that can be used by other scripts, use source lib/scripts/go.sh to use
## Tim Dadd : Genesis June 2020

# go fmt ./...
function goFmt() {
  echo -n "${CYAN}Format check: "
  fmtErrors=$(find . -type f -name \*.go | xargs gofmt -l 2>&1 || true)
  if [ "${fmtErrors}" ]; then
      for f in ${fmtErrors}; do
        case "$f" in
        *vendor/gopkg.in/yaml.v2/* | *lib/*)
#            echo "Ignore formatting of $e"
          ;;
        *)
          echo -n "($f): "
          go fmt $f
        esac
      done
  fi
  echo "${GREEN}ALL FORMATTED$WHITE"
}

# go vet ./...
# Run `go vet` against all targets. If problems are found - print them to stderr (&2)
function goVet() {
  echo -n "${CYAN}Vetting..."
  vetErrors=$(go vet ./... 2>&1 || true)
  if [ -n "${vetErrors}" ]; then
      echo "${RED}FAIL"
      echo "${vetErrors}${WHITE}"
      exit 1
  fi
  echo "${GREEN}VETTED OK$WHITE"
}