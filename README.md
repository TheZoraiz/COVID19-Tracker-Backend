# COVID 19 Tracker Backup Server

A backup server deployed at heroku.com to serve data to the front end in case the official API throws an error.

The saveData.go and data.go are to be used locally, for creating the new backup data to push to heroku server. The server/ directory contains the code for handling requests from front end and serving the required data. All data is stored as JSON in txt files.

You can use the tracker directly from this link:
https://TheZoraiz.github.io/COVID19-Tracker

The code for the front end is available here: https://github.com/TheZoraiz/COVID19-Tracker
