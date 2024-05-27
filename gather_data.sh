#!/bin/sh

./thermpro_exporter -start "${start_year}-${start_month}-${start_day}T00:00:00Z" -end "${end_year}-${end_month}-${end_day}T00:00:00Z"
open out.csv
exit 0
---
[start_year]
  label="Enter Start Year"
  required=true
[start_month]
  label="Enter Start Month"
  required=true
[start_day]
  label="Enter Start Day"
  required=true
[end_year]
  label="Enter End Year"
  required=true
[end_month]
  label="Enter End Month"
  required=true
[end_day]
  label="Enter End Day"
  required=true



