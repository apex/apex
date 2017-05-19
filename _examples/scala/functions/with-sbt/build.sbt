
lazy val root = (project in file("."))
	.settings(
    	name := "apex-scala-example",
    	version := "0.1.0-SNAPSHOT",
    	scalaVersion := "2.12.2",
    	retrieveManaged := true,
    	libraryDependencies += "com.amazonaws" % "aws-lambda-java-core" % "1.1.0",
    	libraryDependencies += "com.amazonaws" % "aws-lambda-java-events" % "1.3.0",
    	cleanFiles <+= baseDirectory { base => base / "build" }
  	)	

assemblyMergeStrategy in assembly := {
  	{
    	case PathList("META-INF", xs@_*) => MergeStrategy.discard
    	case x => MergeStrategy.first
  	}
}

assemblyOutputPath in assembly := baseDirectory.value / "build" / "libs" / "apex.jar"
