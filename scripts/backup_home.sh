#!/bin/bash
notify-send  -u normal "The Scheduled Backup of the Home Directories is starting now." "\nConsider logging off"
rsync -r -t -p -o -g -v --progress -s /myhome/michael sysadm@bkupserv1:/home/sysadm/backup
notify-send  -u normal "Backup of the Home Directories has completed" " \nResume normal processing."