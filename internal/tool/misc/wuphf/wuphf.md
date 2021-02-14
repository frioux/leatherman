Sends notifications on lots (2) of platforms.

Wuphf sends alerts via both `wall` and [pushover](https://pushover.net).  All
arguments are concatenated to produce the sent message.

The following environment variables should be set:

 * LM_PUSHOVER_TOKEN
 * LM_PUSHOVER_USER
 * LM_PUSHOVER_DEVICE

```bash
$ wuphf 'the shoes are on sale'
```
