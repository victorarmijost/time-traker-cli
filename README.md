# Time Tracker CLI

This tool allows to insert time records, into the BairesDev Time Tracker; the login process works on top of the Time Tracker web application; to use the application you will need to login with your Google account, that process is simulated using the Go [Chromedep](https://pkg.go.dev/github.com/chromedp/chromedp#section-readme) package, which required that you have ***Google Chrome installed***.

To compile the application run:
```
cd cmd
go build -o tt *.go
```

It will create a binary called `tt`, use this to execute it:
```
./tt
```

The first time that you open the application, it will create its required folders structure and it will ask you for your some configuration information like your focal point and your project.

The tool works with a REPL (Read-Eval-Print loop), where you have to send a command to perform an action, you can see the list of available commands by using:

```
tt > help
```

To use one of the commands, just send the command name, and the tool will requested you the needed arguments.

Example:

```
tt > rec

- #: bug

- Comment: issue with accounts

**** Record started! ****
```

The tool also allows the following syntax (which we call a short version):

`{command name}; {arg 1}; {arg 2}; {arg 3}...`

Example:
```
tt > rec; bug; issue with accounts

**** Record started! ****
```

The provided arguments will be filled in the required command fields.

Most of the commands works by requesting a set of fields and providing a sigle response after all the fields are provided, we call them *action commands*. Even though, there are some special commands the doesn't follow this pattern, we call them *interactive commands* (i.e. `temp add`)

Some commands like: `rec` or `add` are using templates; a template is a way to predefine the fields for some records,you can create a new template by using `temp add`.

There are some templates already created under the `cmd/templates/rec` folder, but you can delete them and add your own templates. Use `temp list` to see all the created templates.

Workflow
==

This is the workflow which was used to design the tool:

1. Use `rec` to start a time record. It will request a template (you must have some created already); and, if it is required, a description. It will record this information along with the initial time.
2. Use `view` to see your current working task.
3. Use `end` to complete your task, it will use the current time and the recored start time to calculate the worked time.
4. Repeate the process for all you day tasks.
5. Use `list` to see all your recorded time.
6. Use `commit` to send your records to the BairesDev Time Tracker.

Edge cases
==

1. If you missed to record some time, you can use `add` to manually add it, it will request a template, an if it is required, a description and a task duration (in hours).
2. If you missed the start time of a task that you are currently working, use `rec at` to start it at a defined hour.
3. If you missed the end time of a task that you are currently working, use `end at` it will calculate the time base on the end hour that you provide.
4. If started the time record with the wrong template, use `edit` to change the record information, keeping the same start time.

The tool works on today's date by default; but you can use `change date` to change to another date. This can be useful to add missing records for a previous day.


Status bar
==

When you are using the application, the prompt may change adding some status values:
```
[Worked:0.25][Commited:9.00][Tracking:0.25] tt >
```

1. ***Worked***: Is the total time that you have worked during the day and that is not commited. Once you commit this time will be added to Commited.
2. ***Commited***: Is the total time that you have worked during the day that is commited.
3. ***Tracking***: is the time worked on the current task. Once you end the task this time will be added to Worked.

Application folders structure
==

The first time you open the application, it will create it own folder structure. This are the relevant files under those folders:
1. ***.tmp***: it stores temporary files like the application state and a cache.
2. ***local***: it a folders base local db, to store the time records before commiting them, the records are store here on `.json` format.
3. ***templates***: commands like `rec` or `add` uses templates (see above); those templates are saved here on `.json` format. The template fields that starts with "x-" are information fields, which means that they don't store values required by the template.

Besides that, you will see a `config.json` file, which stores your email, project, focal point, and working time information, all this information is required when you first run the application. If this file is removed, all the information will be required again when run the application.
