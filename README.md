# Træfik log multiplexer

It reads Traefik's logs file, listen Docker events and join them.
You can now access docker's labels in your log workflow.
Main target is lots of docker-compose projects behind a Traefik.

There is multiple outputs.

**File**, for writing one log file per project.

**Fluent**, export logs somewhere, through fluent protocol. Target can be fluent-bit, fluentd, Loki…

**stdout**, everybody loves debug.

**prometheus**, for every a prometheus export with a password. Hit per *status*. Latencies per *status family*, per *method*.
