# Snoop
The idea behind snoop is to monitor an http endpoint. The program takes a url and first requests it 3 times. If the 3 responses look identical (they have the same hash), then the endpoint is "observable". Once it has done that, it will request the same endpoint and compare it with the previously obtained hash with a frequency specified by the user.
If at any moment, the hashes don't match, it will print a message to stdout and, if the credentials in the config are correctly set, send an email to whomever the user wants.
