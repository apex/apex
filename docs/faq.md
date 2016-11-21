
## How do you manage multiple environments?

It's highly recommended to create separate AWS accounts for staging and production environments. This provides complete isolation of resources, and allows you to easily provide granular access to environments as required. See [AWS credentials](#aws-credentials) for supplying an account profile.

AWS IAM roles can be used to provide quick access to each environment using a drop-down in the AWS Console.

## Can I test functions locally?

Currently there is no way to run functions locally, it would be a very large task to emulate the AWS. We recommend writing the bulk of your logic as libraries or packages native to your chosen language, using only thin connective layers in the Lambda functions themselves. This makes it easy to unit-test your functions, and makes them more portable if you're worried about vendor lock-in.

## How is this different than Serverless?

Serverless uses CloudFormation to bootstrap resources, which can be great for getting started, but is generally less robust than [Terraform](https://www.terraform.io/) for managing infrastructure throughout its lifetime. For this reason Apex does not currently provide resource management. This may change in the future for light bootstrapping, likely in an optional form.

At the time of writing Serverless does not support shimming for languages which are not supported natively by Lambda, such as Golang. Apex does this for you out of the box.

The structures imposed by each project are different, as well as varying features, see the documentation for each project to see what either supports.

Serverless aims to be provider agnostic, which can be both a pro and a con depending on the level of abstraction you're comfortable with, and if you desire to have a tool modelled closer to a single provider's capabilities. This similar to the contention around ORM vs "raw" queries.

Serverless is written using Node.js, Apex is written in Go. Apex aims to be a simple and robust solution, while Serverless intends on providing a more feature-rich solution, pick your poison.

## Is using the Node.js shim slow?

The shim creates a child process, thus creates a few-hundred millisecond delay for the first invocation. Subsequent calls to a function are likely to hit an active container, unless the function sees very little traffic.

## Do shimmed languages have any limitations?

The shim currently operates using JSON over stdout, because of this you must use stderr for logging.

## Can I manage functions that Apex did not create?

Apex is not designed to handle legacy cases, or functions created by other software. For most operations Apex references the local functions, if it is not defined locally, then you cannot operate against it. This is by-design and is unlikely to change.
