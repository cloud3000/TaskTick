# TaskTick

Three circular task ques, to run tasks based on system ticks, not chronological scheduled times.

The 10 second que is processed every second, with 10 slots in the que. The 1 hour que is processed every minute with 60 slots in the que. The 24 hour que is processed every hour, with 24 slots in the que.

#### As an example:

In the following *clock.json* config file the 10 second que is empty, the 60 minute que has one entry, to run a backup script, on the third minute of every hour. The 24 hour que will run hr1 on the third hour of the day.

~~~JSON
{
  "sec": {
    "title":"Ten Second Que", 
    "tasks": [
  ]},
  "min": {
    "title":"Sixty Minute Que", 
    "tasks": [
    {"tick":3, "script":"./scripts/backup_home.sh"}
  ]},
  "hr": {
    "title":"24 Hour Que", 
    "tasks": [
    {"tick":3, "script":"./scripts/hr1.sh"}
  ]}
}
~~~

> Currently the third hour is NOT 3AM.
> Rather it is the 3rd hour after the ticker is started.
The "Ticks" do not match chronological time!

The intent was to have something to launch various tasks at various intervals, NOT at a specific time. I am on the fence about syncing the ticks to chronological time. However, if I sync the tick or not, these circular ques are not going away, and they will always be based on intervals.

The configuration is stored in JSON format, and is intended to be manually edited. It is recommended to use vim or vi to edit the configuration. The reason for this is because TaskTick monitors this file for changes. Many editors will keep your changes in a hidden file until you exit, and vi seems to do the right thing, whereas other editors don't -- your millage may vary. Personally I've noticed that when editing the congif using VS Code causes an JSON.Unmarshall() warning, so I use vi.

The tasks are launched via GO routine (separate thread), so the main ticker loop will not be blocked. Furthermore these go routines fork/exec the task as a child process. More work will be done to add command-line flags support, system user to run as, monitoring stdin, stdout, and stderr. It will be slow, time is limited, but I am going to have fun adding these features, who knows how far I'll take this. I may even add a chronological schedule, and background tasks also. Who knows, maybe a UI frontend as well.
