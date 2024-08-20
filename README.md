# Go Tour | Web Crawler

This repository contains some different solutions to the `web crawler` final exercise in the go tour [Exercise: Web Crawler](https://go.dev/tour/concurrency/10)

## Background

I like to think of every possible way/syntax/technique to accomplish a solution to a coding problem, so I thought to share some. There are obviously infinite amount of ways to solve this exercise.

The initial solutions I came up with are in inside [sol1/](./sol1/) directory. They were kinda based off the `equivelent binary tress` exercise (from go tour as well). Every other solution is from the community (google searches), with slight modifications.

## Realizations

After looking at other people's solutions, I realized the following points:

- You should check for the _depth_ being `0` before spawning a child goroutine (or recursing), so the program can be more effecient.
- You should check for the _url_ being already _cached_ before spawning a child goroutine (or recursing), so the program can be more effecient.
- You can use a single commuincation _object_ that is _passed_ down to sub goroutines.
- The initial _url_ will never be _cached_, so you might base your code structure on that assumption.
- There are so many _communication_ methods and _combinations_ that you can use to achieve the same result.

## The Best Solution

For me, none of these are the **best**. Obviously, it all comes down to **preference**. However, out of these solutions, I would choose [sol8/main_v2.go](./sol8/main_v2.go) as the best one for its clean data flow and flexibility ([sol3/main_v2.go](./sol3/main_v2.go) and [sol6/main.go](./sol6/main.go) are also pretty nice).

If I had refactored my initial parallel solutions ([sol1/main_v2.go](./sol1/main_v2.go), [sol1/main_v3.go](./sol1/main_v3.go)) arccording to the first three points from the `Realization` section, I would choose it as the best solution for its simplicity and readability.

__NOTE__: I have just added the best solutions (in my opinion). If you strive for effeciency and simplicity, look at [best/main_v1.go](./best/main_v1.go). If you strive for flexibility and OOP style of code, look at [best/main_v4.go](./best/main_v4.go).