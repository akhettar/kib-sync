
![go](go.png) ![opendistro](odfe.jpg)![kibana](kibana.jpg)![arraow](arrow.jpg)![git](git.jpg)
# odfe-kibana-sync 
![main](https://github.com/akhettar/odfe-kibana-sync/workflows/main/badge.svg)

## Overview
This is a command line tool to sync and create the `Kibana object configs` from a given [open distro elk cluster](https://opendistro.github.io/for-elasticsearch-docs/) - see below the help section. This tool assumes that Kibana cluster holds the source of the truth in relation to the configuration files. The sole purpose of this tool is to be run periodically in a given CI pipeline to sync these configuration files with the a given Git repository.

Ideally, this tool should be run in a CI pipeline for a given project that hosts Kibana configuration files - see example Github project.

`Read carefully before usage please`:

* This cli is not a full blown command line interface for the open distro elasticsearch. There is one, see [doc](https://opendistro.github.io/for-elasticsearch-docs/docs/cli/).
* This cli assumes that the all the updates, create and delete are done through the Kiban console.
* This cli is for synching the Kibana configuration files with a given Git repository.
* Deleted configs in Kibana console will be automatically deleted with the `sync folder`.

## Supported Kibana configurations
The following kibana configuration files can be synched and created in a given kiban instance:
1. Monitors
2. Dashboards
3. Saved search
4. Destinations(slack, email etc)
5. Email accounts
6. Email groups
7. Visualization
8. Index pattern


```./odfe-kibana-sync -h
This tool performs the followings:
                1. Fetches configured monitors, dashboards, alert destinations for the given kibana cluster and store them locally as json files.
                2. Pushes the changes done to the monitor's config to Kiban cluster

Usage:
  odfe-kibana-sync [flags]
  odfe-kibana-sync [command]

Available Commands:
  create      create all kiban objects (monitors, dashbaor, etc) present in the config folder
  help        Help about any command
  sync        Fetches Kiban objects (monitor, dashbaord, etc) from Kibana cluster

Flags:
  -h, --help              help for odfe-kibana-sync
      --password string   The kibana cluster password. This is a required argument to connect to the ELK cluster
      --url string        The kibana cluster url. This is required argument to connect to the ELk cluster
      --username string   The kibana cluster username. This is required argument to connect to the ELK cluster
      --workdir string    The working directory where the kibana configuration files will be stored (default "config")

Use "odfe-kibana-sync [command] --help" for more information about a command.
```

## Watch a demo
coming soon!

## Invoking the sync command
The sync command fetches all the monitors config defined in the given kibana cluster and store them locally in the `./config folder`

```
./odfe-kibana-sync sync --username admin --password admin --url https://localhost:9200
```

The above command should produce a `config` folder containing the following structure (of course depending on the kiban objects configured in the cluster)
```
├── config
│   ├── dashboard
│   │   ├── dashboard:722b74f0-b882-11e8-a6d9-e546fe2bba5f.json
│   │   └── dashboard:7adfa750-4c81-11e8-b3d7-01146121b73d.json
│   ├── destination
│   │   ├── 2nvxw3cBIRIg_9xDkrXp.json
│   │   ├── jnv6w3cBIRIg_9xDLPso.json
│   │   └── qXusw3cBIRIg_9xDR7Ru.json
│   ├── email_account
│   │   ├── iHvgw3cBIRIg_9xDSrXh.json
│   │   └── j3v6w3cBIRIg_9xDLftc.json
│   ├── email_group
│   │   └── jHvgw3cBIRIg_9xDmbV5.json
│   ├── index-pattern
│   │   ├── index-pattern:d3d7af60-4c81-11e8-b3d7-01146121b73d.json
│   │   └── index-pattern:ff959d40-b880-11e8-a6d9-e546fe2bba5f.json
│   ├── monitor
│   │   └── qHurw3cBIRIg_9xD6LRW.json
│   ├── search
│   │   ├── search:3ba638e0-b894-11e8-a6d9-e546fe2bba5f.json
│   │   └── search:571aaf70-4c88-11e8-b3d7-01146121b73d.json
│   └── visualization
│       ├── visualization:01c413e0-5395-11e8-99bf-1ba7b1bdaa61.json
│       ├── visualization:08884800-52fe-11e8-a160-89cc2ad9e8e2.json
│       ├── visualization:09ffee60-b88c-11e8-a6d9-e546fe2bba5f.json
│       ├── visualization:10f1a240-b891-11e8-a6d9-e546fe2bba5f.json
│       ├── visualization:129be430-4c93-11e8-b3d7-01146121b73d.json
│       ├── visualization:1c389590-b88d-11e8-a6d9-e546fe2bba5f.json
│       ├── visualization:293b5a30-4c8f-11e8-b3d7-01146121b73d.json
│       ├── visualization:2edf78b0-5395-11e8-99bf-1ba7b1bdaa61.json
│       ├── visualization:334084f0-52fd-11e8-a160-89cc2ad9e8e2.json
│       ├── visualization:37cc8650-b882-11e8-a6d9-e546fe2bba5f.json
│       ├── visualization:45e07720-b890-11e8-a6d9-e546fe2bba5f.json
│       ├── visualization:4b3ec120-b892-11e8-a6d9-e546fe2bba5f.json
│       ├── visualization:707665a0-4c8c-11e8-b3d7-01146121b73d.json
│       ├── visualization:76e3c090-4c8c-11e8-b3d7-01146121b73d.json
│       ├── visualization:8f4d0c00-4c86-11e8-b3d7-01146121b73d.json
│       ├── visualization:9886b410-4c8b-11e8-b3d7-01146121b73d.json
│       ├── visualization:9c6f83f0-bb4d-11e8-9c84-77068524bcab.json
│       ├── visualization:9ca7aa90-b892-11e8-a6d9-e546fe2bba5f.json
│       ├── visualization:aeb212e0-4c84-11e8-b3d7-01146121b73d.json
│       ├── visualization:b72dd430-bb4d-11e8-9c84-77068524bcab.json
│       ├── visualization:b80e6540-b891-11e8-a6d9-e546fe2bba5f.json
│       ├── visualization:bcb63b50-4c89-11e8-b3d7-01146121b73d.json
│       ├── visualization:c8fc3d30-4c87-11e8-b3d7-01146121b73d.json
│       ├── visualization:e6944e50-52fe-11e8-a160-89cc2ad9e8e2.json
│       ├── visualization:ed78a660-53a0-11e8-acbd-0be0ad9d822b.json
│       ├── visualization:ed8436b0-b88b-11e8-a6d9-e546fe2bba5f.json
│       ├── visualization:f8283bf0-52fd-11e8-a160-89cc2ad9e8e2.json
│       └── visualization:f8290060-4c88-11e8-b3d7-01146121b73d.json
```

## Invoking the create command
The create command read all the kibana configs present in the local `./config folder` and push them into Kiban cluster

```
./odfe-kibana-sync push --username admin --password admin --url https://localhost:9200
```

## Experimenting with the tool

You can run the opendistro ELK locally by running the following command

`docker-compose up `

More details on getting started with ELK Open distro can be found [here](https://opendistro.github.io/for-elasticsearch-docs/#get-started)






