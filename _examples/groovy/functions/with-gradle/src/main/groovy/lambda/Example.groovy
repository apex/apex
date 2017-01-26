package lambda

import com.amazonaws.services.lambda.runtime.Context

class Example {

    static class ExampleRequest {
        def hello
    }

    static class ExampleResponse {
        def hello
    }

    ExampleResponse handler(ExampleRequest event, Context context) {
        new ExampleResponse(hello: event.hello)
    }
}
