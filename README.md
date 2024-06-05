# Jrnl

A tiny microblogging website made to try out htmx and PicoCSS.
Backend written in go, with charmbracelet/log for logging :)

This was just a tiny weekend project i did for myself, nothing revolutionary.

Features:
- Post creation
- Tag filtering
- Post ordering
- SPA ( kinda )
- ~~Error handling~~ hahahaha no

# Installation / Running instructions

Just install go and run 
```
$ go run main.go
```

That's it. it's pretty self contained.

# Thoughts

There are some parts of this site that i would not normally write the way i have written them here, such as this site being an SPA when it really doesn't need to be. 
It's an SPA because i wanted to try out htmx's SPA features, but it being an SPA leads to a good bit of jank, like for example having to define html blobs that get replaced when you click on a link rather than just directly sending full html pages.
This also ended up making me having to define a new handler in the backend for every path as all of them functioned *just slightly* differently, rather than just defining one big handler for "/*.html" path that just delivers the entire page as it is.

Another thing that i don't like about this is the way i handled database serialization. Currently the server just overrides the database file ( a json file. it works alright as a simple data store \*shrug\* ) and then just saves the new data in it's place. Ideally i'd like for it to actually append the data to the file rather than overwriting it everytime, as that would not scale.

Also initially the "sort posts" handler did a weird thing where i actually had two handler functions, "/oldest" and "/recent", but both of them pinged the same handler.
The difference was that that handler also took an `s` parameter, which was a boolean and controlled whether or not to sort the output list or not.
The handler calls used closures and looked something like `mux.HandleFunc("/recent", func(res ResponseWriter, req *Request){sortHandler(true, res, req)})`, with true for `/recent` and `false` for `/oldest`.
It was terrible but i kinda liked it. I used it until i remembered URL queries were a thing and then replaced it.

Oh also this is missing error handling, very important but my development cycle did not need proper error handling as there wasn't much that would fail so i didn't bother implementing that.
However that is no excuse, error handling is important and you 🫵 should never skip it.
