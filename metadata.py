#!/usr/bin/env python

import urllib2
import sys

def main():
	data = urllib2.urlopen("http://blids.de/spion/bilder/").read()

	if len(data) == 0:
		sys.exit(23)
	
	files = ["aktkartebeneluxgrau.jpg", "aktkartegbgrau.jpg", "aktkartegergrau.jpg", "aktkartepolengrau.jpg", "aktkarteschweizgrau.jpg"]

	for f in files:
		pos1 = data.find(f)
		if pos1 == -1:
			continue
		
		pos2 = data.find('"right">', pos1)
		if pos2 == -1:
			continue
		pos2 += len('"right">')
		
		pos3 = data.find('</td>', pos2)
		if pos3 == -1:
			continue

		s = data[pos2:pos3]
		s = s.strip()
		a = s.split(' ')
		if len(a) != 2:
			continue
		fmt = a[0] + "," + a[1] + ";"

		fout = open(f + "_metadata.txt", 'w')
		fout.write(fmt)
		fout.close()

if __name__ == "__main__":
	main()

