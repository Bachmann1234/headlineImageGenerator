# NYT Headline Image Generator

This basic script will call out to NYT to get a random popular headline.

It will then ask DALLE*2 to draw that headline.

If DALLE*2 rejects this headline it will ask GBT-3 to rewrite the headline to make it PG rated and remove names. 

Then it will draw the new headline