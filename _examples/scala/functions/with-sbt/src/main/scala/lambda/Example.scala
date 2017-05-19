package lambda

import com.amazonaws.services.lambda.runtime.Context


class Example {

  	case class ExampleRequest(hello: String)

  	case class ExampleResponse(hello: String)

  	def handler(event: ExampleRequest, context: Context) = ExampleResponse(event.hello)

}
