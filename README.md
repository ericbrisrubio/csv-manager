# CSV-MANAGER
-Modify huge csv headers and get basic info from it

-This App assumes the headers exist in the first line of the csv

##Installation
-If you have golang then run:

1. `go get github.com/ericbrisrubio/csvmanager` 
2. ``cd $GOPATH/src/github.com/ericbrisrubio/csvmanager/``
3. `go install`

Once you run go install you have access globally to the binary file


-If you DON'T have golang installed(you need git installed in this case):
1. `git clone http://github.com/ericbrisrubio/csvmanager.git`
2. `cd csvmanager/`

##Commmands
-To get the fileinfo:

no golang installed:
`csvmanager -f "path/to/file" info`

golang installed:
`./csvmanager -f "path/to/file" info`

-To modify the file headers:

no golang installed:
`./csvmanager -c "path/to/configfile" -f "path/to/file" d`

golang installed:
`csvmanager -c "path/to/configfile" -f "path/to/file" d`

##Current features
The headers modification currently changes the headers to uppercase only.
(To have this value just add "Upper" value for key entry "name_changer_instance" in the config.json 
file)

##Extending the headers transformation (Developers note)
For different ways to extend the way the headers are modified
you should extend `IChanger` interface and add a new entry under:
```
var FactoryNameChanger = map[string]fn{
	"Upper": func() IChanger{return UpperChanger{}},
}
``` 
The new entry key could be any picked named and the value should be
a function returning an instance of the new defined class.

To use that new way to change the headers the entry key defined should
be referred in the config.json file as the value of `name_changer_instance`
key


