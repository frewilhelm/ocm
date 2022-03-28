## ocm componentarchive transfer

transfer component archive to some component repository

### Synopsis


Transfer a component archive to some component repository. This might
be a CTF Archive or a regular repository.
Explicitly supported types, so far: OCIRegistry, CTF (directory, tar, tgz).
If the type CTF is specified the target must already exist, if CTF flavor
is specified it will be created if it does not exist.

Besides those explicitly known types a complete repository spec might be configured,
either via inline argument or command configuration file and name.


```
ocm componentarchive transfer [<options>]  <source> <target> [flags]
```

### Options

```
  -h, --help   help for transfer
```

### SEE ALSO

* [ocm componentarchive](ocm_componentarchive.md)	 - 
