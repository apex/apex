The build hook can be used to install required libraries using pip. You need to install these libraries at function level so that they become part of the zip files.

The requirements are mentioned in the requirements.txt and the build hook is used to install the dependencies before the build

```
"hooks":{
    "build": "pip install -r requirements.txt -t ."
  }
```

