cantuse is like caniuse.com, except that it lists the browsers in which a
feature *won't* work, instead of listing in which it will. I find this a more
useful way to look at the data.

Right now this is just a basic commandline utility. A web version will be added
*soon™*.

Install it with `go install zgo.at/cantuse@latest`, which will put the binary in
`~/go/bin`.

It uses the same data as caniuse.com (https://github.com/Fyrd/caniuse); you'll
need to fetch the `data.json` yourself:

    curl https://raw.githubusercontent.com/Fyrd/caniuse/master/data.json > data.json

Use just `cantuse` to display all data, or `cantuse arrow` to search by feature
name. Use `-ignore` to ignore some browsers and always consider them
"supported"; this accepts a comma-separated list of the browser + version as it
appears in the output; for example:

    $ cantuse -ignore 'IE 11'
    $ cantuse -ignore 'IE 11,Opera Mini all'

Browsers with a usage lower than 0.05% are not displayed individually to reduce
noise (chances are you don't care about 0.009% of people reported as using
IE5.5, or 0.004% usage of Firefox 2). It's still counted in the total supported
number though.

Partial support is counted as "supported"; see caniuse.com for more detailed
support notes.
