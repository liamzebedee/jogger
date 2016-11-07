for i in range(1, 100):
	f = open('file'+str(i), 'w')
	f.write("FOOBAR")
	f.close()



# No longer needed.

# from subprocess import call
# call(["dd", "of=sample.txt", "bs=20G", "count=1"])