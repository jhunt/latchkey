latchkey
========

Automation pipelines are great.  Automation pipelines that notify
are even better.  Automation pipelines that notify on
break/recovery vs.  fail/success?  Better still.

`latchkey` aims to help CI/CD software to get to that final stage,
where notifications only go out when the pipeline stops working,
or resumes normal operations.  I'm tired of getting multiple fail
notifications in a row as I try to debug and fix any given
automaton.

So how does it work?

`latchkey` provides its full URI space for you to use as keys.
PUT'ing to a key with a value in the body will set the key, and
atomically let you know if that value is new or not.

Here's an example:

```
curl -X PUT http://localhost/some/pipeline/status \
  --data-binary "ok"
changed

curl -X PUT http://localhost/some/pipeline/status \
  --data-binary "ok"
same

curl -X PUT http://localhost/some/pipeline/status \
  --data-binary "ok"
same

curl -X PUT http://localhost/some/pipeline/status \
  --data-binary "fail"
changed

curl -X PUT http://localhost/some/pipeline/status \
  --data-binary "fail"
same

curl -X PUT http://localhost/some/pipeline/status \
  --data-binary "ok"
changed
```

If `latchkey` responds to your PUT with a 200 and the body
"same\n", then the key did **not** transition to a new value.  If
it responds with a 200 and the body "changed\n", a transition
occurred.

(Note: This means you cannot use the values "same" and "changed"
       for your own purposes.  In practice, this has not been a
       deal-breaker.)

Running latchkey
----------------

```
LATCHKEY_BIND=:80 latchkey
```

By default, `latchkey` binds all interfaces, on port 8080.

Caveats
-------

`latchkey` uses memory for storage, it does not persist the key
values to disk.  If you restart the `latchkey` process,
_EVERYTHING_ resets.
