package config

var testToken = `aHR0cHM6Ly80NS41NS4xNTIuMjAxOjMwMDF8ZDU1ZjU1MTgtYjU2Yi00NTlhLWFhYTMtMmVmN2M5MjQxYmI3fE1tWmhNbU15TldFdFptRTRaUzAwTUdNNExXRTNZMkl0WVRBek56aGpNRFZrWXpZNUNnPT18LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUVyRENDQXBRQ0NRRDhrMEhlYmthbjhqQU5CZ2txaGtpRzl3MEJBUVVGQURBWU1SWXdGQVlEVlFRRERBMDAKTlM0MU5TNHhOVEl1TWpBeE1CNFhEVEUxTURNeU56RTROVEl3T1ZvWERURTJNRE15TmpFNE5USXdPVm93R0RFVwpNQlFHQTFVRUF3d05ORFV1TlRVdU1UVXlMakl3TVRDQ0FpSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnSVBBRENDCkFnb0NnZ0lCQU43UEtpQzY3akRIc3BuazZneE12Qk9nTkJMRWJLYndWVTlTZ0gzeVBLUCtGUVVXL1ZQL2tTMXMKRUlOWFhFRjZkNnlZZXpzY1pZNjMzTE9DUHllanFHYzBZZy96d01nWk1wdXphSnJHYWZUMG4vRnA1Vy9nYm4vSgpFOHBKdG1NU0J0L3VNelJxUWNBWFpyVUtmTnZJaStPUFliVko1SFFNWTdjUzc0blh0bUpWL0MyaFVJMHJIN1R5ClR2dTBOZy9oRlpmYjFpdEE4azNsSVJqdHkreWtZME15RmpJa3plYVJONGNzZTcxWnBibVgxbzF4QVFETEpORnoKSTdOZlhzMktFcGRXRFY1ZjdSRThraVY0Q3dYL0VpVDRzYVE4djg0QUI3MHR0TlZPbkJmS3lQZnYrbkNVQUFLNQpGaVRrNFd2SUlpS0gySlc4bW5qblhlL05Od2RoSGJPVUE4aXBoOXVjb3pNTTlTeG9EZVRpTC8zeEg4QlB3YVBmCk8xSFVFajNNbVJCTWxrRHR4dkZvM2wrbjVyajNjclUrVUxlTVBWckpOUnRmREZmUmVLakxseHlRY2FaNWtJMHEKYVRDRUNzK1A4d0t2UWFoNHdVR3daWDVydkh0eG5FeGFVdjRpQzVCZVBrZ1Y5ZElobTdQVFJQQ0k0K1p3d0QwUApMRmlWZ0ZmNWoyaXZEQzN6MVpPeG44Sy9kUmwrZmc3b0E1N1c2Y3MwUWxmR3lEazEzVmpGeU9KamVwK0VDNVY5CnhlcW11RW9MdENPa1J5aWE3UGwzR3BDa25VTmFsUk5ub3lIbURZck1TcFpIdE5SeUI2UXBOZmJGcUJ1MHNnNW4KdEJ4NEthM2hEcmNCdVV0UWx6Z0pCZjFRNUlVTXFZU3NwSDNFUWxZSWM0OXg0MzY1RWZiN0FnTUJBQUV3RFFZSgpLb1pJaHZjTkFRRUZCUUFEZ2dJQkFHdjdaNCtCSnBldDBKcFcwTWhmWjVDZmdoL3BKeFNqTEpCbUV0NDNWKytIClg3VEhXZ3BCRHZ2QVhSUzRoY001dmwwYjRtRDJ3RzZmTys3Z0labWVpZEk2UHlLenlVSVBrN0NlUzlQU1lqRHgKZEVOeXdob2pWV2Qzb3p1QmZLcjB5MEVtQ1Jqd1NxY1NabWdTQXRwMVlKeE83bmxON3F4bnlNWi8vbzd0RzRPZwo2eTRNZG1XNWFEdmNVQXZ1aG8rSkF4UktLZk5QM290Qkc5alNzNHBDQm51L0ZleEd5VVhUaC9JSS9iQWorMUtFCmw0aHQvdEVTUm9NZEZRN3drMlpzYXB1eStvNklzdDRxVnRYU1ZsU0F1K1hvanB1aVZsNmVhZWltQUozTEJRMXcKNUlnczlsUnNaZEJOR2Y3STlqbDVNZndFMURUbzlnMC9tcTRGcUROZHVVSWJNaHRZUVpxdEhEVWI1N25MajR0YgptVVNVTG9aUk5sSzJyTURSRlhDMERWMU5GRVo1MFIzQXJpdVdVK0JqYTdYeU96WkFWVWZEWlFOajhuczhMR3JSCnpCTXlVWEZIbVpCNmt1NHNEc3EwWnE1cno4RHE3dUhvcUxOa0ZRM05oZ1pzNEgvZGtZRHFRamlrY0taUkVmRG4KL1NJK0lvTW5GbjFLamhXMXlCZ3JFdkNrNzBMSWhRTHhnalJZUDU3VXhaREhEMVBrVVhHY1Z4U3l4dEhjaHIwbwp5c0liQllkWjdGRTFOQXFNTmswR1U0TTB3RUZlekxGRWhlRE9NYndabXpPL2RsaVkzTkE5dFY3OHkwcDJiQWhCCkZIWUQ2V3BIWTNnVUlFcnVUVjRrWktNMk5SdWI5VUcrbGpTYy9iUDk4emJDK0RaN2pNbmtrTkdsdmY2cFRkM2cKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=`