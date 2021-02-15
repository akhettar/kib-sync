
![go](go.png)
# odfe-kibana-sync 
![Master](https://github.com/akhettar/odfe-kibana-sync/workflows/Master/badge.svg)

## Overview
This is a command line tool to sync and create the `Kibana object configs` from a given [open distro elk cluster](https://opendistro.github.io/for-elasticsearch-docs/) - see below the help section. This tool assumes that Kibana cluster holds the source of the truth in relation to the configuration file. The sole purpose of this tool is to be run periodically to sync these configuration files with the a given Git repository.

Idelly, this tool should be run in a CI pipeline for a given project that hosts Kibana configuration files - see example Github project.

The following kibana configuration files can be synched and created in a given kiban instance:
1. Monitors
2. Dashboards
3. Saved search
4. Destinations(slack, email etc)
5. Email accounts
6. Email groups.

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


## Invoking the sync command
The sync command fetches all the monitors config defined in the given kibana cluster and store them locally in the `./config folder`

```
./odfe-kibana-sync sync --username admin --password admin --url https://localhost:9200
```

The above command should produce a `config` folder containing the following structure (of course depending on the kiban objects configured in the cluster)
```
├── config
│   ├── dashboard
│   │   ├── dashboard:4b85e090-f4be-11ea-8342-bf90f7b9d26e.json
│   │   ├── dashboard:722b74f0-b882-11e8-a6d9-e546fe2bba5f.json
│   │   └── dashboard:edf84fe0-e1a0-11e7-b6d5-4dc382ef7f5b.json
│   ├── destination
│   │   └── XILEfXcBZWXOV7PGGr3r.json
│   ├── email_account
│   ├── email_group
│   ├── monitor
│   │   └── 13otoHcBbX-aeATowSlk.json
│   └── search
│       ├── search:3ba638e0-b894-11e8-a6d9-e546fe2bba5f.json
│       └── search:571aaf70-4c88-11e8-b3d7-01146121b73d.json
```

## Invoking the create command
The create command read all the kibana configs present in the local `./config folder` and push them into Kiban cluster

```
./odfe-kibana-sync push --username admin --password admin --url https://localhost:9200
```

## Experementing with the tool

You can run the opendistro ELK locally by running the following command

`docker-compose up `

More details on getting started with ELK Open distro can be found [here](https://opendistro.github.io/for-elasticsearch-docs/#get-started)






