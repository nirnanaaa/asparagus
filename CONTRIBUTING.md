Contributing to Asparagus
========================

Bug reports
---------------
Before you file an issue, please search existing issues in case it has already been filed, or perhaps even fixed. If you file an issue, please include the following.
* Full details of your operating system (or distribution) e.g. 64-bit Ubuntu 14.04.
* The version of Asparagus you are running (or commit)
* Whether you installed it using a pre-built package, or built it from source.
* A small test case, if applicable, that demonstrates the issues.

Remember the golden rule of bug reports: **The easier you make it for us to reproduce the problem, the faster it will get fixed.**
If you have never written a bug report before, or if you want to brush up on your bug reporting skills, we recommend reading [Simon Tatham's essay "How to Report Bugs Effectively."](http://www.chiark.greenend.org.uk/~sgtatham/bugs.html)


Feature requests
---------------
We really like to receive feature requests, as it helps us prioritize our work. Please be clear about your requirements, as incomplete feature requests may simply be closed if we don't understand what you would like to see added to Asparagus.

Contributing to the source code
---------------

Asparagus follows standard Go project structure. This means that all your Go development are done in `$GOPATH/src`. GOPATH can be any directory under which Asparagus and all its dependencies will be cloned. For full details on the project structure, follow along below.

Submitting a pull request
------------
To submit a pull request you should fork the Asparagus repository, and make your change on a feature branch of your fork. Then generate a pull request from your branch against *master* of the Asparagus repository. Include in your pull request details of your change -- the why *and* the how -- as well as the testing your performed. Also, be sure to run the test suite with your change in place. Changes that cause tests to fail cannot be merged.

There will usually be some back and forth as we finalize the change, but once that completes it may be merged.

To assist in review for the PR, please add the following to your pull request comment:

```md
- [ ] CHANGELOG.md updated
- [ ] Rebased/mergable
- [ ] Tests pass
```

Installing Go
-------------
Asparagus requires Go 1.9.2.

We find gvm, a Go version manager, useful for installing Go. For instructions
on how to install it see [the gvm page on github](https://github.com/moovweb/gvm).

After installing gvm you can install and set the default go version by
running the following:

    gvm install go1.9.2
    gvm use go1.9.2 --default

Installing godep
-------------
We use [godep](https://github.com/tools/godep) to manage dependencies.  Install it by running the following:

    go get github.com/tools/godep

Revision Control Systems
-------------
Go has the ability to import remote packages via revision control systems with the `go get` command.  To ensure that you can retrieve any remote package, be sure to install the following rcs software to your system.
Currently the project only depends on `git` and `mercurial`.

* [Install Git](http://git-scm.com/book/en/Getting-Started-Installing-Git)
* [Install Mercurial](http://mercurial.selenic.com/wiki/Download)


Getting the source
------
Setup the project structure and fetch the repo like so:

```bash
    mkdir $HOME/gocodez
    export GOPATH=$HOME/gocodez
    go get github.com/nirnanaaa/asparagus
```

You can add the line `export GOPATH=$HOME/gocodez` to your bash/zsh file to be set for every shell instead of having to manually run it everytime.


Cloning a fork
-------------
If you wish to work with fork of Asparagus, your own fork for example, you must still follow the directory structure above. But instead of cloning the main repo, instead clone your fork. Follow the steps below to work with a fork:

```bash
    export GOPATH=$HOME/gocodez
    mkdir -p $GOPATH/src/github.com/myname
    cd $GOPATH/src/github.com/myname
    git clone git@github.com:<myname>/asparagus
```

Retaining the directory structure `$GOPATH/src/github.com/myname` is necessary so that Go imports work correctly.


Build and Test
-----

Make sure you have Go installed and the project structure as shown above. To then get the dependencies for the project, execute the following commands:

```bash
cd $GOPATH/src/github.com/nirnanaaa/asparagus
godep restore
```

To then build and install the binaries, run the following command.
```bash
./build.py
```
The binary will be located in `./build`.

To run the tests, execute the following command:

```bash
cd $GOPATH/src/github.com/nirnanaaa/asparagus
go test -v ./...
```


Continuous Integration testing
-----
Asparagus uses CircleCI for continuous integration testing. To see how the code is built and tested, check out [this file](./scripts/test-reporter.sh). It closely follows the build and test process outlined above. You can see the exact version of Go Asparagus uses for testing by consulting that file.
