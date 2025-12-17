I would like to create a Go application http server that responds to different endpoints in different ways based on the configuration.

Only use Go default imports if at all possible. Ask me about adding additional libraries.

The application is intended to be run from a container. 

The application has the following command line arguments.
command line arguments are either --arg or -a with the shortened form taken from the first letter of the full argument.
Take the first n letters if there is an overlap to make them unique.

* help
provide a list of all commands and arguments

* file (required)
specify a configuration file.

* port (optional)
the port that the server will listen to with a default of 8888

* wonky (optional)
a value between 1 and 100 representing the percentage likelihood that any request will randomly get one of: error (500), delay of 5 seconds, or slow (405). Default is 0 (disabled).

The application will log all messages to standard out.

# input file format
The input file format will be a JSON object consisting of an array of objects.

All of the elements in the object in the array are required except for headers which is optional.

An empty array is invalid and will error on startup.
A missing file is invalid and will error on startup. 

{
    [
        {"verb": "GET", "url": "/foo", "code":"200", "response":"{}","headers":["application/json]}
    ]
}

The application will receive http requests on the specified or default port and if the verb and url match then return the code, response, and headers as specified.

If the application receives a verb and url request that doesn't match then return a 404 response code.

The requesting url can have url parameters that change the response behavior.

* error
if the url matches then return a 500 instead of the specified resposne code

* slow
if the url matches then return a 405 to fast response code

* delay=amout unit
if delay is specified then delay the response by the amount specified.
the first value is the number and the second is the unit
100m means 100 milliseconds
10s means 10 seconds
1M means 1 minute.


# Create unit tests
# Creat test examples using curl


Create a Github action to build the application into a container and publish as a GitHub artifact
Make sure that the container is created for ARM64 as well as X86 targets.

Include instructions on how to pull and execute the container. 