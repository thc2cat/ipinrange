# ipinrange

Finds network addresses in the standard input, and filters the output based on their belonging to a network block

## Context

When searching for ip source of trouble, i needed an ip extractor ( see also extractip ), that would not only print ip, but check that these are not private, reserved or local ips.

## Usage

```shell
$ go mod tidy & &go build

$ rg -z someone mtaauth.log.gz  | ipinrange -n local 

Will print only lines containing ip addresses not matching my networks      
                     
```
