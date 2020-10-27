#!/usr/bin/env bash
openssl req -newkey rsa:4096 -nodes -sha512 -x509 -days 3650 -nodes -out public.pem -keyout private.pem
