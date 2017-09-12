import sys
import asyncio

async def test_async_generator():
    for i in range(5):
        yield i
        await asyncio.sleep(0.01)

def handle(event, context):
    name = "Fred"
    formatted = f"He said his name is {name}."
    primes: List[int] = [2, 3, 5, 7]
    one_thousand = 1_000
    one_million_underscore = '{:_}'.format(1000000)

    return {
        'event': event,
        'formatted': formatted,
        'primes': primes,
        'one_thousand': one_thousand,
        'one_million_underscore': one_million_underscore,
        'sys_version': sys.version,
        'sys_version_info': sys.version_info,
    }
