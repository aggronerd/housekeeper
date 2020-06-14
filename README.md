# housekeeper

Command line tool to help with housekeeping tasks that are boooorrinnnngggg.

## Configuration

Ensure you have a config file in `~/.housekeeper.yaml` this needs to have your Jira credentials, for
example:

```
jira:
  username: example@example.com
  password: apasswordortokenfromjira
  url: https://yourdomain.atlassian.net
  auth: basic
```

Currently only `basic` authentication is supported. Be sure to `chmod 600` this file.

## Usage

The following example prints out a daily report of work logged against tickets:

```bash
housekeeper time report

  ISSUE ID |       SUMMARY        | TIME SPENT
-----------+----------------------+-------------
  TP-1     | Example task         | 2h0m0s
  TP-2     | Another example task | 1h30m0s
  TP-3     | A third dask         | 2h10m0s
-----------+----------------------+-------------
                    TOTAL         |  5H40M0S
           -----------------------+-------------
```

To see a full list of commands run: `housekeeper help`.