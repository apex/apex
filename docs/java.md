# Requirements

Apex expects pom.xml file to be in function directory. If not found, a generic pom.xml will be used.

Requirements for pom.xml to work with apex:
- mvn package produces valid jar
- jar.finalName prop have to set name jar name
- has AWS Lambda dependencies (com.amazonaws.aws-lambda-java-core)

# Generic pom.xml
When no pom.xml file in function main directory, Apex will generate one with default properties.

Generated pom.xml
```xml
 <project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
           xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/maven-v4_0_0.xsd">
     <modelVersion>4.0.0</modelVersion>

     <!-- Generic values -->
     <groupId>plugin-result</groupId>
     <artifactId>run.apex</artifactId>
     <packaging>jar</packaging>
     <version>apex</version>
     <name>apex-plugin-result</name>

     <properties>
         <!-- This prop is set when apex mvn package -->
         <jar.finalName>apex-plugin-result</jar.finalName>
     </properties>

     <dependencies>
         <dependency>
             <groupId>com.amazonaws</groupId>
             <artifactId>aws-lambda-java-core</artifactId>
             <version>1.1.0</version>
         </dependency>
     </dependencies>

     <build>
         <plugins>
             <plugin>
                 <groupId>org.apache.maven.plugins</groupId>
                 <artifactId>maven-shade-plugin</artifactId>
                 <version>2.3</version>
                 <configuration>
                     <createDependencyReducedPom>false</createDependencyReducedPom>
                 </configuration>
                 <executions>
                     <execution>
                         <phase>package</phase>
                         <goals>
                             <goal>shade</goal>
                         </goals>
                     </execution>
                 </executions>
             </plugin>
             <plugin>
                 <groupId>org.apache.maven.plugins</groupId>
                 <artifactId>maven-jar-plugin</artifactId>
                 <version>2.3.2</version>
             </plugin>
         </plugins>
     </build>
 </project>
```
