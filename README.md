mkp-kibana-monitors

# Overview
This is a command line tool to sync and update the `Kibana object configs` - see below the help section. This tool is run as part of the scheduled pipeline to run every hour. It downloads the monitors config from kiban cluster
then push the latest configuration into Gitlab - this project. The kibana objects supported are:
1. Monitors
2. Dashboards
3. Saved search
4. Destinations(slack, email etc)

```./kibana-sync -h
This tool performs the followings:
                1. Fetches configured monitors for the given kibana cluster and store them locally as json files.
                2. Pushes the changes done to the monitor's config to Kiban cluster

Usage:
  kibana-sync [flags]
  kibana-sync [command]

Available Commands:
  help        Help about any command
  sync        Fetch all configured monitors from Kibana cluster
  push      pushes the monitor's config to Kibana cluster

Flags:
  -h, --help              help for kibana-sync
      --password string   The kibana cluster password. This is a required argument to connect to the ELK cluster
      --url string        The kibana cluster url. This is required argument to connect to the ELk cluster
      --username string   The kibana cluster username. This is required argument to connect to the ELK cluster

Use "kibana-sync [command] --help" for more information about a command.
```

## The pipeline

The pipeline of this project is scheduled to run every hour, it's not triggered by push to master. So any changes to Kibana monitor config will get picked up every hour or so.


## Limitations

The limitations of this tool are highlighted below. Hopefully, subsequent releases of this tool will address some of them

1. The push command push all the configuration files present in the config folder regardless if there is a change or no. This is not an issue at all, but we can improve.

## Invoking the sync command
The sync command fetches all the monitors config defined in the given kibana cluster and store them locally in the `./config folder`

```
./kibana-sync sync --username VF_Kibana_EMEA --password Vfkibana***** --url https://vpc-vf-sysint-emea-es-ir-jiqfwlmrnmrkm7ydovqpsvqueu.eu-west-1.es.amazonaws.com
```

## Invoking the push command
The push command read all the monitor configs present in the local `./config folder` and push them into Kiban cluster

```
./kibana-sync push --username VF_Kibana_EMEA --password Vfkibana****** --url https://vpc-vf-sysint-emea-es-ir-jiqfwlmrnmrkm7ydovqpsvqueu.eu-west-1.es.amazonaws.com
```




