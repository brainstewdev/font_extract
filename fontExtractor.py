from sys import argv, exit
import os
import subprocess
import io
import datetime
import clr

from PIL import Image


def main():
    print("Starting...\n\n")
    print("getting files from input folder...")
    f = os.listdir("input_files")
    print("got ", f)
    for file in f:
        fil = open("input_files"+ os.sep + file, "rb")
        bm = Image.new( 'RGBA', (24,25), (255,255,255))
        pos = 0
        #byte = fil.read(1)
        #while True:
        mappa = bm.load()
        for y in range (25):
            for x in range(24):
                if pos % 2 == 0:
                    byte = fil.read(1)
                halfInteger = int.from_bytes(byte,byteorder='little')
                alpha = 0
                if pos % 2 == 0:
                    alpha = ((halfInteger & 0xF) << 4) | (halfInteger & 0xF)
                else:
                    alpha = (halfInteger & 0xF0) | ((halfInteger & 0xF0) >> 4)
                # bm.SetPixel(x,y, alpha)
                mappa[x,y] = (255,255,255, alpha)
                pos += 1

        bm.save("output_files" + os.sep + file.split(".bin")[0] +".png")
            #if byte == b'':
            #    break
        

    

if __name__ == "__main__":
    main()