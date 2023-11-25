#!/bin/bash

psql -d postgres -a -f ./postgres/init_users_table.sql
psql -d postgres -a -f ./postgres/init_sessions_table.sql