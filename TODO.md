# TODO and BUGS

## TODO CLI

## TODO API

## TODO GENERAL

* [ ] Add testing

## BUGS

* [ ] Saving kubeconfig fails from time to time. One: Depends on machine clock (wsl lags time sync), Two: looks like it doesn't completely overwrite opened file (saved data is corrupted if existing kubeconfig contains data in it)

## IMPLEMENTED

* [x] Add Github Actions workflow for building/publishing
* [x] Implement create/delete for workergroups (`cl shoot create --workergroup`). Only `--cluster` option is currently available.
* [x] Add API calls to hibernate/wake-up given shoot clusters.
* [x] Implement hibernate/wake-up functionality for shoot clusters.
