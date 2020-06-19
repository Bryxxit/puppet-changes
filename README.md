# puppet changes
Puppet changes is a small cli that polls puppetdb to look for recurring changes.

For now it just checks daily/hourly/weekly and continious changes. And prints them.
Works on windows and linux. 

# install
Get latest release from here https://github.com/Bryxxit/puppet-changes/releases

# usage
```
.\puppet-changes.exe -h
Scans changes for hosts and displays if they are recurring.

Usage:
  puppet-changes [flags]

Flags:
  -C, --ca string       The ca certificate.
  -c, --cert string     The certificate.
      --config string   config file (default is $HOME/.puppet-changes.yaml)
  -h, --help            help for puppet-changes
  -H, --host string     The puppetdb host. (default "localhost")
  -k, --key string      The private key.
  -n, --node string     If you only want changes for a specific node/certname.
  -p, --port int        The puppetdb port. (default 8080)
  -t, --toggle          Help message for toggle
```
In order to use puppetdb with ssl you must supply key, ca and cert path.

poll all
```
.\puppet-changes.exe

certname: node1 | message: notice /Stage[main]/something/Exec[something_install]/returns executed successfully (corrective) | pattern: continious
....
```
poll one
```
.\puppet-changes.exe -n node1
certname: node1 | message: notice /Stage[main]/something/Exec[something_install]/returns executed successfully (corrective) | pattern: continious
...
```
