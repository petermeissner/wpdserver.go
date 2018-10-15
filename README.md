

# WPDAccessServer

Simple Go server that routes to simple SQL queries that are executed on page views database server and pushs back results to user as JSON. 



## Routes


### Page Views for Articles

**[path:]** `http://SERVER:8880/search/{lang}/{articlename}`

**[lang:]** one of 20 wikipedia language shorthands: cs, da, de, el, en, es, et, fi, fr, hu, it, no, pl, pt, ru, sk, sl, sv, tr, zh

**[articlename:]** Name of a Wikipedia article in lower case latters as it is found in the URL of an Wikipedia article.



### Page Views for Article Searched by Regex

**[path:]** `http://SERVER:8880/search/{lang}/{regex}`

**[lang:]** one of 20 wikipedia language shorthands: cs, da, de, el, en, es, et, fi, fr, hu, it, no, pl, pt, ru, sk, sl, sv, tr, zh

**[regex:]** a Postgres compliant regular expression as described here: https://www.postgresql.org/docs/current/static/functions-matching.html




### Article Name Search via Regex

**[path:]** `http://SERVER:8880/search/{lang}/{regex}`

**[lang:]** one of 20 wikipedia language shorthands: cs, da, de, el, en, es, et, fi, fr, hu, it, no, pl, pt, ru, sk, sl, sv, tr, zh

**[regex:]** a Postgres compliant regular expression as described here: https://www.postgresql.org/docs/current/static/functions-matching.html
