#!/bin/bash
mongo --eval "
  use shelob,
  db.proxy.createIndex({ ip_address: 1, port: 1 }, { unique: true });
"