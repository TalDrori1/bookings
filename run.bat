go build -o bookings.exe .\cmd\web\.
bookings.exe --dbname=bookings --dbuser=postgres --dbpass=password --cache=false --production=false