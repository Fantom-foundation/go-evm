#!/usr/bin/env sh

SCRIPT="`readlink -f "$0"`"
SCRIPTPATH="${SCRIPT%/*}"
REPO="${SCRIPTPATH%/*}"
TMPDIR="${TMPDIR%/:-'/tmp/'}"
BINDIR="$REPO/usr_bins"

tarball="$TMPDIR"base.tgz
extract_to="$TMPDIR"netbsd_sets_base

[ -f "$tarball" ] || wget ftp://ftp.netbsd.org/pub/NetBSD/NetBSD-8.0/amd64/binary/sets/base.tgz -O "$tarball" --show-progress

mkdir -p "$BINDIR" "$extract_to"

[ -n "`ls -A $extract_to`" ] || tar xf "$tarball" -C "$extract_to"

for bin in ls mv rm mkdir sh; do
  [ -f "$BINDIR/$bin" ] || cp "$extract_to/bin/$bin" "$BINDIR/";
done

for bin in env; do
  [ -f "$BINDIR/$bin" ] || cp "$extract_to/usr/bin/$bin" "$BINDIR/"; 
done
