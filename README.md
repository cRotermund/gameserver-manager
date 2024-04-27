# Game Server Manager

A simple set of libraries curated to manage my gameserver in AWS via a lambda gateway of my own design.
The implementation contains a command line application and a very simple discord bot.


## Local Development

Local development relies on a virtual environment to share the module code.  Poetry is used
in the CI/CD platform.

To get set up writing and running the code locally:
* Create a virtual environment in the root directory of the project: `.venv`
* Activate your virtual environment
* Ensure your IDE is configured to your virtual environment
* To share the library code and actively develop, navigate /src/libs/gsmclient and run a `pip install -e .`
* You can now develop the library and use it in the other projects
* Run `pip install .` in any of the other projects to build and pull in dependencies.

To add a dependency:
* Navigate to the root directory of the project you wish to add the dependency (e.g. `src/libs/gsmclient`)
* Execute `poetry add <DEPENDENCY>` (i.e. `poetry add requests`)

## Deployment

### Daemon
* GitHub action builds image pushes to git hub packages
* EC2 server pulls, runs in swarm

## License

[MIT](https://choosealicense.com/licenses/mit/)