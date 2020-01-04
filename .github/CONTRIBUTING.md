# How to contribute

If you would like to contribute code to this project, fork the repository and send a pull request.

## Prerequisite

In this project, use `go mod` as the package management tool and make sure your Go version is higher then `Go 1.13`.

## Fork

Before contributing, you need to fork [faas-akash](https://github.com/vitwit/faas-akash) to your GitHub account.

## Contribution flow

```bash
$ git remote add faas-akash https://github.com/vitwit/faas-akash.git
# sync with the remote master
$ git checkout master
$ git fetch faas-akash
$ git rebase faas-akash/master
$ git push origin master
# create a PR branch
$ git checkout -b your_branch   
# do something
$ git add [your change files]
$ git commit -sm "xxx"
$ git push origin your_branch
```

## Configure Jetbrains - GoLand

`faas-akash` uses `go mod` to manage dependencies, so make sure your IDE enables `Go Modules(vgo)`.

To configure annotation processing in GoLand, follow the steps below.

1. To open the **Go Modules Settings** window, in GoLand, click **Preferences** > **Go** > **Go Modules(vgo)**.

2. Select the **Enable Go Modules(vgo) integration** checkbox.

3. Click **Apply** and **OK**.

## Code style

We use Go Community Style Guide in `faas-akash`. 
For more information, see [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).
Always use gofmt and linter (linter config is provided), before submitting any code changes.

To make your pull request easy to review, maintain and develop, follow this style.

## Test Scripts
To make sure that the changes/features/bug-fixes submitted by you does not break anything, always include test scripts 
for your work.

## Update dependencies

`faas-akash` uses [Go 1.13 module](https://github.com/golang/go/wiki/Modules) to manage dependencies.
To add or update a dependency, use the `go mod edit` command to change the dependencies
