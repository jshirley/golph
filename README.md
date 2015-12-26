# Golph: Go Library for Phabricator

_Yes, I hate the ph too!_

This library is heavily inspired by the excellent Digital Ocean library [`godo`](https://github.com/digitalocean/godo).

## Go 1.2+

This is tested from Go 1.2 - 1.5 (and tip).

## Usage

```go
import "github.com/jshirley/golph"

client := golph.NewClient("api-token", "https://phabricator.example.com")
```

## Examples

### Listing Users

For pagination, you can iterate through while there are additional pages.

```go
func UserList(client *golph.Client) ([]golph.User, error) {
    // create a list to hold our droplets
    list := []golph.User{}

    // create options. initially, these will be blank
    opt := &golph.ListOptions{}
    for {
        users, resp, err := client.Users.List(opt)
        if err != nil {
            return nil, err
        }

        // append the current page's users to our list
        for _, d := range users {
            list = append(list, d)
        }

        // if we are at the last page, break out the for loop
        if resp.Links == nil || resp.Links.IsLastPage() {
            break
        }

        page, err := resp.Links.CurrentPage()
        if err != nil {
            return nil, err
        }

        // set the page we want for the next request
        opt.Page = page + 1
    }

    return list, nil
}
```

# Contributing

Help me make this library awesome! Please see the [contributing guidelines](./CONTRIBUTING.md).
