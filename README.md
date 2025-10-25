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
└── Close()
    └── Delete temporary file
Journal
├── Lock()
│   └── Obtain permission to interact with multiple files
├── Commit()
│   └── Commit changes to files and journal
└── Close()
    └── Delete temporary files and journal
Shield
└── GetJournal()
    └── Get journal of the current process
```

> How do I abort changes to file/journal?

Changes are only made if you commit. If you want to abort, just close the file/journal.