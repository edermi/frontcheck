# frontcheck
Checks if you can domain front a site, but curl doesn't fit your needs and you are too lazy to write something from scratch.

Assume you want to check if domain fronting still works with $provider.
You set up a service on your own infrastructure, get an account at $provider and configure your fronting.
Now, you want to find out under which circumstances you can fetch your own content while calling a foreign, trusted site also hosted at $provider.
The USP of this tool over `curl -H "Host: mysite" https://othersite.com` is (currently) that you have full control over SNI (have fun with curl).

This is a quick hack that tends to be more flexible to adapt to your needs compared to tools like curl.
Go's standard library allows you to to fancy things that are more difficult to achieve in other languages, e.g. because you have to deal with openssl there.

This is just a hacky small tool to ease my work, don't expect updates or support.
It may be a good starting point for your modifications, I'm open to contributions if they make sense to me.

## Example

Let's say I want to front using cloudflare and I'm "protecting" my site `demo.com` there, too.
There are tons of sites that also use cloudflare, for example `digitalocean.com`

~~~bash
$ ./frontcheck -url=digitalocean.com
Status: 200 OK
Proto: HTTP/1.1
Up to first 1000 body bytes:
<!DOCTYPE html><html lang="en"><head><script>
              !function(){var analytics=window.analytics=window.analytics||[];if(!analytics.initialize)if(analytics.invoked)window.console&&console.error&&console.error('Segment snippet included twice.');else{analytics.invoked=!0;analytics.methods=['trackSubmit','trackClick','trackLink','trackForm','pageview','identify','reset','group','track','ready','alias','debug','page','once','off','on','addSourceMiddleware','addIntegrationMiddleware','setAnonymousId','addDestinationMiddleware'];analytics.factory=function(e){return function(){var t=Array.prototype.slice.call(arguments);t.unshift(e);analytics.push(t);return analytics}};for(var e=0;e<analytics.methods.length;e++){var key=analytics.methods[e];analytics[key]=analytics.factory(key)}analytics.load=function(key,e){var t=document.createElement('script');t.type='text/javascript';t.async=!0;t.src='https://cdn.segment.com/analytics.js/v1/' + key + '/analytics.min.js';var n=document.getElementsByT
~~~

There's some header data and the first 1000 bytes of the response body. 
Now, using my own site `demo.com`, which just returns a unfriendly nginx error page:

~~~bash
$ ./frontcheck -url=digitalocean.com -sni=demo.com -host=demo.com
Status: 200 OK
Proto: HTTP/1.1
Up to first 1000 body bytes:
<!DOCTYPE html>
<html>
<head>
<title>Go away!</title>
<style>
    body {
        width: 35em;
        margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }
</style>
</head>
<body>
<h1>This is a private page!</h1>
<p>If you see this page, you should strive on because there is nothing to see here.</p>

</body>
</html>
~~~

So, allegedly my site can be fronted via sites also running behind cloudflare if I just set the corret `Host`-Header and SNI name.
Without SNI, the following will happen (same as `curl -H "Host: demo.com" https://digitalocean.com`):

~~~bash
$ ./frontcheck -url=digitalocean.com -host=demo.com
Status: 403 Forbidden
Proto: HTTP/1.1
Up to first 1000 body bytes:
<html>
<head><title>403 Forbidden</title></head>
<body>
<center><h1>403 Forbidden</h1></center>
<hr><center>cloudflare</center>
</body>
</html>
~~~

