# Welcome to lil TUI: A Simple and Flexible Interface for List and Get Operations

Welcome to lil TUI, a lightweight terminal user interface (TUI) application designed to provide a seamless, generic interface for any service that supports list and get operations.

## Features
User-friendly terminal interface
Easily configurable for a variety of services
Simple plug-and-play setup


## Demo
![Demo](https://github.com/andrievsky/lil/raw/master/docs/demo.gif "Demo")

## Configuration
To connect a new service to lil, follow these steps:

Create a hidden directory in your home folder called .lil.
Inside the .lil directory, create a new folder with a name that corresponds to your desired service. This folder should contain three executable files:

```
~/
  .lil/
    your_service/
      init
      list
      get
```

These files serve the following purposes:

- init: Runs once, returns the base path, and optionally initializes the service.
- list: Accepts a path as a parameter and returns a list of traversable and final items (similar to folders and files). By convention, traversable items should end with a /.
- get: Accepts a path as a parameter and returns the content of a final item.

### Acknowledgments

We would like to express our gratitude to the Tcell library (https://github.com/gdamore/tcell) for providing a fantastic minimalistic TUI foundation. We built a few minimalistic custom components because we were unaware of prebuilt options, which actually makes the application even simpler.