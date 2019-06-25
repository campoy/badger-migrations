# badger-migrations
A repository to experiment with badger and all its versions

All of the folders in this repository implement a very similar CLI tool
that is able to set, get, delete, and list values from a Badger DB.

They all store their data under `db`, and the idea is to be able to see
whether different version are data compatible.

For each Badger version there are two directories: one that vendors its
dependencies, and one that uses Go Modules which have a `-mods` suffix
on the folder name.