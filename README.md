# Ultrastar RSS adder

This script will add all torrents from an [ultrastar-es.org] RSS feed to transmission, using its [RPC interface](https://github.com/transmission/transmission/blob/master/extras/rpc-spec.txt).

## Usage

The simplest invocation is

`ULTRASTAR_RSS="<your rss url here>" ultrastar`

or, if you want these torrents to be added to a specific folder different from your default transmission folder,

`ULTRASTAR_RSS="<your rss url here>" ULTRASTAR_TARGET="$HOME/Music/Ultrastar" ultrastar`

## Installing

`go get -u git.sr.ht/~vicentereyes/ultrastar`

## Needed configuration

If you get a connection refused error, it's probably because transmission isn't running, or because it doesn't have remote access enabled.
Go to Edit->Preferences->Remote->Allow remote access and leave the rest as default.