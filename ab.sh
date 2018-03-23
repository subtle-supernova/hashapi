# put post data in the file postdata
ab -n 100 -c 10 -T  'application/x-www-form-urlencoded' -p postdata http://localhost/hash
