# Requirements

Apex expects that your JVM project generate a fat jar, with all dependencies included. Requirements for your jar to work with Apex:

- `apex.jar` has to be in either `target/` or `build/libs/`
- has AWS Lambda dependencies (`com.amazonaws.aws-lambda-java-core`)
