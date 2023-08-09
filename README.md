This code fetches 100 integers using the endpoint `/integers/:int`.
It will report a list of which integers are even and which are odd. 

For each integer fetched, the server will send back the value as JSON. For example:
`{ 
    "value": 1
}`

To avoid triggering a rate limit, the code only makes 10 concurrent requests at a time.

This repo contains coding mistakes and overall bad form. How would you improve it as a whole to make it ready for production? Feel free to rewrite as much as you’d like!

################ New stuff ################
- Introduced a `DOMAIN` environment variable, which defaults to `localhost:8080` if it is not set
- Response should be JSON of the form
```
{
    "even":[2,4,6,8,etc.],
    "odd":[1,3,5,7,etc.]
}
```
- a list of `errors` will be appended to the payload if any were found
