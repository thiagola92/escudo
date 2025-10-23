# escudo

## development

### rules
```
xxx_internal.go
├── Private functions
└── Use "assert.NoError()" functions
xxx.go
├── Public functions
├── Use "defer assert.Catch()"
└── Don't call public functions from the same file
    └── To avoid recursion
```

### main methods
```
File
├── Lock()
│   └── Obtain permission to interact with file
├── Commit()
│   └── Update original file but keep editing
└── Push()
    └── Update original file and close it
Journal
├── Lock()
│   └── Obtain permission to interact with multiple files
├── Commit()
│   └── Update all original files but keep editing
└── Push()
    └── Update all original files and close them
```