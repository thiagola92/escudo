# escudo

## development

### rules
```
xxx_internal.go
├── Private functions
├── Use "defer assert.Catch()"
└── Use "assert.NoError()" functions
xxx.go
├── Public functions
└── Don't call public functions from the same file
    └── To avoid recursion
```

### main methods
```
File
├── Lock()
│   └── Obtain permission to interact with file
├── Commit()
│   └── Commit changes to file
└── Push()
    └── Commit changes to file, then close it
Journal
├── Lock()
│   └── Obtain permission to interact with multiple files
├── Commit()
│   └── Commit changes to files and journal
└── Push()
    └── Commit changes to files and journal, then close them
```

> If I want to abort changing the file?

Just close it. If you don't commit or push changes, nothing will happen.