# puppet changes
Scans changes for hosts and displays if they are recurring. Can also show a history of changes for all or single nodes.


For now it just checks daily/hourly/weekly and continious changes. And prints them.
Works on windows and linux. 

Can also display a list of changes ordered by time. Warnings/errors are filtered out but options can be given to also include them.
Do note when getting history for all your hosts it may take a while as all data needs to be retrieved and sorted.

# install
Get latest release from here https://github.com/Bryxxit/puppet-changes/releases

# usage
```
Scans changes for hosts and displays if they are recurring. Can also show a history of changes for all or single nodes.

Usage:
  puppet-changes [flags]

Flags:
  -C, --ca string       The ca certificate.
  -c, --cert string     The certificate.
      --config string   config file (default is $HOME/.puppet-changes.yaml)
  -h, --help            help for puppet-changes
  -r, --history         Show all changes by time.
  -H, --host string     The puppetdb host. (default "localhost")
  -k, --key string      The private key.
  -n, --node string     If you only want changes for a specific node/certname.
  -p, --port int        The puppetdb port. (default 8080)
  -E, --show-errors     Show the errors as well.
  -W, --show-warnings   Show the warnings as well.
  -t, --toggle          Help message for toggle
```
In order to use puppetdb with ssl you must supply key, ca and cert path.

## poll all
```
.\puppet-changes.exe

certname: node1 | message: notice /Stage[main]/something/Exec[something_install]/returns executed successfully (corrective) | pattern: continious
....
```
## poll one
```
.\puppet-changes.exe -n node1
certname: node1 | message: notice /Stage[main]/something/Exec[something_install]/returns executed successfully (corrective) | pattern: continious
...
```
## history all
```
.\puppet-changes.exe -r
2020-06-22 02:55:52.373 +0200 CEST certname: node1 message: notice /Stage[main]/something/Exec[something_install]/returns executed successfully (corrective)
...
```
you can display errors/warnigns by using
```
.\puppet-changes.exe -r -E -W
2020-06-22 02:55:52.373 +0200 CEST certname: node1 message: warning /Stage[main]/Test::Service/Service[test] Skipping because of failed dependencies
...
```


## history one
```
.\puppet-changes.exe -r -n node1
2020-06-22 02:55:52.373 +0200 CEST certname: node1 message: notice /Stage[main]/something/Exec[something_install]/returns executed successfully (corrective)
...
```