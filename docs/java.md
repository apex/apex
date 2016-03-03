# Requirements

Apex expects `pom.xml` file to be in function directory. Requirements for `pom.xml` to work with apex:

- `mvn package` produces valid jar
- `jar.finalName` prop has to set name jar name
- has AWS Lambda dependencies (`com.amazonaws.aws-lambda-java-core`)
