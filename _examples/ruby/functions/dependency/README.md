The build hook can be used to install required libraries using bundle. You need to install these libraries at function level so that they become part of the zip files.

The requirements are mentioned in the Gemfile and the build hook is used to install the dependencies before the build

```
"hooks":{
    "build": "bundle install && bundle install --deployment"
  }
```

